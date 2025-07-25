package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type WhatsAppClient struct {
	baseURL       string
	accessToken   string
	phoneNumberID string
	httpClient    *http.Client
}

type WhatsAppMessage struct {
	To       string                 `json:"to"`
	Type     string                 `json:"type"`
	Template WhatsAppTemplate       `json:"template"`
}

type WhatsAppTemplate struct {
	Name       string                    `json:"name"`
	Language   WhatsAppLanguage          `json:"language"`
	Components []WhatsAppTemplateComponent `json:"components"`
}

type WhatsAppLanguage struct {
	Code string `json:"code"`
}

type WhatsAppTemplateComponent struct {
	Type       string                      `json:"type"`
	Parameters []WhatsAppTemplateParameter `json:"parameters"`
}

type WhatsAppTemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type WhatsAppResponse struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	Contacts []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
}

type WhatsAppError struct {
	Error struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Code      int    `json:"code"`
		ErrorData struct {
			MessagingProduct string `json:"messaging_product"`
			Details          string `json:"details"`
		} `json:"error_data"`
	} `json:"error"`
}

func NewWhatsAppClient() *WhatsAppClient {
	return &WhatsAppClient{
		baseURL:       "https://graph.facebook.com/v19.0",
		accessToken:   os.Getenv("WHATSAPP_API_TOKEN"),
		phoneNumberID: os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (w *WhatsAppClient) SendVerificationCode(ctx context.Context, phoneNumber, code string) error {
	if w.accessToken == "" || w.phoneNumberID == "" {
		return fmt.Errorf("WhatsApp API credentials not configured")
	}

	// Clean phone number (remove any formatting)
	cleanPhone := w.cleanPhoneNumber(phoneNumber)
	
	message := WhatsAppMessage{
		To:   cleanPhone,
		Type: "template",
		Template: WhatsAppTemplate{
			Name: "hello_world",
			Language: WhatsAppLanguage{
				Code: "en_US",
			},
			Components: []WhatsAppTemplateComponent{
				{
					Type: "body",
					Parameters: []WhatsAppTemplateParameter{
						{
							Type: "text",
							Text: fmt.Sprintf("Your InfluenzaNet verification code is: %s. This code will expire in 10 minutes.", code),
						},
					},
				},
			},
		},
	}

	return w.sendMessage(ctx, message)
}

func (w *WhatsAppClient) sendMessage(ctx context.Context, message WhatsAppMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", w.baseURL, w.phoneNumberID)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.accessToken))

	log.Printf("Sending WhatsApp message to %s", message.To)
	log.Printf("Request URL: %s", url)
	log.Printf("Request Body: %s", string(jsonData))

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("WhatsApp API Response Status: %d", resp.StatusCode)
	log.Printf("WhatsApp API Response Body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		var whatsappError WhatsAppError
		if err := json.Unmarshal(body, &whatsappError); err == nil {
			return fmt.Errorf("WhatsApp API error: %s (code: %d)", 
				whatsappError.Error.Message, whatsappError.Error.Code)
		}
		return fmt.Errorf("WhatsApp API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var response WhatsAppResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Messages) == 0 {
		return fmt.Errorf("no message ID returned from WhatsApp API")
	}

	log.Printf("WhatsApp message sent successfully. Message ID: %s", response.Messages[0].ID)
	return nil
}

func (w *WhatsAppClient) cleanPhoneNumber(phoneNumber string) string {
	// Remove any non-digit characters except the leading +
	cleaned := ""
	for i, char := range phoneNumber {
		if i == 0 && char == '+' {
			cleaned += string(char)
		} else if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}
	return cleaned
}

func (w *WhatsAppClient) SendWithRetry(ctx context.Context, phoneNumber, code string, maxRetries int) error {
	var lastErr error
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 30s, 1m, 2m
			delays := []time.Duration{30 * time.Second, 60 * time.Second, 120 * time.Second}
			delay := delays[attempt-1]
			if attempt-1 >= len(delays) {
				delay = delays[len(delays)-1]
			}
			
			log.Printf("Retrying WhatsApp send in %v (attempt %d/%d)", delay, attempt+1, maxRetries+1)
			
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
		
		err := w.SendVerificationCode(ctx, phoneNumber, code)
		if err == nil {
			return nil
		}
		
		lastErr = err
		log.Printf("WhatsApp send attempt %d failed: %v", attempt+1, err)
	}
	
	return fmt.Errorf("failed to send WhatsApp message after %d attempts: %w", maxRetries+1, lastErr)
}