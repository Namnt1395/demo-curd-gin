package main

import (
	"context"
	v1 "demo-curd/api/v1"
	"demo-curd/config"
	"demo-curd/database"
	"demo-curd/i18n"
	"demo-curd/model"
	"demo-curd/rabbitmq"
	"demo-curd/router"
	"demo-curd/util"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/swag/example/basic/docs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Config    config.Config
	Database  *database.Database
	Router    *router.Router
	RabbitMQ  *rabbitmq.RabbitMQ
	I18n      *i18n.I18n
	CurdV1Api *v1.CurdV1Api
}

func (r App) Start() error {
	// setup routers
	r.SetupRouters()

	// migration
	r.Database.DB.AutoMigrate(&model.Curd{})

	// run Gin engine
	util.CheckError(r.Router.Engine.Run(fmt.Sprintf(":%s", r.Config.Server.Port)))

	gracefulShutdown(&http.Server{
		Addr:    fmt.Sprintf(":%s", r.Config.Server.Port),
		Handler: r.Router.Engine,
	})

	return nil
}

func (r App) Stop() {
	if err := r.Database.Close(); err != nil {
		panic(err)
	}
	if err := r.RabbitMQ.Close(); err != nil {
		panic(err)
	}
}

func (r App) SetupRouters() {
	// test group
	// public api v1
	//groupPublicV1 := r.Router.Engine.Group("/api/public/v1")
	//{
	//
	//}
	// authorized api v1
	groupV1 := r.Router.Engine.Group("/api/v1")
	groupV1.Use(r.Router.AuthMiddleware.MiddlewareFunc())
	{
		// foo API
		groupV1.POST("curd", r.CurdV1Api.Create)
	}

	// init swagger
	r.InitSwagger(r.Config)
}

func (r App) InitSwagger(c config.Config) {
	docs.SwaggerInfo.Host = c.Swagger.Url
	r.Router.InitSwagger(c)
}

// @title Sample Service API
// @version 1.0
// @description This is Sample Service API.

// @contact.name Namnt
// @contact.email namnguyenthanh024@gmail.com

// @host localhost:8099
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// wire
	app, err := InitApp()
	util.CheckError(err)

	err = app.Start()
	util.CheckError(err)

	defer app.Stop()

	log.Info().Msg("App started")
}

func gracefulShutdown(srv *http.Server) {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
