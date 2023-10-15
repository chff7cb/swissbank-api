package svc

import (
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/providers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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

// AccountsHandler implements HTTP handler methods for accounts
type AccountsHandler struct {
	service         core.AccountsService
	wrapperProvider providers.GinWrapperProvider
}

func NewAccountsHandler(service core.AccountsService, wrapperProvider providers.GinWrapperProvider) *AccountsHandler {
	return &AccountsHandler{service, wrapperProvider}
}

// CreateAccount handles a request for creating an account
func (h *AccountsHandler) CreateAccount(ctx *gin.Context) {
	ginWrapper := h.wrapperProvider.Wrap(ctx)

	accountForm := CreateAccountForm{}
	if err := ginWrapper.ShouldBindJSON(&accountForm); err != nil {
		ginWrapper.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// build new account using form data
	accountData := core.Account{
		AccountID:      uuid.NewString(),
		DocumentNumber: accountForm.DocumentNumber,
	}

	newAccount, err := h.service.CreateAccount(ctx, &accountData)
	if err != nil {
		ginWrapper.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ginWrapper.JSON(http.StatusOK, AccountResponse{
		AccountID:      newAccount.AccountID,
		DocumentNumber: newAccount.DocumentNumber,
	})
}

func (h *AccountsHandler) GetAccountByID(ctx *gin.Context) {
	ginWrapper := h.wrapperProvider.Wrap(ctx)

	accountID := ginWrapper.Param("account_id")

	if accountID == "" {
		ginWrapper.JSON(http.StatusBadRequest, "account_id cannot be empty")
		return
	}

	accountData, err := h.service.GetAccountByID(ctx, accountID)
	if err != nil {
		ginWrapper.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ginWrapper.JSON(http.StatusOK, &AccountResponse{
		AccountID:      accountData.AccountID,
		DocumentNumber: accountData.DocumentNumber,
	})
}
