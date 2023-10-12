package data_test

import (
	"context"
	"github.com/chff7cb/swissbank/core"
	"github.com/google/uuid"
	"testing"
)

func TestAccountsData_CreateAccount(t *testing.T) {
	ctx := context.Background()
	account := core.Account{
		AccountID:      uuid.NewString(),
		DocumentNumber: "12312312312",
	}

	repos := AccountsData{}
	err := repos.CreateAccount(ctx, &account)
	if err != nil {
		t.Fatal(err)
	}
}
