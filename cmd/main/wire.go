//go:build wireinject
// +build wireinject

package main

import (
	"gin_template/internal/handler/v1"
	"gin_template/internal/middleware"
	"gin_template/internal/repository"
	"gin_template/internal/service"
	"gin_template/pkg/cache"
	"gin_template/pkg/jwt"
	"gin_template/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func newWire(client *redis.Client, rLock *redsync.Redsync) (*gin.Engine, func(), error) {
	panic(wire.Build(
		middleware.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		v1.ProviderSet,
		routes.ProviderSet,
		jwt.NewJwt,
		cache.InitLocalCache,
	))
}
