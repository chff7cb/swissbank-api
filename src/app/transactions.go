package app

import (
	"context"
	"github.com/chff7cb/swissbank/core"
	"github.com/google/uuid"
	"time"
)

type CreateTransactionForm struct {
	AccountID       string  `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

type TransactionResponse struct {
	TransactonID    string    `json:"transaction_id"`
	AccountID       string    `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventTimestamp  time.Time `json:"event_timestamp"`
	Description     string    `json:"description"`
}

type TransactionsUseCase interface {
	CreateTransaction(ctx context.Context, form *CreateTransactionForm) (*TransactionResponse, error)
}

type transactionsUseCaseImpl struct {
	service core.TransactionsService
}

func NewTransactionsUseCase(service core.TransactionsService) TransactionsUseCase {
	return &transactionsUseCaseImpl{service}
}

func (uc *transactionsUseCaseImpl) CreateTransaction(ctx context.Context, form *CreateTransactionForm) (*TransactionResponse, error) {
	transactionData := core.Transaction{
		TransactionID:   uuid.NewString(),
		AccountID:       form.AccountID,
		Amount:          form.Amount,
		OperationTypeID: core.OperationTypeID(form.OperationTypeID),
		EventTimestamp:  time.Now(),
	}

	newTransaction, err := uc.service.CreateTransaction(ctx, &transactionData)
	if err != nil {
		return nil, err
	}

	return &TransactionResponse{
		TransactonID:    newTransaction.TransactionID,
		AccountID:       newTransaction.AccountID,
		OperationTypeID: int(newTransaction.OperationTypeID),
		Amount:          newTransaction.Amount,
		EventTimestamp:  newTransaction.EventTimestamp,
		Description:     newTransaction.GetDescription(),
	}, nil
}
