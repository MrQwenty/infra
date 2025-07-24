package messaging_service

import (
	"context"

	"github.com/coneno/logger"
	loggingAPI "github.com/influenzanet/logging-service/pkg/api"
)

func (s *messagingServer) SaveLogEvent(instanceID, userID string, eventType loggingAPI.LogEventType, eventName, msg string) {
	_, err := s.clients.LoggingService.SaveLogEvent(context.Background(), &loggingAPI.NewLogEvent{
		Origin:     "messaging-service",
		InstanceId: instanceID,
		UserId:     userID,
		EventType:  eventType,
		EventName:  eventName,
		Msg:        msg,
	})
	if err != nil {
		logger.Error.Printf("failed to save log: %v", err)
	}
}
