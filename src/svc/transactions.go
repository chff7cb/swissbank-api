package svc

import (
	"net/http"
	"time"

	"github.com/chff7cb/swissbank/app"

	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/providers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTransactionForm input data necessary required to register a new transaction
type CreateTransactionForm struct {
	AccountID       string  `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

// TransactionResponse data output model of existing transactions
type TransactionResponse struct {
	TransactionID   string    `json:"transaction_id"`
	AccountID       string    `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventTimestamp  time.Time `json:"event_timestamp"`
	Description     string    `json:"description"`
}

// TransactionsHandler handler class of transaction requests
type TransactionsHandler struct {
	useCase         app.TransactionsUseCase
	wrapperProvider providers.GinWrapperProvider
}

// NewTransactionsHandler instantiates a new handler of transaction requests
func NewTransactionsHandler(useCase app.TransactionsUseCase, provider providers.GinWrapperProvider) *TransactionsHandler {
	return &TransactionsHandler{useCase, provider}
}

// CreateTransaction handles a request for creating a new transaction
// @Summary      Create a transaction
// @Description  Registers a new transaction associated with a given account
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transactionInfo body CreateTransactionForm true "Required data of the new transaction"
// @Success      200  {object} TransactionResponse
// @Router       /transactions [post]
func (h *TransactionsHandler) CreateTransaction(ctx *gin.Context) {
	ginWrapper := h.wrapperProvider.Wrap(ctx)

	transactionForm := CreateTransactionForm{}

	if err := ginWrapper.ShouldBindJSON(&transactionForm); err != nil {
		ginWrapper.JSON(http.StatusBadRequest, err.Error())
		return
	}

	transactionData := core.Transaction{
		TransactionID:   uuid.NewString(),
		AccountID:       transactionForm.AccountID,
		Amount:          transactionForm.Amount,
		OperationTypeID: core.OperationTypeID(transactionForm.OperationTypeID),
		EventTimestamp:  time.Now(),
	}

	newTransaction, err := h.useCase.CreateTransaction(ctx, &transactionData)
	if err != nil {
		ginWrapper.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ginWrapper.JSON(http.StatusOK, &TransactionResponse{
		TransactionID:   newTransaction.TransactionID,
		AccountID:       newTransaction.AccountID,
		OperationTypeID: int(newTransaction.OperationTypeID),
		Amount:          newTransaction.Amount,
		EventTimestamp:  newTransaction.EventTimestamp,
		Description:     newTransaction.GetDescription(),
	})
}
