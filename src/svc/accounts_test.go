package svc_test

import (
	"context"
	"errors"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type accountsHandlerTestSuite struct {
	suite.Suite
	service         *mocks.AccountsService
	ginWrapper      *mocks.GinWrapper
	wrapperProvider *mocks.GinWrapperProvider
	handler         *svc.AccountsHandler
	ginContext      *gin.Context
	accountData     svc.AccountResponse
}

func TestAccountsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	suite.Run(t, new(accountsHandlerTestSuite))
}

func (s *accountsHandlerTestSuite) SetupTest() {
	s.service = mocks.NewAccountsService(s.T())
	s.ginWrapper = mocks.NewGinWrapper(s.T())
	s.wrapperProvider = mocks.NewGinWrapperProvider(s.T())
	s.handler = svc.NewAccountsHandler(s.service, s.wrapperProvider)
	s.ginContext, _ = gin.CreateTestContext(httptest.NewRecorder())
	s.accountData = svc.AccountResponse{
		AccountID:      uuid.NewString(),
		DocumentNumber: "12312312312",
	}
}

func (s *accountsHandlerTestSuite) TestAccountsHandler_CreateAccount() {
	s.wrapperProvider.On("Wrap", s.ginContext).Return(s.ginWrapper)

	s.ginWrapper.
		On("ShouldBindJSON", mock.Anything).
		Return(func(v any) error {
			form, isAccountForm := v.(*svc.CreateAccountForm)
			assert.Equal(s.T(), true, isAccountForm)
			form.DocumentNumber = s.accountData.DocumentNumber
			return nil
		})

	// we expect that the service method is called with correct input parameters
	s.service.
		On("CreateAccount", mock.Anything, mock.Anything).
		Return(func(_ context.Context, data *core.Account) (*core.Account, error) {
			// the generated accountID should have the same length
			assert.Equal(s.T(), len(s.accountData.AccountID), len(data.AccountID))
			assert.Equal(s.T(), s.accountData.DocumentNumber, data.DocumentNumber)
			// save the generated accountID so we may check in the final response
			s.accountData.AccountID = data.AccountID
			return data, nil
		})

	s.ginWrapper.
		On("JSON", http.StatusOK, mock.Anything).
		Return(func(_ int, response svc.AccountResponse) {
			assert.Equal(s.T(), s.accountData.AccountID, response.AccountID)
			assert.Equal(s.T(), s.accountData.DocumentNumber, response.DocumentNumber)
		})

	s.handler.CreateAccount(s.ginContext)
}

func (s *accountsHandlerTestSuite) TestAccountsHandler_CreateAccount2() {
	s.wrapperProvider.On("Wrap", s.ginContext).Return(s.ginWrapper)

	bindingError := errors.New("could not parse request")

	// the json decoding should fail
	s.ginWrapper.
		On("ShouldBindJSON", mock.Anything).
		Return(bindingError)

	// an error response is expected
	s.ginWrapper.
		On("JSON", http.StatusBadRequest, bindingError.Error())

	s.handler.CreateAccount(s.ginContext)
}

func (s *accountsHandlerTestSuite) TestAccountsHandler_CreateAccount3() {
	s.wrapperProvider.On("Wrap", s.ginContext).Return(s.ginWrapper)

	domainError := errors.New("validation error")

	s.ginWrapper.
		On("ShouldBindJSON", mock.Anything).
		Return(nil)

	s.service.
		On("CreateAccount", mock.Anything, mock.Anything).
		Return(nil, domainError)

	// an error response is expected
	s.ginWrapper.
		On("JSON", http.StatusBadRequest, domainError.Error())

	s.handler.CreateAccount(s.ginContext)
}

func (s *accountsHandlerTestSuite) TestAccountsHandler_GetAccountByID() {
	s.wrapperProvider.On("Wrap", s.ginContext).Return(s.ginWrapper)

	s.ginWrapper.On("Param", "account_id").
		Return(s.accountData.AccountID)

	s.service.On("GetAccountByID", mock.Anything, s.accountData.AccountID).
		Return(&core.Account{
			AccountID:      s.accountData.AccountID,
			DocumentNumber: s.accountData.DocumentNumber,
		}, nil)

	s.ginWrapper.On("JSON", http.StatusOK, mock.Anything).
		Return(func(_ int, response *svc.AccountResponse) {
			assert.Equal(s.T(), s.accountData.AccountID, response.AccountID)
			assert.Equal(s.T(), s.accountData.DocumentNumber, response.DocumentNumber)
		})

	s.handler.GetAccountByID(s.ginContext)
}

func (s *accountsHandlerTestSuite) TestAccountsHandler_GetAccountByID2() {
	s.wrapperProvider.On("Wrap", s.ginContext).Return(s.ginWrapper)

	// intentionally return an invalid account_id
	s.ginWrapper.On("Param", "account_id").
		Return("")

	s.ginWrapper.On("JSON", http.StatusBadRequest, mock.Anything)

	s.handler.GetAccountByID(s.ginContext)
}

func (s *accountsHandlerTestSuite) TestAccountsHandler_GetAccountByID3() {
	s.wrapperProvider.On("Wrap", s.ginContext).Return(s.ginWrapper)

	s.ginWrapper.On("Param", "account_id").
		Return(s.accountData.AccountID)

	domainError := errors.New("validation error")
	s.service.On("GetAccountByID", mock.Anything, s.accountData.AccountID).
		Return(nil, domainError)

	s.ginWrapper.On("JSON", http.StatusBadRequest, domainError.Error())

	s.handler.GetAccountByID(s.ginContext)
}
