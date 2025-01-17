// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gin_template/internal/handler/v1/user_handler"
	"gin_template/internal/middleware"
	"gin_template/internal/repository"
	"gin_template/internal/repository/user_repository"
	"gin_template/internal/service/user_service"
	"gin_template/pkg/cache"
	"gin_template/pkg/jwt"
	"gin_template/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

// Injectors from wire.go:

func newWire(client *redis.Client, rLock *redsync.Redsync) (*gin.Engine, func(), error) {
	recovery := middleware.NewRecoveryM()
	cors := middleware.NewCorsM()
	logM := middleware.NewLogM()
	jwtJWT := jwt.NewJwt()
	authM := middleware.NewAuthM(jwtJWT)
	cacheCache := cache.InitLocalCache(client)
	db, cleanup, err := repository.InitDB()
	if err != nil {
		return nil, nil, err
	}
	transaction := repository.NewTransaction(db)
	userRepository := user_repository.NewUserRepository(db)
	userService := user_service.NewUserService(cacheCache, rLock, transaction, userRepository, jwtJWT)
	userHandler := user_handler.NewUserHandler(userService)
	engine := routes.NewRouter(recovery, cors, logM, authM, userHandler)
	return engine, func() {
		cleanup()
	}, nil
}
