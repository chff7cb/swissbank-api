package main

import (
	"context"
	"github.com/chff7cb/swissbank/core"
	"github.com/chff7cb/swissbank/providers"
	"github.com/chff7cb/swissbank/svc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
	"net/http"
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
		fx.Invoke(func(lc fx.Lifecycle, cfg *viper.Viper, r *gin.Engine, _ gin.IRoutes, shutdowner fx.Shutdowner) {
			srv := &http.Server{
				Addr:    cfg.GetString(providers.ConfigKeyHttpListAddress),
				Handler: r,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						log.Println("http.Server listening on address", srv.Addr)
						log.Println(srv.ListenAndServe())
						log.Println(shutdowner.Shutdown())
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown(ctx)
				},
			})
		}),
	).Run()
}
