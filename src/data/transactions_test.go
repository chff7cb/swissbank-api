package data_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/data"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	errPutTransactionItem = errors.New("putitem error")
)

type transactionsDataTestSuite struct {
	suite.Suite
	driver          *mocks.DynamoDBAPI
	proxy           core.TransactionsDataProxy
	transactionData core.Transaction
}

// TestTransactionsDataSuite entry function that will run our test suite
func TestTransactionsDataSuite(t *testing.T) {
	suite.Run(t, new(transactionsDataTestSuite))
}

// SetupTest setup test state data before running each test case
func (s *transactionsDataTestSuite) SetupTest() {
	s.driver = mocks.NewDynamoDBAPI(s.T())
	s.proxy = data.NewTransactionsData(s.driver, "transactions")
	s.transactionData = core.Transaction{
		TransactionID:   uuid.NewString(),
		AccountID:       uuid.NewString(),
		Amount:          0.1 + 0.2,
		OperationTypeID: 4,
		EventTimestamp:  time.Now(),
	}
}

// TestCreateTransaction test the creating of a new transaction and assert output data values
func (s *transactionsDataTestSuite) TestCreateTransaction() {
	ctx := context.Background()

	s.driver.On("PutItem", mock.Anything).
		Return(func(req *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
			assert.NotEqual(s.T(), nil, req.Item["TransactionID"].S)
			assert.NotEqual(s.T(), nil, req.Item["AccountID"].S)
			assert.NotEqual(s.T(), nil, req.Item["OperationTypeID"].N)
			assert.NotEqual(s.T(), nil, req.Item["Amount"].N)
			assert.NotEqual(s.T(), nil, req.Item["EventTimestamp"].N)

			assert.Equal(s.T(), s.transactionData.TransactionID, *req.Item["TransactionID"].S)
			assert.Equal(s.T(), s.transactionData.AccountID, *req.Item["AccountID"].S)

			operationTypeUnwrapped, err := strconv.ParseInt(*req.Item["OperationTypeID"].N, 10, 64)
			assert.Equal(s.T(), nil, err)
			assert.Equal(s.T(), int64(s.transactionData.OperationTypeID), operationTypeUnwrapped)

			eventTimestampUnwrapped, err := strconv.ParseInt(*req.Item["EventTimestamp"].N, 10, 64)
			assert.Equal(s.T(), nil, err)
			assert.Equal(s.T(), s.transactionData.EventTimestamp.Unix(), eventTimestampUnwrapped)

			amountUnwrapped, err := strconv.ParseFloat(*req.Item["Amount"].N, 64)
			assert.Equal(s.T(), nil, err)
			assert.Equal(s.T(), s.transactionData.Amount, amountUnwrapped)

			return &dynamodb.PutItemOutput{}, nil
		})

	newTransaction, err := s.proxy.CreateTransaction(ctx, &s.transactionData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.transactionData.TransactionID, newTransaction.TransactionID)
	assert.Equal(s.T(), s.transactionData.AccountID, newTransaction.AccountID)
	assert.Equal(s.T(), s.transactionData.Amount, newTransaction.Amount)
	assert.Equal(s.T(), s.transactionData.OperationTypeID, newTransaction.OperationTypeID)
	assert.Equal(s.T(), s.transactionData.EventTimestamp, newTransaction.EventTimestamp)
}

// TestCreateTransaction2 test the failed creation of a transaction due to a database error
func (s *transactionsDataTestSuite) TestCreateTransaction2() {
	ctx := context.Background()

	s.driver.On("PutItem", mock.Anything).
		Return(nil, errPutTransactionItem)

	_, err := s.proxy.CreateTransaction(ctx, &s.transactionData)

	assert.Equal(s.T(), errPutTransactionItem, err)
}
