package messagedb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageDBService struct {
	DBClient   *mongo.Client
	timeout    int
	DBNamePrefix string
}

func NewMessageDBService(
	connectionStr string,
	username string,
	password string,
	connectionPrefix string,
	dbNamePrefix string,
) *MessageDBService {
	var err error
	dbClient, err := mongo.NewClient(options.Client().ApplyURI(connectionStr))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = dbClient.Connect(ctx)
	if err != nil {
		panic(err)
	}

	return &MessageDBService{
		DBClient:     dbClient,
		timeout:      30,
		DBNamePrefix: dbNamePrefix,
	}
}

func (dbService *MessageDBService) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(dbService.timeout)*time.Second)
}

func (dbService *MessageDBService) collectionRefOutgoingWhatsApp(instanceID string) *mongo.Collection {
	return dbService.DBClient.Database(dbService.DBNamePrefix + "_" + instanceID).Collection("outgoing_whatsapp")
}

func (dbService *MessageDBService) collectionRefSentWhatsApp(instanceID string) *mongo.Collection {
	return dbService.DBClient.Database(dbService.DBNamePrefix + "_" + instanceID).Collection("sent_whatsapp")
}

func (dbService *MessageDBService) FindEmailTemplateByType(instanceID, messageType, studyKey string) (*EmailTemplate, error) {
	return &EmailTemplate{
		MessageType:     messageType,
		DefaultLanguage: "en",
		Translations: []LocalizedTemplate{
			{
				Lang:        "en",
				Subject:     "Notification",
				TemplateDef: "VGVzdCBtZXNzYWdl", // base64 encoded "Test message"
			},
		},
	}, nil
}

func (dbService *MessageDBService) AddToOutgoingEmails(instanceID string, email OutgoingEmail) (OutgoingEmail, error) {
	return email, nil
}

func (dbService *MessageDBService) AddToSentEmails(instanceID string, email OutgoingEmail) (OutgoingEmail, error) {
	return email, nil
}

type EmailTemplate struct {
	MessageType     string
	StudyKey        string
	DefaultLanguage string
	HeaderOverrides HeaderOverrides
	Translations    []LocalizedTemplate
}

type LocalizedTemplate struct {
	Lang        string
	Subject     string
	TemplateDef string
}

type HeaderOverrides struct {
	From      string
	Sender    string
	ReplyTo   []string
	NoReplyTo bool
}

func (h HeaderOverrides) ToEmailClientAPI() interface{} {
	return h
}

type OutgoingEmail struct {
	MessageType     string
	To              []string
	HeaderOverrides HeaderOverrides
	Subject         string
	Content         string
	HighPrio        bool
}
