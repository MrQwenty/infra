package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/coneno/logger"
	"github.com/influenzanet/messaging-service/pkg/dbs/globaldb"
	"github.com/influenzanet/messaging-service/pkg/dbs/messagedb"
	"github.com/influenzanet/messaging-service/pkg/grpc/clients"
	"github.com/influenzanet/messaging-service/pkg/retry"
	"github.com/influenzanet/messaging-service/pkg/types"
	whatsappAPI "github.com/influenzanet/messaging-service/pkg/api/whatsapp_client_service"
)

func main() {
	logger.Info.Println("Starting WhatsApp message scheduler...")

	globalDBCon := os.Getenv("GLOBAL_DB_CONNECTION_STR")
	globalDBUsername := os.Getenv("GLOBAL_DB_USERNAME")
	globalDBPassword := os.Getenv("GLOBAL_DB_PASSWORD")
	globalDBConnectionPrefix := os.Getenv("GLOBAL_DB_CONNECTION_PREFIX")

	globalDB := globaldb.NewGlobalDBService(
		globalDBCon,
		globalDBUsername,
		globalDBPassword,
		globalDBConnectionPrefix,
		"globalinfodb",
	)

	messageDBCon := os.Getenv("MESSAGE_DB_CONNECTION_STR")
	messageDBUsername := os.Getenv("MESSAGE_DB_USERNAME")
	messageDBPassword := os.Getenv("MESSAGE_DB_PASSWORD")
	messageDBConnectionPrefix := os.Getenv("MESSAGE_DB_CONNECTION_PREFIX")

	messageDB := messagedb.NewMessageDBService(
		messageDBCon,
		messageDBUsername,
		messageDBPassword,
		messageDBConnectionPrefix,
		"messagedb",
	)

	whatsappClientAddr := os.Getenv("WHATSAPP_CLIENT_SERVICE_LISTEN_PORT")
	if whatsappClientAddr == "" {
		whatsappClientAddr = "localhost:5007"
	} else {
		whatsappClientAddr = "localhost:" + whatsappClientAddr
	}

	whatsappClient, closeWhatsAppClient := clients.ConnectToWhatsAppClientService(whatsappClientAddr)
	defer closeWhatsAppClient()

	intervalStr := os.Getenv("WHATSAPP_MESSAGE_SCHEDULER_INTERVAL")
	if intervalStr == "" {
		intervalStr = "10"
	}
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		logger.Error.Fatalf("Invalid WHATSAPP_MESSAGE_SCHEDULER_INTERVAL: %v", err)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			instances, err := globalDB.GetAllInstances()
			if err != nil {
				logger.Error.Printf("Failed to get instances: %v", err)
				continue
			}

			for _, instance := range instances {
				processInstanceWhatsAppMessages(messageDB, whatsappClient, instance.InstanceID)
			}
		}
	}
}

func processInstanceWhatsAppMessages(messageDB *messagedb.MessageDBService, whatsappClient whatsappAPI.WhatsAppClientServiceApiClient, instanceID string) {
	batchSize := 50
	retryDelay := int64(300)

	messages, err := messageDB.FetchOutgoingWhatsApp(instanceID, batchSize, retryDelay, true)
	if err != nil {
		logger.Error.Printf("Failed to fetch high priority WhatsApp messages for %s: %v", instanceID, err)
		return
	}

	for _, message := range messages {
		processWhatsAppMessage(messageDB, whatsappClient, instanceID, message)
	}

	messages, err = messageDB.FetchOutgoingWhatsApp(instanceID, batchSize, retryDelay, false)
	if err != nil {
		logger.Error.Printf("Failed to fetch normal priority WhatsApp messages for %s: %v", instanceID, err)
		return
	}

	for _, message := range messages {
		processWhatsAppMessage(messageDB, whatsappClient, instanceID, message)
	}
}

func processWhatsAppMessage(messageDB *messagedb.MessageDBService, whatsappClient whatsappAPI.WhatsAppClientServiceApiClient, instanceID string, message types.OutgoingWhatsApp) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := whatsappClient.SendWhatsApp(ctx, &whatsappAPI.SendWhatsAppReq{
		To:       message.To,
		Content:  message.Content,
		HighPrio: message.HighPrio,
	})

	if err != nil {
		logger.Error.Printf("Failed to send WhatsApp message %s: %v", message.ID.Hex(), err)
		
		errorCategory := retry.CategorizeError(err)
		
		if !retry.ShouldRetry(errorCategory) {
			logger.Warning.Printf("WhatsApp message %s has non-retryable error, deleting", message.ID.Hex())
			if delErr := messageDB.DeleteOutgoingWhatsApp(instanceID, message.ID.Hex()); delErr != nil {
				logger.Error.Printf("Failed to delete WhatsApp message %s: %v", message.ID.Hex(), delErr)
			}
			return
		}
		
		if message.RetryCount >= message.MaxRetries {
			logger.Warning.Printf("WhatsApp message %s exceeded max retries, deleting", message.ID.Hex())
			if delErr := messageDB.DeleteOutgoingWhatsApp(instanceID, message.ID.Hex()); delErr != nil {
				logger.Error.Printf("Failed to delete WhatsApp message %s: %v", message.ID.Hex(), delErr)
			}
		} else {
			baseDelay := message.BaseDelaySeconds
			if baseDelay <= 0 {
				baseDelay = retry.DefaultBaseDelay
			}
			
			delaySeconds := retry.CalculateNextRetryDelay(message.RetryCount, baseDelay, errorCategory)
			nextRetryAt := time.Now().Unix() + delaySeconds
			
			if retryErr := messageDB.IncrementWhatsAppRetryCountWithBackoff(instanceID, message.ID.Hex(), nextRetryAt, string(errorCategory)); retryErr != nil {
				logger.Error.Printf("Failed to increment retry count for WhatsApp message %s: %v", message.ID.Hex(), retryErr)
			} else {
				logger.Info.Printf("WhatsApp message %s scheduled for retry in %d seconds (attempt %d/%d)", 
					message.ID.Hex(), delaySeconds, message.RetryCount+1, message.MaxRetries)
			}
		}
		return
	}

	logger.Debug.Printf("WhatsApp message sent successfully: %s", message.ID.Hex())
	
	if delErr := messageDB.DeleteOutgoingWhatsApp(instanceID, message.ID.Hex()); delErr != nil {
		logger.Error.Printf("Failed to delete sent WhatsApp message %s: %v", message.ID.Hex(), delErr)
	}

	if _, addErr := messageDB.AddToSentWhatsApp(instanceID, message); addErr != nil {
		logger.Error.Printf("Failed to add WhatsApp message to sent collection: %v", addErr)
	}
}
