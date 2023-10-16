package providers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/chff7cb/swissbank/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errDynamoResourceInUseException = errors.New(dynamodb.ErrCodeResourceInUseException)

func TestCreateTables(t *testing.T) {
	ctx := context.Background()
	ddb := mocks.NewDynamoDBAPI(t)
	cfg := providers.ViperConfigProvider()

	ddb.On("CreateTable", mock.Anything).
		Return(nil, nil)
	ddb.
		On("CreateTable", mock.Anything).
		Return(nil, nil)

	assert.Equal(t, nil, providers.CreateTables(ctx, ddb, cfg))
}

func TestCreateTables2(t *testing.T) {
	ctx := context.Background()
	ddb := mocks.NewDynamoDBAPI(t)
	cfg := providers.ViperConfigProvider()

	ddb.On("CreateTable", mock.Anything).
		Return(nil, errDynamoResourceInUseException)
	ddb.
		On("CreateTable", mock.Anything).
		Return(nil, errDynamoResourceInUseException)

	assert.Equal(t, nil, providers.CreateTables(ctx, ddb, cfg))
}
