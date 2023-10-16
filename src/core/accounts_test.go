package core_test

import (
	"context"
	"testing"

	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type accountsServiceTestSuite struct {
	suite.Suite
	proxy       *mocks.AccountsDataProxy
	service     core.AccountsService
	accountData core.Account
}

func TestAccountsServiceSuite(t *testing.T) {
	suite.Run(t, new(accountsServiceTestSuite))
}

// SetupTest setup test state before running each test case
func (s *accountsServiceTestSuite) SetupTest() {
	s.proxy = mocks.NewAccountsDataProxy(s.T())
	s.service = core.NewAccountsService(s.proxy)
	s.accountData = core.Account{
		AccountID:      uuid.NewString(),
		DocumentNumber: "12312312312",
	}
}

// TestAccountsServiceImpl_CreateAccount tests a successful account creation
func (s *accountsServiceTestSuite) TestAccountsServiceImpl_CreateAccount() {
	s.proxy.
		On("CreateAccount", context.Background(), &s.accountData).
		Return(&s.accountData, nil)

	newAccount, err := s.service.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.accountData.AccountID, newAccount.AccountID)
	assert.Equal(s.T(), s.accountData.DocumentNumber, newAccount.DocumentNumber)
}

// TestAccountsServiceImpl_CreateAccount2 validates cases where account creation should fail due to malformed parameters
func (s *accountsServiceTestSuite) TestAccountsServiceImpl_CreateAccount2() {
	// test for too few characters in DocumentNumber
	s.accountData.DocumentNumber = "123"
	_, err := s.service.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), core.ErrInvalidDocumentLength, err)

	// test for too many characters in DocumentNumber
	for i := 0; i < 10; i++ {
		s.accountData.DocumentNumber += s.accountData.DocumentNumber
	}

	_, err = s.service.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), core.ErrInvalidDocumentLength, err)

	// test for empty AccountID
	s.accountData.AccountID = ""
	s.accountData.DocumentNumber = "12312312312"
	_, err = s.service.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), core.ErrInvalidAccountID, err)
}

// TestAccountsServiceImpl_GetAccountByID test for retrieving account data by a given accountID
func (s *accountsServiceTestSuite) TestAccountsServiceImpl_GetAccountByID() {
	s.proxy.
		On("GetAccountByID", context.Background(), s.accountData.AccountID).
		Return(&s.accountData, nil)

	account, err := s.service.GetAccountByID(context.Background(), s.accountData.AccountID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.accountData.DocumentNumber, account.DocumentNumber)
	assert.Equal(s.T(), s.accountData.AccountID, account.AccountID)
}
