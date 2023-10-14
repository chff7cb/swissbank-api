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
	"go.uber.org/zap"
)

const (
	ConfigKeyAccountsTableName     = "SWISSBANK_ACCOUNTS_TABLE_NAME"
	ConfigKeyTransactionsTableName = "SWISSBANK_TRANSACTIONS_TABLE_NAME"
	ConfigKeyAWSProfile            = "SWISSBANK_AWS_PROFILE"
	ConfigKeyAWSRegion             = "SWISSBANK_AWS_REGION"
	ConfigKeyDynamoDBEndpoint      = "SWISSBANK_DYNAMODB_ENDPOINT_URL"
)

func ViperConfigProvider() *viper.Viper {
	cfg := viper.New()
	cfg.AutomaticEnv()

	// initialize default values
	cfg.SetDefault(ConfigKeyAccountsTableName, "Accounts")
	cfg.SetDefault(ConfigKeyTransactionsTableName, "Transactions")

	return cfg
}

func LoggingProvider(_ *viper.Viper) (*zap.Logger, error) {
	return zap.NewProduction()
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
