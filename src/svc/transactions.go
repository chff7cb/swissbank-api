package svc

import (
	"github.com/chff7cb/swissbank/app"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionsHandler struct {
	uc app.TransactionsUseCase
}

func NewTransactionsHandler(uc app.TransactionsUseCase) *TransactionsHandler {
	return &TransactionsHandler{uc}
}

func (h *TransactionsHandler) CreateTransaction(ctx *gin.Context) {
	transactionFormData := app.CreateTransactionForm{}

	if err := ctx.ShouldBindJSON(&transactionFormData); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.uc.CreateTransaction(ctx, &transactionFormData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, response)
}
