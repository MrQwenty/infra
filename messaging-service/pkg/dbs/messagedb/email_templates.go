package messagedb

import (
	"errors"

	"github.com/influenzanet/messaging-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *MessageDBService) SaveEmailTemplate(instanceID string, template types.EmailTemplate) (types.EmailTemplate, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{
		"messageType": template.MessageType,
		"studyKey":    template.StudyKey,
	}
	if template.StudyKey == "" {
		filter["studyKey"] = bson.M{"$exists": false}
	}

	upsert := true
	rd := options.After
	options := options.FindOneAndReplaceOptions{
		Upsert:         &upsert,
		ReturnDocument: &rd,
	}
	elem := types.EmailTemplate{}
	err := dbService.collectionRefEmailTemplates(instanceID).FindOneAndReplace(
		ctx, filter, template, &options,
	).Decode(&elem)
	return elem, err
}

func (dbService *MessageDBService) DeleteEmailTemplate(instanceID string, messageType string, studyKey string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{
		"messageType": messageType,
		"studyKey":    studyKey,
	}
	if studyKey == "" {
		filter["studyKey"] = bson.M{"$exists": false}
	}
	res, err := dbService.collectionRefEmailTemplates(instanceID).DeleteOne(ctx, filter)
	if res.DeletedCount < 1 {
		err = errors.New("not found")
	}
	return err
}

func (dbService *MessageDBService) FindEmailTemplateByType(instanceID string, messageType string, studyKey string) (types.EmailTemplate, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{
		"messageType": messageType,
		"studyKey":    studyKey,
	}
	if studyKey == "" {
		filter["studyKey"] = bson.M{"$exists": false}
	}

	elem := types.EmailTemplate{}
	err := dbService.collectionRefEmailTemplates(instanceID).FindOne(ctx, filter).Decode(&elem)
	return elem, err
}

func (dbService *MessageDBService) FindAllEmailTempates(instanceID string) (templates []types.EmailTemplate, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{}
	cur, err := dbService.collectionRefEmailTemplates(instanceID).Find(
		ctx,
		filter,
	)

	if err != nil {
		return templates, err
	}
	defer cur.Close(ctx)

	templates = []types.EmailTemplate{}
	for cur.Next(ctx) {
		var result types.EmailTemplate
		err := cur.Decode(&result)
		if err != nil {
			return templates, err
		}

		templates = append(templates, result)
	}
	if err := cur.Err(); err != nil {
		return templates, err
	}

	return templates, nil
}
