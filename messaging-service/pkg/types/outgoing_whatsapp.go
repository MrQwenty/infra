package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OutgoingWhatsApp struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	MessageType      string             `bson:"messageType"`
	To               string             `bson:"to"`
	Content          string             `bson:"content"`
	AddedAt          int64              `bson:"addedAt"`
	HighPrio         bool               `bson:"highPrio"`
	LastSendAttempt  int64              `bson:"lastSendAttempt"`
	RetryCount       int                `bson:"retryCount"`
	MaxRetries       int                `bson:"maxRetries"`
	NextRetryAt      int64              `bson:"nextRetryAt"`
	BaseDelaySeconds int                `bson:"baseDelaySeconds"`
	LastErrorType    string             `bson:"lastErrorType"`
}
