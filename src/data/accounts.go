package data

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chff7cb/swissbank/core"
	"github.com/spf13/viper"
)

// AccountsData implements operations for storing accounts' data on DynamoDB
type AccountsData struct {
	driver        *dynamodb.DynamoDB
	accountsTable string
}

func NewDynamoAccountsData(driver *dynamodb.DynamoDB, cfg *viper.Viper) core.AccountsDataProxy {
	return &AccountsData{
		driver:        driver,
		accountsTable: cfg.GetString("SWISSBANK_ACCOUNTS_TABLE_NAME"),
	}
}

func (d *AccountsData) CreateAccount(_ context.Context, account *core.Account) (itemResult *core.Account, err error) {
	if _, err := d.driver.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.accountsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"account_id": {
				S: aws.String(account.AccountID),
			},
			"document_number": {
				S: aws.String(account.DocumentNumber),
			},
		},
	}); err != nil {
		return nil, err
	}

	// TODO: retrieve data from PutItemResult
	itemResult = account

	return itemResult, nil
}
