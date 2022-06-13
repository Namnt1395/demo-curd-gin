//+build wireinject

package main

import (
	v1 "demo-curd/api/v1"
	"demo-curd/config"
	"demo-curd/dao"
	"demo-curd/database"
	"demo-curd/i18n"
	"demo-curd/rabbitmq"
	"demo-curd/router"
	"demo-curd/service"
	"github.com/google/wire"
)

func InitApp() (App, error) {
	panic(wire.Build(
		// infrastructure
		config.LoadConfig,
		database.NewDatabase,
		i18n.NewI18n,
		rabbitmq.NewRabbitMQ,
		router.NewRouterWithoutAuthMw,
		// dao
		wire.Struct(new(dao.CurdDao), "*"),
		//service
		wire.Struct(new(service.CurdService), "*"),
		// api
		wire.Struct(new(v1.CurdV1Api), "*"),
	// app
	wire.Struct(new(App), "*")))
	return App{}, nil
}
