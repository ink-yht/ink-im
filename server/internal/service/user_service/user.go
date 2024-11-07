package user_service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"ink-im-server/internal/domain"
	"ink-im-server/internal/repository/user_repo"
)

var ErrDuplicate = user_repo.ErrDuplicate

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
}

type userService struct {
	repo user_repo.UserRepository
}

func NewUserService(repo user_repo.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, user domain.User) error {
	// 加密，然后存起来
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// 存起来
	return svc.repo.Create(ctx, user)
}
