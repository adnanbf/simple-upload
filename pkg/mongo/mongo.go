package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type MongoDB interface {
	Close(ctx context.Context) error
	Find(ctx context.Context, tableName string, filter bson.M, result interface{}) error
	Insert(ctx context.Context, tableName string, data interface{}) error
	Aggregate(ctx context.Context, tableName string, pipeline []bson.M, result interface{}) error
}

type MongoDBImpl struct {
	Client     Client
	ConnString connstring.ConnString
}

func NewMongoDB() (MongoDB, error) {
	clientOption, connString, err := GetClientOptions()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *clientOption.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database")
	}

	return &MongoDBImpl{
		Client:     client,
		ConnString: *connString,
	}, nil
}

func (m *MongoDBImpl) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

func (q *MongoDBImpl) Find(ctx context.Context, tableName string, filter bson.M, result interface{}) error {
	// log.Debugf("mongo query => find: %s, query: %s \n", tableName, utility.ConvertBsonMToJSONString(filter))

	csr, err := q.Client.
		Database(q.ConnString.Database).
		Collection(tableName).
		Find(ctx, filter)
	if csr != nil {
		defer csr.Close(ctx)
	}
	if err != nil {
		return err
	}

	err = csr.All(ctx, result)
	if err != nil {
		return err
	}

	return nil
}

func (q *MongoDBImpl) Insert(ctx context.Context, tableName string, data interface{}) error {
	// log.Debugf("mongo query => insert: %s, data: %s \n", tableName, utility.ConvertBsonMToJSONString(data))

	_, err := q.Client.
		Database(q.ConnString.Database).
		Collection(tableName).
		InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (q *MongoDBImpl) Aggregate(ctx context.Context, tableName string, pipeline []bson.M, result interface{}) error {
	// log.Debugf("mongo query => aggregate: %s, pipeline: %s \n", tableName, utility.ConvertBsonMToJSONString(pipeline))

	csr, err := q.Client.
		Database(q.ConnString.Database).
		Collection(tableName).
		Aggregate(ctx, pipeline)
	if csr != nil {
		defer csr.Close(ctx)
	}
	if err != nil {
		return err
	}

	err = csr.All(ctx, result)
	if err != nil {
		return err
	}

	return nil
}
