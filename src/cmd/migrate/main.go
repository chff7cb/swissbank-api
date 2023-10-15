package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chff7cb/swissbank/providers"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
)

func migrate(_ context.Context, ddb dynamodbiface.DynamoDBAPI, cfg *viper.Viper) (err error) {
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

	return nil
}

func main() {
	fx.New(
		// here we will wire up application layers
		fx.Provide(
			// load base configuration
			providers.ViperConfigProvider,
			// instantiate database clients
			providers.DynamoDBProvider,
		),
		fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner, ddb dynamodbiface.DynamoDBAPI, cfg *viper.Viper) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					defer func() { _ = shutdowner.Shutdown() }()
					return migrate(ctx, ddb, cfg)
				},
			})
		}),
	).Run()
}
