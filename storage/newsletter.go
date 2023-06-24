package storage

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"canvas/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"
)

// SignupForNewsletter with the given email. Returns a token used for confirming the email address.
const tableName = "newsletter_subscribers"

type NewsletterEntity struct {
	Email     string    `dynamodbav:"email"`
	Token     string    `dynamodbav:"token"`
	Confirmed bool      `dynamodbav:"confirmed"`
	Active    bool      `dynamodbav:"active"`
	Created   time.Time `dynamodbav:"created"`
	Updated   time.Time `dynamodbav:"updated"`
}

func (d *Database) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	token, err := createSecret()
	if err != nil {
		return "", err
	}
	av, err := dynamodbattribute.MarshalMap(NewsletterEntity{
		Email:     email.String(),
		Token:     token,
		Confirmed: false,
		Active:    true,
		Created:   time.Now(),
		Updated:   time.Now(),
	})

	if err != nil {
		d.log.Error("Got error marshalling new movie item", zap.Error(err))
	}

	_, err = d.DB.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})

	if err != nil {
		d.log.Error("Error PutItem", zap.Error(err))
	}
	return token, err
}

func createSecret() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", secret), nil
}
