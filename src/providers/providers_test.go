package providers_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chff7cb/swissbank/mocks"
	"github.com/chff7cb/swissbank/providers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestViperConfigProvider(t *testing.T) {
	cfg := providers.ViperConfigProvider()

	// make sure defaults were set
	assert.Equal(t, true, cfg.IsSet(providers.ConfigKeyAccountsTableName))
	assert.Equal(t, true, cfg.IsSet(providers.ConfigKeyTransactionsTableName))
	assert.Equal(t, true, cfg.IsSet(providers.ConfigKeyHTTPListenAddress))
	assert.Equal(t, true, cfg.IsSet(providers.ConfigKeyReadHeaderTimeout))
}

func TestGinProvider(t *testing.T) {
	cfg := providers.ViperConfigProvider()

	gin.SetMode(gin.TestMode)

	engine := providers.GinProvider(cfg)

	// no real configuration is done as of right now
	assert.Equal(t, "/", engine.BasePath())
}

func TestDynamoDBProvider(t *testing.T) {
	cfg := providers.ViperConfigProvider()

	ddb := providers.DynamoDBProvider(cfg)

	result, err := ddb.DescribeEndpoints(&dynamodb.DescribeEndpointsInput{})
	assert.Equal(t, nil, err)
	assert.Greater(t, len(result.Endpoints), 0)
}

func TestAccountsDataProvider(t *testing.T) {
	cfg := providers.ViperConfigProvider()
	ddb := mocks.NewDynamoDBAPI(t)

	proxy := providers.AccountsDataProvider(ddb, cfg)

	assert.NotEqual(t, nil, proxy)
}

func TestTransactionsDataProvider(t *testing.T) {
	cfg := providers.ViperConfigProvider()
	ddb := mocks.NewDynamoDBAPI(t)

	proxy := providers.TransactionsDataProvider(ddb, cfg)

	assert.NotEqual(t, nil, proxy)
}
