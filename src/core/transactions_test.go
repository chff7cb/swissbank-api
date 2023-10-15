package core_test

import (
	"context"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type transactionsServiceTestSuite struct {
	suite.Suite
	proxy           *mocks.TransactionsDataProxy
	service         core.TransactionsService
	transactionData core.Transaction
}

func TestTransactionsServiceSuite(t *testing.T) {
	suite.Run(t, new(transactionsServiceTestSuite))
}

// SetupTest setup test state before running each test case
func (s *transactionsServiceTestSuite) SetupTest() {
	s.proxy = mocks.NewTransactionsDataProxy(s.T())
	s.service = core.NewTransactionsService(s.proxy)
	s.transactionData = core.Transaction{
		TransactionID:   uuid.NewString(),
		AccountID:       uuid.NewString(),
		Amount:          0.1 + 0.2,
		OperationTypeID: 4,
		EventTimestamp:  time.Time{},
	}
}

// TestTransactionsServiceImpl_CreateTransaction test default case for creating a new transaction
func (s *transactionsServiceTestSuite) TestTransactionsServiceImpl_CreateTransaction() {
	s.proxy.
		On("CreateTransaction", context.Background(), &s.transactionData).
		Return(&s.transactionData, nil)

	_, err := s.service.CreateTransaction(context.Background(), &s.transactionData)

	assert.Equal(s.T(), nil, err)
}

// TestTransactionsServiceImpl_CreateTransaction2 test error cases of invalid transaction data input
func (s *transactionsServiceTestSuite) TestTransactionsServiceImpl_CreateTransaction2() {
	transactionID := s.transactionData.TransactionID
	accountID := s.transactionData.AccountID

	s.transactionData.TransactionID = ""
	_, err := s.service.CreateTransaction(context.Background(), &s.transactionData)
	s.transactionData.TransactionID = transactionID
	assert.Equal(s.T(), core.EmptyTransactionIDError, err)

	s.transactionData.AccountID = ""
	_, err = s.service.CreateTransaction(context.Background(), &s.transactionData)
	s.transactionData.AccountID = accountID

	assert.Equal(s.T(), core.EmptyAccountIDError, err)

	s.transactionData.OperationTypeID = core.OperationTypeID(9999)
	_, err = s.service.CreateTransaction(context.Background(), &s.transactionData)

	assert.Equal(s.T(), core.InvalidOperationTypeError, err)
	assert.Equal(s.T(), "TIPO DESCONHECIDO", s.transactionData.GetDescription())

	allOperationTypes := []core.OperationTypeID{
		core.OperationIDPurchaseIntegral,
		core.OperationIDPurchaseInstallments,
		core.OperationIDAccountWithdraw,
		core.OperationIDCashPayment,
	}

	s.proxy.
		On("CreateTransaction", context.Background(), &s.transactionData).
		Return(&s.transactionData, nil)
	for i := range allOperationTypes {
		s.transactionData.OperationTypeID = allOperationTypes[i]
		newTransaction, err := s.service.CreateTransaction(context.Background(), &s.transactionData)
		assert.Equal(s.T(), nil, err)
		assert.NotEqual(s.T(), "TIPO DESCONHECIDO", newTransaction.GetDescription())
	}
}
