package data

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chff7cb/swissbank/core"
)

// accountsDataImpl implements operations for storing accounts' data on DynamoDB
type accountsDataImpl struct {
	driver        dynamodbiface.DynamoDBAPI
	accountsTable string
}

// NewAccountsData instantiates a new data proxy for account information
func NewAccountsData(driver dynamodbiface.DynamoDBAPI, accountsTable string) core.AccountsDataProxy {
	return &accountsDataImpl{
		driver:        driver,
		accountsTable: accountsTable,
	}
}

// CreateAccount creates a new account using DynamoDB
func (d *accountsDataImpl) CreateAccount(_ context.Context, account *core.Account) (*core.Account, error) {
	if _, err := d.driver.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.accountsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"AccountID":      {S: aws.String(account.AccountID)},
			"DocumentNumber": {S: aws.String(account.DocumentNumber)},
		},
	}); err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByID lookup account data information by a given AccountID
func (d *accountsDataImpl) GetAccountByID(_ context.Context, accountID string) (itemResult *core.Account, err error) {
	// retrieve an account item matching the given accountID
	result, err := d.driver.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.accountsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"AccountID": {S: aws.String(accountID)},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, core.ErrInvalidAccountID
	}

	// build the output model from result attributes
	return &core.Account{
		AccountID:      *result.Item["AccountID"].S,
		DocumentNumber: *result.Item["DocumentNumber"].S,
	}, nil
}
