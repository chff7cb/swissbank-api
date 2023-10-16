package providers_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"testing"

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

func TestAWSConfigProvider(t *testing.T) {
	cfg := providers.ViperConfigProvider()

	dynamodbEndpointURL := "http://localhost:8000"
	awsRegion := "us-east-1"

	cfg.Set(providers.ConfigKeyDynamoDBEndpoint, dynamodbEndpointURL)
	cfg.Set(providers.ConfigKeyAWSRegion, awsRegion)

	configProvider := providers.AWSConfigProvider(cfg)
	awsConfig := configProvider.ClientConfig(dynamodb.ServiceName)

	assert.Equal(t, dynamodbEndpointURL, awsConfig.Endpoint)
	assert.Equal(t, awsRegion, awsConfig.SigningRegion)
}

type awsConfigProviderImpl struct{}

func (a awsConfigProviderImpl) ClientConfig(string, ...*aws.Config) client.Config {
	return client.Config{Config: &aws.Config{}}
}

func TestDynamoDBProvider(t *testing.T) {
	assert.NotEqual(t, nil, providers.DynamoDBProvider(&awsConfigProviderImpl{}))
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
