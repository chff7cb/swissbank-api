package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chff7cb/swissbank/app"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/data"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func setupConfig() *viper.Viper {
	viper.AutomaticEnv()
	return viper.GetViper()
}

func setupLogging() (*zap.Logger, error) {
	return zap.NewProduction()
}

func setupGinServer() *gin.Engine {
	return gin.Default()
}

func setupRoutes(r *gin.Engine, accountsHandler *svc.AccountsHandler) gin.IRoutes {
	return r.Group("/v1").
		POST("/accounts", accountsHandler.CreateAccount)
}

func setupDatabase(cfg *viper.Viper) *dynamodb.DynamoDB {
	sdkSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           cfg.GetString("SWISSBANK_AWS_PROFILE"),
	}))

	return dynamodb.New(sdkSession)
}

func main() {
	fx.New(
		// here we will wire in our application layers
		fx.Provide(
			setupConfig,
			setupLogging,
			setupDatabase,
			// data layer
			data.NewDynamoAccountsData,
			// domain service layer
			core.NewAccountsService,
			// use cases layer
			app.NewAccountsUseCase,
			// service handlers
			svc.NewAccountsHandler,
			// build our API endpoint structure
			setupGinServer,
			setupRoutes,
		),
		fx.Invoke(func(r *gin.Engine, _ gin.IRoutes) error {
			return r.Run()
		}),
	).Run()
}
