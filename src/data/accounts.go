package data

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chff7cb/swissbank/core"
)

// AccountsData implements operations for storing accounts' data on DynamoDB
type AccountsData struct {
	driver        dynamodbiface.DynamoDBAPI
	accountsTable string
}

func NewAccountsData(driver dynamodbiface.DynamoDBAPI, accountsTable string) core.AccountsDataProxy {
	return &AccountsData{
		driver:        driver,
		accountsTable: accountsTable,
	}
}

func (d *AccountsData) CreateAccount(_ context.Context, account *core.Account) (itemResult *core.Account, err error) {
	if _, err := d.driver.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.accountsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"AccountID": {
				S: aws.String(account.AccountID),
			},
			"DocumentNumber": {
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

func (d *AccountsData) GetAccountByID(_ context.Context, accountID string) (itemResult *core.Account, err error) {
	// retrieve an account item matching the given accountID
	result, err := d.driver.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.accountsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"AccountID": {
				S: aws.String(accountID),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// build the output model from result attributes
	return &core.Account{
		AccountID:      *result.Item["AccountID"].S,
		DocumentNumber: *result.Item["DocumentNumber"].S,
	}, nil
}
