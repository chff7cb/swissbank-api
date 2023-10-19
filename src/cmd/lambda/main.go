package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/chff7cb/swissbank/app"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/providers"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func setupRoutes(ginEngine *gin.Engine, cfg *viper.Viper,
	accountsHandler *svc.AccountsHandler,
	transactionsHandler *svc.TransactionsHandler) gin.IRoutes {
	baseRouter := ginEngine.Group(cfg.GetString("API_BASE_PATH"))
	v1Router := baseRouter.Group("/v1")

	v1Router.
		POST("/accounts", accountsHandler.CreateAccount).
		GET("/accounts/:account_id", accountsHandler.GetAccountByID)
	v1Router.
		POST("/transactions", transactionsHandler.CreateTransaction)

	return v1Router
}

func main() {
	fx.New(
		// here we will wire up application layers
		fx.Provide(
			// load base configuration
			providers.ViperConfigProvider,
			// setup AWS client from our base configuration
			providers.AWSConfigProvider,
			// instantiate database clients
			providers.DynamoDBProvider,
			// wrapper interface for gin methods
			providers.NewGinWrapperProvider,
			// application data layer
			providers.AccountsDataProvider,
			providers.TransactionsDataProvider,
			// application domain service layer
			core.NewAccountsService,
			core.NewTransactionsService,
			// use cases setup
			app.NewTransactionsUseCase,
			// request service handlers
			svc.NewAccountsHandler,
			svc.NewTransactionsHandler,
			// base HTTP server
			providers.GinProvider,
			// setup endpoint structure
			setupRoutes,
		),
		fx.Invoke(func(lcx fx.Lifecycle, cfg *viper.Viper, ginEngine *gin.Engine, _ gin.IRoutes, shutdownHandler fx.Shutdowner) {
			lcx.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ginLambda := ginadapter.NewV2(ginEngine)
					go func() {
						lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
							return ginLambda.ProxyWithContext(ctx, req)
						})
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
		fx.Invoke(func(lc fx.Lifecycle, shutdownHandler fx.Shutdowner, ddb dynamodbiface.DynamoDBAPI, cfg *viper.Viper) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return providers.CreateTables(ctx, ddb, cfg)
				},
			})
		}),
	).Run()
}
