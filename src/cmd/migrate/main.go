package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chff7cb/swissbank/providers"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
)

func main() {
	fx.New(
		// here we will wire up application layers
		fx.Provide(
			// load base configuration
			providers.ViperConfigProvider,
			// instantiate database clients
			providers.DynamoDBProvider,
		),
		fx.Invoke(func(ddb dynamodbiface.DynamoDBAPI, cfg *viper.Viper) (err error) {
			_, err = ddb.CreateTable(&dynamodb.CreateTableInput{
				TableName: aws.String(cfg.GetString(providers.ConfigKeyAccountsTableName)),
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(2),
					WriteCapacityUnits: aws.Int64(2),
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
				TableName: aws.String(cfg.GetString(providers.ConfigKeyTransactionsTableName)),
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(2),
					WriteCapacityUnits: aws.Int64(2),
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
						KeyType:       aws.String("SORT"),
					},
				},
			})
			if err != nil {
				log.Println("Could not create transactions table", err)
			}

			return
		}),
	).Run()
}
