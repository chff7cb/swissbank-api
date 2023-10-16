package data

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chff7cb/swissbank/core"
)

type transactionsDataImpl struct {
	driver            dynamodbiface.DynamoDBAPI
	transactionsTable string
}

func NewTransactionsData(driver dynamodbiface.DynamoDBAPI, transactionsTableName string) core.TransactionsDataProxy {
	return &transactionsDataImpl{driver, transactionsTableName}
}

func (d *transactionsDataImpl) CreateTransaction(_ context.Context, transaction *core.Transaction) (*core.Transaction, error) {
	if _, err := d.driver.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.transactionsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"TransactionID":   {S: aws.String(transaction.TransactionID)},
			"AccountID":       {S: aws.String(transaction.AccountID)},
			"OperationTypeID": {N: aws.String(strconv.FormatInt(int64(transaction.OperationTypeID), 10))},
			"Amount":          {N: aws.String(strconv.FormatFloat(transaction.Amount, 'f', -1, 64))},
			"EventTimestamp":  {N: aws.String(strconv.FormatInt(transaction.EventTimestamp.Unix(), 10))},
		},
	}); err != nil {
		return nil, err
	}

	return transaction, nil
}
