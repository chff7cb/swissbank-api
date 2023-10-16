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
	errGetAccountItem = errors.New("getitem error")
	errPutAccountItem = errors.New("putitem error")
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

// TestCreateAccount will test the creation of a new account and the output data
func (s *accountsDataTestSuite) TestCreateAccount() {
	s.driver.On("PutItem", mock.Anything).
		Return(func(req *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
			assert.NotEqual(s.T(), nil, req.Item["AccountID"].S)
			assert.NotEqual(s.T(), nil, req.Item["DocumentNumber"].S)

			assert.Equal(s.T(), s.accountData.AccountID, *req.Item["AccountID"].S)
			assert.Equal(s.T(), s.accountData.DocumentNumber, *req.Item["DocumentNumber"].S)

			return &dynamodb.PutItemOutput{}, nil
		})

	newAccount, err := s.proxy.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.accountData.AccountID, newAccount.AccountID)
	assert.Equal(s.T(), s.accountData.DocumentNumber, newAccount.DocumentNumber)
}

// TestCreateAccount2 tests the case of a PutItem error when creating a new account
func (s *accountsDataTestSuite) TestCreateAccount2() {
	s.driver.On("PutItem", mock.Anything).
		Return(nil, errPutAccountItem)

	_, err := s.proxy.CreateAccount(context.Background(), &s.accountData)

	assert.Equal(s.T(), errPutAccountItem, err)
}

// TestGetAccountByID will test the returned data of the requests account
func (s *accountsDataTestSuite) TestGetAccountByID() {
	getItemOutput := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"AccountID":      {S: &s.accountData.AccountID},
			"DocumentNumber": {S: &s.accountData.DocumentNumber},
		},
	}

	s.driver.On("GetItem", mock.Anything).
		Return(func(req *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
			assert.NotEqual(s.T(), nil, req.Key["AccountID"].S)
			assert.Equal(s.T(), s.accountData.AccountID, *req.Key["AccountID"].S)
			return &getItemOutput, nil
		}, nil)

	accountData, err := s.proxy.GetAccountByID(context.Background(), s.accountData.AccountID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.accountData.AccountID, accountData.AccountID)
	assert.Equal(s.T(), s.accountData.DocumentNumber, accountData.DocumentNumber)
}

// TestGetAccountByID2 test when retrieving account data fails due to GetItem error
func (s *accountsDataTestSuite) TestGetAccountByID2() {
	s.driver.On("GetItem", mock.Anything).
		Return(nil, errGetAccountItem)

	_, err := s.proxy.GetAccountByID(context.Background(), s.accountData.AccountID)

	assert.Equal(s.T(), errGetAccountItem, err)
}

// TestGetAccountByID3 test when retrieving account data fails because of an invalid AccountID
func (s *accountsDataTestSuite) TestGetAccountByID3() {
	s.driver.On("GetItem", mock.Anything).
		Return(&dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{},
		}, nil)

	_, err := s.proxy.GetAccountByID(context.Background(), s.accountData.AccountID)

	assert.Equal(s.T(), core.ErrInvalidAccountID, err)
}
