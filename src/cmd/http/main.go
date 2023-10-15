package main

import (
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/providers"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func setupRoutes(r *gin.Engine,
	accountsHandler *svc.AccountsHandler,
	transactionsHandler *svc.TransactionsHandler) gin.IRoutes {
	v1 := r.Group("/v1")
	v1.
		POST("/accounts", accountsHandler.CreateAccount).
		GET("/accounts/:account_id", accountsHandler.GetAccountByID)
	v1.
		POST("/transactions", transactionsHandler.CreateTransaction)
	return v1
}

func main() {
	fx.New(
		// here we will wire up application layers
		fx.Provide(
			// load base configuration
			providers.ViperConfigProvider,
			// setup logging infrastructure
			providers.LoggingProvider,
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
			// request service handlers
			svc.NewAccountsHandler,
			svc.NewTransactionsHandler,
			// base HTTP server
			providers.GinProvider,
			// setup endpoint structure
			setupRoutes,
		),
		fx.Invoke(func(r *gin.Engine, _ gin.IRoutes, cfg *viper.Viper) error {
			var addrs []string

			if cfg.IsSet("SWISSBANK_HTTP_ADDRESS") {
				addrs = append(addrs, cfg.GetString("SWISSBANK_HTTP_ADDRESS"))
			}

			return r.Run(addrs...)
		}),
	).Run()
}
