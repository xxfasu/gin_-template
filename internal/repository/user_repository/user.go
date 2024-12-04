package user_repository

import (
	"context"
	"gin_template/internal/model"
	"gin_template/internal/repository/gen"
	"gorm.io/gorm"
)

func NewUserRepository(
	db *gorm.DB,
) UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	err := gen.User.WithContext(ctx).Create(user)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	err := gen.User.WithContext(ctx).Save(user)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, userId string) (*model.User, error) {
	user, err := gen.User.WithContext(ctx).Where(gen.User.UserID.Eq(userId)).First()
	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := gen.User.WithContext(ctx).Where(gen.User.Email.Eq(email)).First()
	return user, err
}
