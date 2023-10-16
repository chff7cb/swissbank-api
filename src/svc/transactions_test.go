package svc_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type transactionsHandlerTestSuite struct {
	suite.Suite
	service         *mocks.TransactionsService
	ginWrapper      *mocks.GinWrapper
	wrapperProvider *mocks.GinWrapperProvider
	handler         *svc.TransactionsHandler
	ginContext      *gin.Context
	transactionData svc.TransactionResponse
}

func TestTransactionsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	suite.Run(t, new(transactionsHandlerTestSuite))
}

func (s *transactionsHandlerTestSuite) SetupTest() {
	s.service = mocks.NewTransactionsService(s.T())
	s.ginWrapper = mocks.NewGinWrapper(s.T())
	s.wrapperProvider = mocks.NewGinWrapperProvider(s.T())
	s.handler = svc.NewTransactionsHandler(s.service, s.wrapperProvider)
	s.ginContext, _ = gin.CreateTestContext(httptest.NewRecorder())
	s.transactionData = svc.TransactionResponse{
		TransactionID:   uuid.NewString(),
		AccountID:       uuid.NewString(),
		Amount:          0.1 + 0.2,
		OperationTypeID: 4,
		EventTimestamp:  time.Now(),
		Description:     "",
	}
}

func (s *transactionsHandlerTestSuite) TestTransactionsHandler_CreateTransaction() {
	s.wrapperProvider.On("Wrap", mock.Anything).Return(s.ginWrapper)

	s.ginWrapper.On("ShouldBindJSON", mock.Anything).
		Return(func(v any) error {
			form, isTransactionForm := v.(*svc.CreateTransactionForm)
			assert.Equal(s.T(), true, isTransactionForm)
			form.AccountID = s.transactionData.AccountID
			form.Amount = s.transactionData.Amount
			form.OperationTypeID = s.transactionData.OperationTypeID
			return nil
		})

	s.service.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(func(_ context.Context, transaction *core.Transaction) (*core.Transaction, error) {
			assert.Equal(s.T(), len(s.transactionData.TransactionID), len(transaction.TransactionID))
			assert.Equal(s.T(), s.transactionData.AccountID, transaction.AccountID)
			assert.Equal(s.T(), s.transactionData.Amount, transaction.Amount)
			assert.Equal(s.T(), core.OperationTypeID(s.transactionData.OperationTypeID), transaction.OperationTypeID)
			assert.GreaterOrEqual(s.T(), s.transactionData.EventTimestamp.Unix(), transaction.EventTimestamp.Unix())

			s.transactionData.TransactionID = transaction.TransactionID
			s.transactionData.EventTimestamp = transaction.EventTimestamp
			s.transactionData.Description = transaction.GetDescription()

			return transaction, nil
		})

	s.ginWrapper.On("JSON", http.StatusOK, mock.Anything).
		Return(func(_ int, response *svc.TransactionResponse) {
			assert.Equal(s.T(), s.transactionData.TransactionID, response.TransactionID)
			assert.Equal(s.T(), s.transactionData.AccountID, response.AccountID)
			assert.Equal(s.T(), s.transactionData.Amount, response.Amount)
			assert.Equal(s.T(), s.transactionData.EventTimestamp, response.EventTimestamp)
			assert.Equal(s.T(), s.transactionData.OperationTypeID, response.OperationTypeID)
			assert.Equal(s.T(), s.transactionData.Description, response.Description)
		})

	s.handler.CreateTransaction(s.ginContext)
}

func (s *transactionsHandlerTestSuite) TestTransactionsHandler_CreateTransaction2() {
	s.wrapperProvider.On("Wrap", mock.Anything).Return(s.ginWrapper)

	decodingError := errors.New("failed to decode json")

	s.ginWrapper.On("ShouldBindJSON", mock.Anything).
		Return(decodingError)

	s.ginWrapper.On("JSON", http.StatusBadRequest, decodingError.Error())

	s.handler.CreateTransaction(s.ginContext)
}

func (s *transactionsHandlerTestSuite) TestTransactionsHandler_CreateTransaction3() {
	s.wrapperProvider.On("Wrap", mock.Anything).Return(s.ginWrapper)

	domainError := errors.New("validation error")

	s.ginWrapper.On("ShouldBindJSON", mock.Anything).
		Return(func(v any) error {
			form, isTransactionForm := v.(*svc.CreateTransactionForm)
			assert.Equal(s.T(), true, isTransactionForm)
			form.AccountID = s.transactionData.AccountID
			form.Amount = s.transactionData.Amount
			form.OperationTypeID = s.transactionData.OperationTypeID
			return nil
		})

	s.service.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(nil, domainError)

	s.ginWrapper.On("JSON", http.StatusBadRequest, domainError.Error())

	s.handler.CreateTransaction(s.ginContext)
}
