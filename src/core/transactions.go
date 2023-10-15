package core

import (
	"context"
	"errors"
	"time"
)

// OperationTypeID is an operation type identifier
type OperationTypeID int

const (
	// OperationIDPurchaseIntegral operation type of an integral purchase
	OperationIDPurchaseIntegral OperationTypeID = iota + 1
	// OperationIDPurchaseInstallments operation type of a purchase whose payment is split up in multiple installments
	OperationIDPurchaseInstallments
	// OperationIDAccountWithdraw operation type of value withdrawn from an account
	OperationIDAccountWithdraw
	// OperationIDCashPayment operation type of bill payment
	OperationIDCashPayment
)

var (
	EmptyTransactionIDError   = errors.New("TransactionID cannot be empty")
	EmptyAccountIDError       = errors.New("AccountID cannot be empty")
	InvalidOperationTypeError = errors.New("OperationTypeID is invalid")
)

// Transaction models a transaction associated with a given account
type Transaction struct {
	// TransactionID the unique identifier of this transaction
	TransactionID string
	// AccountID reference to the account associated with this transaction
	AccountID string
	// Amount the value amount of the transaction
	Amount float64
	// OperationTypeID identifies the transaction kind
	OperationTypeID OperationTypeID
	// EventTimestamp the time when this transaction took effect
	EventTimestamp time.Time
}

// GetDescription returns the description of a transaction according to its associated operation type
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

	return "TIPO DESCONHECIDO"
}

// TransactionsDataProxy data proxy interface for Transactions management
type TransactionsDataProxy interface {
	CreateTransaction(context.Context, *Transaction) (*Transaction, error)
}

// TransactionsService domain service for managing Transactions
type TransactionsService interface {
	TransactionsDataProxy
}

// transactionsServiceImpl implements methods for creating and validating Transactions
type transactionsServiceImpl struct {
	proxy TransactionsDataProxy
}

// NewTransactionsService instantiates a new domain service
func NewTransactionsService(proxy TransactionsDataProxy) TransactionsService {
	return &transactionsServiceImpl{proxy}
}

// CreateTransaction creates a transaction after validating associated data
func (s *transactionsServiceImpl) CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	if transaction.AccountID == "" {
		return nil, EmptyAccountIDError
	}

	if transaction.TransactionID == "" {
		return nil, EmptyTransactionIDError
	}

	switch transaction.OperationTypeID {
	case OperationIDPurchaseIntegral,
		OperationIDPurchaseInstallments,
		OperationIDAccountWithdraw,
		OperationIDCashPayment:
	default:
		return nil, InvalidOperationTypeError
	}

	return s.proxy.CreateTransaction(ctx, transaction)
}
