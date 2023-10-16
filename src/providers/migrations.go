package providers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/spf13/viper"
)

const (
	dynamoDBDefaultRCU int64 = 2
	dynamoDBDefaultWCU int64 = 2
)

// CreateTables create the necessary DynamoDB table schema
func CreateTables(_ context.Context, ddb dynamodbiface.DynamoDBAPI, cfg *viper.Viper) (err error) {
	_, err = ddb.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(cfg.GetString(ConfigKeyAccountsTableName)),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(dynamoDBDefaultRCU),
			WriteCapacityUnits: aws.Int64(dynamoDBDefaultWCU),
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("AccountID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("AccountID"),
				KeyType:       aws.String("HASH"),
			},
		},
	})
	if err != nil {
		log.Println("Could not create accounts table", err)
	}

	_, err = ddb.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(cfg.GetString(ConfigKeyTransactionsTableName)),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(dynamoDBDefaultRCU),
			WriteCapacityUnits: aws.Int64(dynamoDBDefaultWCU),
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("AccountID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("TransactionID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("TransactionID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("AccountID"),
				KeyType:       aws.String("RANGE"),
			},
		},
	})
	if err != nil {
		log.Println("Could not create transactions table", err)
	}

	return nil
}
