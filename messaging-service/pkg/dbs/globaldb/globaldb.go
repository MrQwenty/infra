package globaldb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GlobalDBService struct {
	DBClient     *mongo.Client
	timeout      int
	DBNamePrefix string
}

type Instance struct {
	InstanceID string `bson:"instanceId"`
}

func NewGlobalDBService(
	connectionStr string,
	username string,
	password string,
	connectionPrefix string,
	dbNamePrefix string,
) *GlobalDBService {
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

	return &GlobalDBService{
		DBClient:     dbClient,
		timeout:      30,
		DBNamePrefix: dbNamePrefix,
	}
}

func (dbService *GlobalDBService) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(dbService.timeout)*time.Second)
}

func (dbService *GlobalDBService) GetAllInstances() ([]Instance, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	collection := dbService.DBClient.Database(dbService.DBNamePrefix).Collection("instances")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var instances []Instance
	if err = cursor.All(ctx, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}
