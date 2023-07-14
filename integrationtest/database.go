package integrationtest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"canvas/storage"
)

var once sync.Once

// CreateDatabase for testing.
// Usage:
//
//	db, cleanup := CreateDatabase()
//	defer cleanup()
//	…

type MockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m *MockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

func CreateDatabase() (*storage.Database, func()) {

	once.Do(initDatabase)

	db, cleanup := connect("dynamicdb")
	defer cleanup()

	dropConnections(db.DB)

	dropConnections(db.DB)

	return connect("dynamodb")
}

func initDatabase() {
	db, cleanup := connect("template1")
	defer cleanup()

	for err := db.Ping(context.Background()); err != nil; {
		time.Sleep(100 * time.Millisecond)
	}
}

func connect(name string) (*storage.Database, func()) {
	db := storage.NewDatabase(storage.NewDatabaseOptions{DB: &MockDynamoDBClient{}})
	if err := db.Connect(); err != nil {
		fmt.Println("Error connecting to database:", err)
		panic(err)
	}
	return db, func() {}
}

func dropConnections(db dynamodbiface.DynamoDBAPI) {}
