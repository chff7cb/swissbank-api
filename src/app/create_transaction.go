package app

import (
	"context"

	"github.com/chff7cb/swissbank/core"
)

// TransactionsUseCase interface for creating transactions
type TransactionsUseCase interface {
	CreateTransaction(context.Context, *core.Transaction) (*core.Transaction, error)
}

// TransactionsUseCaseImpl provides methods for creating transaction
type TransactionsUseCaseImpl struct {
	transactions core.TransactionsService
	accounts     core.AccountsService
}

// NewTransactionsUseCase construct a new use case for creating transactions
func NewTransactionsUseCase(transactions core.TransactionsService, accounts core.AccountsService) TransactionsUseCase {
	return &TransactionsUseCaseImpl{transactions, accounts}
}

// CreateTransaction performs logic to create a transaction after validating the associated account
func (uc *TransactionsUseCaseImpl) CreateTransaction(ctx context.Context, transactionData *core.Transaction) (response *core.Transaction, err error) {
	// we may only create transactions for existing and valid accounts
	if _, err = uc.accounts.GetAccountByID(ctx, transactionData.AccountID); err != nil {
		return nil, core.ErrInvalidAccountID
	}

	return uc.transactions.CreateTransaction(ctx, transactionData)
}
