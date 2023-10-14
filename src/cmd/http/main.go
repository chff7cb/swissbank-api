package main

import (
	"github.com/chff7cb/swissbank/app"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/providers"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func setupRoutes(r *gin.Engine, accountsHandler *svc.AccountsHandler) gin.IRoutes {
	return r.Group("/v1").
		POST("/accounts", accountsHandler.CreateAccount).
		GET("/accounts/:account_id", accountsHandler.GetAccountByID)
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
			// application data layer
			providers.AccountsDataProvider,
			// application domain service layer
			core.NewAccountsService,
			// application use cases layer
			app.NewAccountsUseCase,
			// request service handlers
			svc.NewAccountsHandler,
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
