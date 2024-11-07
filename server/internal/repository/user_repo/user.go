package user_repo

import (
	"context"
	"database/sql"
	"ink-im-server/internal/domain"
	"ink-im-server/internal/repository/dao"
	"time"
)

var ErrDuplicate = dao.ErrDuplicate

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type CacheUserRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &CacheUserRepository{
		dao: dao,
	}
}

func (repo *CacheUserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}

func (repo *CacheUserRepository) domainToEntity(u domain.User) dao.UserModel {
	return dao.UserModel{
		Id:    u.Id,
		Ctime: u.Ctime.UnixMilli(),
		Email: sql.NullString{
			String: u.Email,
			// 确实有手机号
			Valid: u.Email != "",
		},

		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Nickname: u.Nickname,
		Abstract: u.Abstract,
		Avatar:   u.Avatar,
		IP:       u.IP,
		Addr:     u.Addr,
		Role:     u.Role,
		OpenID:   u.OpenID,
	}
}

func (repo *CacheUserRepository) entityToDomain(u dao.UserModel) domain.User {
	return domain.User{
		Id:       u.Id,
		Ctime:    time.UnixMilli(u.Ctime),
		Email:    u.Email.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Abstract: u.Abstract,
		Avatar:   u.Avatar,
		IP:       u.IP,
		Addr:     u.Addr,
		Role:     u.Role,
		OpenID:   u.OpenID,
	}
}
