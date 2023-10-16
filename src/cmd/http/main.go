package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chff7cb/swissbank/app"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/docs"
	"github.com/chff7cb/swissbank/providers"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

func setupRoutes(ginEngine *gin.Engine,
	accountsHandler *svc.AccountsHandler,
	transactionsHandler *svc.TransactionsHandler) gin.IRoutes {
	v1Router := ginEngine.Group("/v1")
	v1Router.
		POST("/accounts", accountsHandler.CreateAccount).
		GET("/accounts/:account_id", accountsHandler.GetAccountByID)
	v1Router.
		POST("/transactions", transactionsHandler.CreateTransaction)

	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return v1Router
}

// @title           SwissBank challenge API
// @version         1.0
// @description     This is the SwissBank API for managing Accounts and Transactions.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT License
// @license.url   https://opensource.org/license/mit/

// @host      localhost:8182
// @BasePath  /v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
			srv := &http.Server{
				Addr:              cfg.GetString(providers.ConfigKeyHTTPListenAddress),
				Handler:           ginEngine,
				ReadHeaderTimeout: time.Duration(cfg.GetInt(providers.ConfigKeyReadHeaderTimeout)) * time.Second,
			}

			addrParts := strings.Split(srv.Addr, ":")
			if len(addrParts) > 1 {
				docs.SwaggerInfo.Host = "localhost:" + addrParts[1]
			}

			lcx.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						log.Println("http.Server listening on address", srv.Addr)
						log.Println(srv.ListenAndServe())
						log.Println(shutdownHandler.Shutdown())
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown(ctx)
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
