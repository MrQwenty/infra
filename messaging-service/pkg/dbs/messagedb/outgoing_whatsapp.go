package messagedb

import (
	"errors"
	"time"

	"github.com/influenzanet/messaging-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (dbService *MessageDBService) AddToOutgoingWhatsApp(instanceID string, message types.OutgoingWhatsApp) (types.OutgoingWhatsApp, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	if message.AddedAt <= 0 {
		message.AddedAt = time.Now().Unix()
	}
	if message.MaxRetries <= 0 {
		message.MaxRetries = 5
	}

	res, err := dbService.collectionRefOutgoingWhatsApp(instanceID).InsertOne(ctx, message)
	if err != nil {
		return message, err
	}
	message.ID = res.InsertedID.(primitive.ObjectID)
	return message, nil
}

func (dbService *MessageDBService) AddToSentWhatsApp(instanceID string, message types.OutgoingWhatsApp) (types.OutgoingWhatsApp, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()
	message.AddedAt = time.Now().Unix()
	message.Content = ""

	message.ID = primitive.NilObjectID
	res, err := dbService.collectionRefSentWhatsApp(instanceID).InsertOne(ctx, message)
	if err != nil {
		return message, err
	}
	message.ID = res.InsertedID.(primitive.ObjectID)
	return message, nil
}

func (dbService *MessageDBService) FetchOutgoingWhatsApp(instanceID string, amount int, olderThan int64, onlyHighPrio bool) (messages []types.OutgoingWhatsApp, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	counter := 0
	for counter < amount {
		var newMessage types.OutgoingWhatsApp
		update := bson.M{"$set": bson.M{"lastSendAttempt": time.Now().Unix()}}
		filter := bson.M{
			"lastSendAttempt": bson.M{"$lt": time.Now().Unix() - olderThan},
			"retryCount":      bson.M{"$lt": "$maxRetries"},
		}
		if onlyHighPrio {
			filter["highPrio"] = true
		}
		if err := dbService.collectionRefOutgoingWhatsApp(instanceID).FindOneAndUpdate(ctx, filter, update).Decode(&newMessage); err != nil {
			break
		}
		messages = append(messages, newMessage)
		counter += 1
	}
	return messages, nil
}

func (dbService *MessageDBService) IncrementWhatsAppRetryCount(instanceID string, id string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{
		"$inc": bson.M{"retryCount": 1},
		"$set": bson.M{"lastSendAttempt": 0},
	}

	res, err := dbService.collectionRefOutgoingWhatsApp(instanceID).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount < 1 {
		return errors.New("no outgoing WhatsApp message found with the given id")
	}
	return nil
}

func (dbService *MessageDBService) ResetLastSendAttemptForOutgoingWhatsApp(instanceID string, id string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"lastSendAttempt": 0}}

	res, err := dbService.collectionRefOutgoingWhatsApp(instanceID).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount < 1 {
		return errors.New("no outgoing WhatsApp message found with the given id")
	}
	return nil
}

func (dbService *MessageDBService) DeleteOutgoingWhatsApp(instanceID string, id string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	res, err := dbService.collectionRefOutgoingWhatsApp(instanceID).DeleteOne(ctx, filter, nil)
	if err != nil {
		return err
	}
	if res.DeletedCount < 1 {
		return errors.New("no outgoing WhatsApp message found with the given id")
	}
	return nil
}
