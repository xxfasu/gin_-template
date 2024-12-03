//go:build wireinject
// +build wireinject

package main

import (
	"gin_template/internal/handler"
	"gin_template/internal/middleware"
	"gin_template/internal/repository"
	"gin_template/internal/service"
	"gin_template/pkg/logs"
	"gin_template/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func newWire(*logs.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		middleware.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,
		routes.ProviderSet,
	))
}
