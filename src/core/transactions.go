package core

import (
	"context"
	"errors"
	"time"
)

type OperationTypeID int

const (
	OperationIDPurchaseIntegral OperationTypeID = iota + 1
	OperationIDPurchaseInstallments
	OperationIDAccountWithdraw
	OperationIDCashPayment
)

type Transaction struct {
	TransactionID   string
	AccountID       string
	Amount          float64
	OperationTypeID OperationTypeID
	EventTimestamp  time.Time
}

func (t *Transaction) GetDescription() string {
	switch t.OperationTypeID {
	case OperationIDPurchaseIntegral:
		return "COMPRA A VISTA"
	case OperationIDPurchaseInstallments:
		return "COMPRA PARCELADA"
	case OperationIDAccountWithdraw:
		return "SAQUE"
	case OperationIDCashPayment:
		return "PAGAMENTO"
	}

	return "TIPO NAO REGISTRADO"
}

type TransactionsDataProxy interface {
	CreateTransaction(context.Context, *Transaction) (*Transaction, error)
}

type TransactionsService interface {
	TransactionsDataProxy
}

type transactionsServiceImpl struct {
	proxy TransactionsDataProxy
}

func NewTransactionsService(proxy TransactionsDataProxy) TransactionsService {
	return &transactionsServiceImpl{proxy}
}

func (s *transactionsServiceImpl) CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	if transaction.AccountID == "" {
		return nil, errors.New("AccountID cannot be empty")
	}

	if transaction.TransactionID == "" {
		return nil, errors.New("TransactionID cannot be empty")
	}

	switch transaction.OperationTypeID {
	case OperationIDPurchaseIntegral,
		OperationIDPurchaseInstallments,
		OperationIDAccountWithdraw,
		OperationIDCashPayment:
	default:
		return nil, errors.New("OperationTypeID is invalid")
	}

	return s.proxy.CreateTransaction(ctx, transaction)
}
