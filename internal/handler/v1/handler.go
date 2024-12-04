package v1

import (
	"gin_template/internal/handler/v1/user_handler"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(user_handler.NewUserHandler)
