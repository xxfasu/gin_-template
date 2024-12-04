package user_service

import (
	"context"
	"gin_template/internal/model/request"
)

type UserService interface {
	Register(ctx context.Context, req *request.Register) error
	Login(ctx context.Context, req *request.Login) (string, error)
}
