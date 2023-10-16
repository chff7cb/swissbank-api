package data_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/data"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	errPutItem = errors.New("putitem error")
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
		EventTimestamp:  time.Time{},
	}
}

// TestTransactionsDataImpl_CreateTransaction test the creating of a new transaction and assert output data values
func (s *transactionsDataTestSuite) TestTransactionsDataImpl_CreateTransaction() {
	ctx := context.Background()

	s.driver.On("PutItem", mock.Anything).
		Return(nil, nil)

	newTransaction, err := s.proxy.CreateTransaction(ctx, &s.transactionData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.transactionData.TransactionID, newTransaction.TransactionID)
	assert.Equal(s.T(), s.transactionData.AccountID, newTransaction.AccountID)
	assert.Equal(s.T(), s.transactionData.Amount, newTransaction.Amount)
	assert.Equal(s.T(), s.transactionData.OperationTypeID, newTransaction.OperationTypeID)
	assert.Equal(s.T(), s.transactionData.EventTimestamp, newTransaction.EventTimestamp)
}

// TestTransactionsDataImpl_CreateTransaction2 test the failed creation of a transaction due to a database error
func (s *transactionsDataTestSuite) TestTransactionsDataImpl_CreateTransaction2() {
	ctx := context.Background()

	s.driver.On("PutItem", mock.Anything).
		Return(nil, errPutItem)

	_, err := s.proxy.CreateTransaction(ctx, &s.transactionData)

	assert.Equal(s.T(), errPutItem, err)
}
