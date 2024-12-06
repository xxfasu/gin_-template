package user_repository

import (
	"context"
	"gin_template/internal/data/service_data"
	"gin_template/internal/model"
	"gin_template/internal/repository/gen"
)

type Reader interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByCondition(ctx context.Context, condition service_data.Condition) (*model.User, error)
}

type Writer interface {
	Create(ctx context.Context, user *model.User) error
	CreateTx(ctx context.Context, query *gen.Query, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

type UserRepository interface {
	Reader
	Writer
}
