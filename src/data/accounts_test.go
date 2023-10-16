package data_test

import (
	"context"
	"errors"
	"testing"

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
	errGetItem      = errors.New("getitem error")
	errPutItemError = errors.New("putitem error")
)

// accountsDataTestSuite test suite for accounts data proxy
type accountsDataTestSuite struct {
	suite.Suite
	driver      *mocks.DynamoDBAPI
	proxy       core.AccountsDataProxy
	accountData core.Account
}

// TestAccountsDataSuite entry function that will run our test suite
func TestAccountsDataSuite(t *testing.T) {
	suite.Run(t, new(accountsDataTestSuite))
}

// SetupTest setup test state data before running each test case
func (s *accountsDataTestSuite) SetupTest() {
	s.driver = mocks.NewDynamoDBAPI(s.T())
	s.proxy = data.NewAccountsData(s.driver, "accounts")
	s.accountData = core.Account{
		AccountID:      uuid.NewString(),
		DocumentNumber: "12312312312",
	}
}

// TestNewAccountsData tests the instantiation of a new accounts data proxy
func (s *accountsDataTestSuite) TestNewAccountsData() {
	proxy := data.NewAccountsData(s.driver, "accounts")
	assert.NotEqual(s.T(), nil, proxy)
}

// TestAccountsDataImpl_CreateAccount will test the creation of a new account and the output data
func (s *accountsDataTestSuite) TestAccountsDataImpl_CreateAccount() {
	s.driver.On("PutItem", mock.Anything).Return(nil, nil)

	newAccount, err := s.proxy.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.accountData.AccountID, newAccount.AccountID)
	assert.Equal(s.T(), s.accountData.DocumentNumber, newAccount.DocumentNumber)
}

// TestAccountsDataImpl_CreateAccount2 tests the case of a PutItem error when creating a new account
func (s *accountsDataTestSuite) TestAccountsDataImpl_CreateAccount2() {
	s.driver.On("PutItem", mock.Anything).Return(nil, errPutItemError)

	_, err := s.proxy.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), errPutItemError, err)
}

// TestAccountsDataImpl_GetAccountByID will test the returned data of the requests account
func (s *accountsDataTestSuite) TestAccountsDataImpl_GetAccountByID() {
	getItemOutput := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"AccountID":      {S: &s.accountData.AccountID},
			"DocumentNumber": {S: &s.accountData.DocumentNumber},
		},
	}

	s.driver.On("GetItem", mock.Anything).
		Return(&getItemOutput, nil)

	accountData, err := s.proxy.GetAccountByID(context.Background(), s.accountData.AccountID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.accountData.AccountID, accountData.AccountID)
	assert.Equal(s.T(), s.accountData.DocumentNumber, accountData.DocumentNumber)
}

// TestAccountsDataImpl_GetAccountByID2 test when retrieving account data fails due to GetItem error
func (s *accountsDataTestSuite) TestAccountsDataImpl_GetAccountByID2() {
	s.driver.On("GetItem", mock.Anything).
		Return(nil, errGetItem)

	_, err := s.proxy.GetAccountByID(context.Background(), s.accountData.AccountID)

	assert.Equal(s.T(), errGetItem, err)
}
