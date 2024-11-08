package user_service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"ink-im-server/internal/domain"
	"ink-im-server/internal/repository/user_repo"
	"ink-im-server/pkg/logger"
)

var (
	ErrDuplicate             = user_repo.ErrDuplicate
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
}

type userService struct {
	repo user_repo.UserRepository
	l    logger.Logger
}

func NewUserService(repo user_repo.UserRepository, l logger.Logger) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	// err 两种情况
	// 1.系统错误
	// 2.用户没找到

	if err == user_repo.ErrRecordNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return domain.User{}, err
	}

	// 密码校验
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (svc *userService) SignUp(ctx context.Context, user domain.User) error {
	// 加密，然后存起来
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		svc.l.Error("加密失败", logger.String("bcrypt", err.Error()))
		return err
	}
	user.Password = string(hash)

	// 存起来
	return svc.repo.Create(ctx, user)
}
