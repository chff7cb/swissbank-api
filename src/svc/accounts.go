package svc

import (
	"github.com/chff7cb/swissbank/app"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AccountsHandler implements HTTP handler methods for accounts
type AccountsHandler struct {
	uc app.AccountsUseCase
}

func NewAccountsHandler(uc app.AccountsUseCase) *AccountsHandler {
	return &AccountsHandler{uc}
}

// CreateAccount handles a request for creating an account
func (h *AccountsHandler) CreateAccount(ctx *gin.Context) {
	newAccountData := app.CreateAccountForm{}
	if err := ctx.ShouldBindJSON(&newAccountData); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.uc.CreateAccount(ctx, &newAccountData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	ctx.JSON(http.StatusOK, response)
}
