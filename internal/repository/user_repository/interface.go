package user_repository

import (
	"context"
	"gin_template/internal/model"
)

type Reader interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type Writer interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

type UserRepository interface {
	Reader
	Writer
}
