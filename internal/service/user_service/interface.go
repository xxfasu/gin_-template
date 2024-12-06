package user_service

import (
	"context"
	"gin_template/internal/model"
	"gin_template/internal/validation"
)

type UserService interface {
	Register(ctx context.Context, req *validation.Register) error
	Login(ctx context.Context, req *validation.Login) (string, error)
	FindUser(ctx context.Context, req *validation.FindUser) (*model.User, error)
}
