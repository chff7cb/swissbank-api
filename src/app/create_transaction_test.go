package app_test

import (
	"context"
	"github.com/chff7cb/swissbank/app"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type transactionsUseCaseTestSuite struct {
	suite.Suite
	accountsService     *mocks.AccountsService
	transactionsService *mocks.TransactionsService
	useCase             app.TransactionsUseCase
	transactionData     core.Transaction
}

// TestTransactionsUseCaseSuite entrypoint for running our test suite
func TestTransactionsUseCaseSuite(t *testing.T) {
	suite.Run(t, new(transactionsUseCaseTestSuite))
}

// SetupTest setup data state before running each test case
func (s *transactionsUseCaseTestSuite) SetupTest() {
	s.accountsService = mocks.NewAccountsService(s.T())
	s.transactionsService = mocks.NewTransactionsService(s.T())
	s.useCase = app.NewTransactionsUseCase(s.transactionsService, s.accountsService)
	s.transactionData = core.Transaction{
		TransactionID:   uuid.NewString(),
		AccountID:       uuid.NewString(),
		Amount:          0.1 + 0.2,
		OperationTypeID: 4,
		EventTimestamp:  time.Now(),
	}
}

// TestCreateTransaction test regular use case of transaction creation
func (s *transactionsUseCaseTestSuite) TestCreateTransaction() {
	s.accountsService.On("GetAccountByID", mock.Anything, s.transactionData.AccountID).
		Return(nil, nil)

	s.transactionsService.On("CreateTransaction", mock.Anything, &s.transactionData).
		Return(&s.transactionData, nil)

	newTransaction, err := s.useCase.CreateTransaction(context.Background(), &s.transactionData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), &s.transactionData, newTransaction)
}

// TestCreateTransaction2 test case when creating a transaction with an invalid AccountID
func (s *transactionsUseCaseTestSuite) TestCreateTransaction2() {
	s.accountsService.On("GetAccountByID", mock.Anything, s.transactionData.AccountID).
		Return(nil, core.ErrInvalidAccountID)

	_, err := s.useCase.CreateTransaction(context.Background(), &s.transactionData)

	assert.Equal(s.T(), core.ErrInvalidAccountID, err)
}
