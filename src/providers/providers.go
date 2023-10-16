package providers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/data"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	// ConfigKeyAccountsTableName is the account information table name
	ConfigKeyAccountsTableName = "SWISSBANK_ACCOUNTS_TABLE_NAME"
	// ConfigKeyTransactionsTableName is the transaction table name
	ConfigKeyTransactionsTableName = "SWISSBANK_TRANSACTIONS_TABLE_NAME"
	// ConfigKeyAWSProfile AWS profile to be passed down to the AWS SDK
	ConfigKeyAWSProfile = "SWISSBANK_AWS_PROFILE"
	// ConfigKeyAWSRegion the AWS region for DynamoDB tables
	ConfigKeyAWSRegion = "SWISSBANK_AWS_REGION"
	// ConfigKeyDynamoDBEndpoint endpoint URL used when running DynamoDB locally
	ConfigKeyDynamoDBEndpoint = "SWISSBANK_DYNAMODB_ENDPOINT_URL"
	// ConfigKeyHTTPListAddress listen address for the main API
	ConfigKeyHTTPListAddress = "SWISSBANK_HTTP_LISTEN_ADDRESS"

	defaultHTTPReadHeaderTimeout = 15
)

func ViperConfigProvider() *viper.Viper {
	cfg := viper.New()
	cfg.AutomaticEnv()

	// initialize default values
	cfg.SetDefault(ConfigKeyAccountsTableName, "Accounts")
	cfg.SetDefault(ConfigKeyTransactionsTableName, "Transactions")
	cfg.SetDefault(ConfigKeyHTTPListAddress, ":8182")
	cfg.SetDefault("HTTP_READ_HEADER_TIMEOUT", defaultHTTPReadHeaderTimeout)

	return cfg
}

func GinProvider(_ *viper.Viper) *gin.Engine {
	return gin.Default()
}

func DynamoDBProvider(cfg *viper.Viper) dynamodbiface.DynamoDBAPI {
	sdkSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(cfg.GetString(ConfigKeyDynamoDBEndpoint)),
			Region:   aws.String(cfg.GetString(ConfigKeyAWSRegion)),
		},
		SharedConfigState: session.SharedConfigEnable,
		Profile:           cfg.GetString(ConfigKeyAWSProfile),
	}))

	return dynamodb.New(sdkSession)
}

func AccountsDataProvider(driver dynamodbiface.DynamoDBAPI, cfg *viper.Viper) core.AccountsDataProxy {
	return data.NewAccountsData(driver, cfg.GetString(ConfigKeyAccountsTableName))
}

func TransactionsDataProvider(driver dynamodbiface.DynamoDBAPI, cfg *viper.Viper) core.TransactionsDataProxy {
	return data.NewTransactionsData(driver, cfg.GetString(ConfigKeyTransactionsTableName))
}
