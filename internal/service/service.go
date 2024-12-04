package service

import (
	"gin_template/internal/service/user_service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(user_service.NewUserService)
