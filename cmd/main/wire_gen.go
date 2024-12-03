// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gin_template/internal/handler"
	"gin_template/internal/middleware"
	"gin_template/internal/repository"
	"gin_template/internal/service"
	"gin_template/pkg/logs"
	"gin_template/routes"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func newWire(logger *logs.Logger) (*gin.Engine, func(), error) {
	db := repository.NewDB(logger)
	repositoryRepository := repository.NewRepository(logger, db)
	userRepository := repository.NewUserRepository(repositoryRepository)
	transaction := repository.NewTransaction(repositoryRepository)
	userService := service.NewUserService(logger, userRepository, transaction)
	userHandler := handler.NewUserHandler(logger, userService)
	recovery := middleware.NewRecoveryM(logger)
	cors := middleware.NewCorsM()
	logM := middleware.NewLogM(logger)
	engine := routes.NewRouter(userHandler, recovery, cors, logM)
	return engine, func() {
	}, nil
}