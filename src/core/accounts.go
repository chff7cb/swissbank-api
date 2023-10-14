package core

import (
	"context"
	"net/http"
)

const (
	// DocumentNumberMinLength minimum length of document identifier of an account
	DocumentNumberMinLength = 11
	// DocumentNumberMaxLength maximum length of document identifier of an account
	DocumentNumberMaxLength = 48
)

// Account models an account
type Account struct {
	// AccountID
	AccountID string
	// DocumentNumber is a number that identifies the account's holder
	DocumentNumber string
}

// AccountsDataProxy adapter interface for account data operations
type AccountsDataProxy interface {
	CreateAccount(context.Context, *Account) (*Account, error)
	GetAccountByID(context.Context, string) (*Account, error)
}

// AccountsService domain service interface for managing accounts
type AccountsService interface {
	AccountsDataProxy
}

// AccountsServiceImpl implements domain service for account data management
type AccountsServiceImpl struct {
	proxy AccountsDataProxy
}

func NewAccountsService(proxy AccountsDataProxy) AccountsService {
	return &AccountsServiceImpl{proxy}
}

func (s *AccountsServiceImpl) CreateAccount(ctx context.Context, accountData *Account) (r *Account, err error) {
	if len(accountData.DocumentNumber) < DocumentNumberMinLength {
		return nil, http.ErrAbortHandler
	}
	if len(accountData.DocumentNumber) > DocumentNumberMaxLength {
		return nil, http.ErrAbortHandler
	}
	return s.proxy.CreateAccount(ctx, accountData)
}

func (s *AccountsServiceImpl) GetAccountByID(ctx context.Context, accountID string) (r *Account, err error) {
	return s.proxy.GetAccountByID(ctx, accountID)
}
