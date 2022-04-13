package mongo

import (
	"context"
	"os"
	"simple-upload/pkg/util"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Client interface {
	StartSession(...*options.SessionOptions) (mongo.Session, error)
	Disconnect(context.Context) error
	NumberSessionsInProgress() int
	Database(string, ...*options.DatabaseOptions) *mongo.Database
}

func GetClientOptions() (*options.ClientOptions, *connstring.ConnString, error) {
	mongoDSN := os.Getenv("STORAGE_MONGO_DSN")
	connString, err := connstring.ParseAndValidate(mongoDSN)
	if err != nil {
		return nil, nil, err
	}

	opt := options.Client().ApplyURI(mongoDSN)

	timeoutRaw, _ := strconv.Atoi(os.Getenv("STORAGE_MONGO_TIMEOUT_IN_SECONDS"))
	opt.ConnectTimeout = util.ConvertSecondsToDuration(timeoutRaw)

	return opt, &connString, nil
}
