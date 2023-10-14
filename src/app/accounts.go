package app

import (
	"context"
	"github.com/chff7cb/swissbank/core"
	"github.com/google/uuid"
)

// CreateAccountForm data required for creating an account
type CreateAccountForm struct {
	DocumentNumber string `json:"document_number"`
}

// AccountResponse response data of account information
type AccountResponse struct {
	AccountID      string `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

// AccountsUseCase interface for account use cases
type AccountsUseCase interface {
	CreateAccount(context.Context, *CreateAccountForm) (*AccountResponse, error)
	GetAccountByID(context.Context, string) (*AccountResponse, error)
}

// AccountsUseCaseImpl provides methods containing application use cases for accounts
type AccountsUseCaseImpl struct {
	service core.AccountsService
}

func NewAccountsUseCase(service core.AccountsService) AccountsUseCase {
	return &AccountsUseCaseImpl{service}
}

// CreateAccount performs logic to create an account
func (uc *AccountsUseCaseImpl) CreateAccount(ctx context.Context, accountData *CreateAccountForm) (response *AccountResponse, err error) {
	// construct new account using form data
	newAccount := core.Account{
		AccountID:      uuid.NewString(),
		DocumentNumber: accountData.DocumentNumber,
	}

	if _, err = uc.service.CreateAccount(ctx, &newAccount); err != nil {
		return nil, err
	}

	// fill in response from the new account
	response = &AccountResponse{
		AccountID:      newAccount.AccountID,
		DocumentNumber: newAccount.DocumentNumber,
	}
	return
}

func (uc *AccountsUseCaseImpl) GetAccountByID(ctx context.Context, accountID string) (*AccountResponse, error) {
	accountData, err := uc.service.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &AccountResponse{
		AccountID:      accountData.AccountID,
		DocumentNumber: accountData.DocumentNumber,
	}, nil
}
