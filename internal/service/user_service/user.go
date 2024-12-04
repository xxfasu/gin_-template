package user_service

import (
	"context"
	"gin_template/internal/model"
	"gin_template/internal/model/request"
	"gin_template/internal/repository"
	"gin_template/internal/repository/user_repository"
	"gin_template/pkg/logs"

	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewUserService(
	logger *logs.Logger,
	tm repository.Transaction,
	userRepo user_repository.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
		tx:       tm,
	}
}

type userService struct {
	userRepo user_repository.UserRepository
	logger   *logs.Logger
	tx       repository.Transaction
}

func (s *userService) Register(ctx context.Context, req *request.Register) error {
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if err == nil && user != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	userId := ""
	user = &model.User{
		UserId:   userId,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	// Transaction demo
	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *request.Login) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}
	token, err := s.jwt.GenToken(user.UserId, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}
