package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.uber.org/zap"
)

// Database is the relational storage abstraction.
type Database struct {
	DB  *dynamodb.DynamoDB
	log *zap.Logger
}

// NewDatabaseOptions for NewDatabase.
type NewDatabaseOptions struct {
	Log *zap.Logger
}

// NewDatabase with the given options.
// If no logger is provided, logs are discarded.
func NewDatabase(opts NewDatabaseOptions) *Database {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	return &Database{
		log: opts.Log,
	}
}

// Connect to the database.
func (d *Database) Connect() error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error

	mySession := session.Must(session.NewSession())

	if err != nil {
		return err
	}

	d.DB = dynamodb.New(mySession, aws.NewConfig().WithRegion("eu-central-1"))

	d.log.Debug("Setting connection to dynamodb")

	return nil
}
