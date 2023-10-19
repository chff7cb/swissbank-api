package svc

import (
	"errors"
	"net/http"

	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/providers"
	"github.com/gin-gonic/gin"
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

// AccountsHandler implements HTTP handler methods for accounts
type AccountsHandler struct {
	service         core.AccountsService
	wrapperProvider providers.GinWrapperProvider
}

// NewAccountsHandler instantiates a new handler of account related requests
func NewAccountsHandler(service core.AccountsService, wrapperProvider providers.GinWrapperProvider) *AccountsHandler {
	return &AccountsHandler{service, wrapperProvider}
}

// CreateAccount handles a request for creating an account
// @Summary      Create an account
// @Description  Create an account for a customer identified by a document number
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account_data body CreateAccountForm true "Information for the new account"
// @Success      200  {object}  AccountResponse
// @Failure      400  {string}  "Invalid account information provided"
// @Failure      500  {string}  "The server could not complete the request due to an internal error"
// @Router       /v1/accounts [post]
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

// GetAccountByID handles a request for account information
// @Summary      Retrieve data of a given account
// @Description  Retrieve data of an account identified by its AccountID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account_id path string true "ID of the account"
// @Success      200  {object}  AccountResponse
// @Failure      404  {string}  "No account with the given AccountID could be found"
// @Failure      500  {string}  "The server could not complete the request due to an internal error"
// @Router       /v1/accounts/{account_id} [get]
func (h *AccountsHandler) GetAccountByID(ctx *gin.Context) {
	ginWrapper := h.wrapperProvider.Wrap(ctx)

	accountID := ginWrapper.Param("account_id")

	if accountID == "" {
		ginWrapper.JSON(http.StatusBadRequest, "account_id cannot be empty")
		return
	}

	accountData, err := h.service.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, core.ErrInvalidAccountID) {
			ginWrapper.JSON(http.StatusNotFound, err.Error())
		} else {
			ginWrapper.JSON(http.StatusBadRequest, err.Error())
		}

		return
	}

	ginWrapper.JSON(http.StatusOK, &AccountResponse{
		AccountID:      accountData.AccountID,
		DocumentNumber: accountData.DocumentNumber,
	})
}
