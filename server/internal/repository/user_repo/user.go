package user_repo

import (
	"context"
	"database/sql"
	"ink-im-server/internal/domain/user_domain"
	"ink-im-server/internal/repository/dao/user_dao"
	"time"
)

var (
	ErrDuplicate      = user_dao.ErrDuplicate
	ErrRecordNotFound = user_dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user user_domain.User) error
	FindByEmail(ctx context.Context, email string) (user_domain.User, error)
	FindById(ctx context.Context, uid uint) (user_domain.User, error)
	UpdateInfo(ctx context.Context, user user_domain.User) (user_domain.User, error)
}

type CacheUserRepository struct {
	dao user_dao.UserDao
}

func NewUserRepository(dao user_dao.UserDao) UserRepository {
	return &CacheUserRepository{
		dao: dao,
	}
}

func (repo *CacheUserRepository) UpdateInfo(ctx context.Context, user user_domain.User) (user_domain.User, error) {
	u, err := repo.dao.UpdateNonZeroFields(ctx, repo.domainToEntity(user))
	return repo.entityToDomain(u), err
}

func (repo *CacheUserRepository) FindById(ctx context.Context, uid uint) (user_domain.User, error) {
	user, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return user_domain.User{}, err
	}
	return repo.entityToDomain(user), err
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (user_domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return user_domain.User{}, err
	}
	return repo.entityToDomain(user), err
}

func (repo *CacheUserRepository) Create(ctx context.Context, user user_domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}

func (repo *CacheUserRepository) domainToEntity(u user_domain.User) user_dao.UserModel {
	data := user_dao.UserModel{
		Id:         u.Id,
		CreateTime: u.CreateTime.UnixMilli(),
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
		UserConfModel: &user_dao.UserConfModel{
			Id:            u.UserConf.Id,
			CreateTime:    u.CreateTime.UnixMilli(),
			UserID:        u.Id,
			RecallMessage: u.UserConf.RecallMessage,
			FriendOnline:  u.UserConf.FriendOnline,
			Sound:         u.UserConf.Sound,
			SecureLink:    u.UserConf.SecureLink,
			SavePwd:       u.UserConf.SavePwd,
			SearchUser:    u.UserConf.SearchUser,
			Verification:  u.UserConf.Verification,
			Problem1:      u.UserConf.Problem1,
			Problem2:      u.UserConf.Problem2,
			Problem3:      u.UserConf.Problem3,
			Answer1:       u.UserConf.Answer1,
			Answer2:       u.UserConf.Answer2,
			Answer3:       u.UserConf.Answer3,
			Online:        u.UserConf.Online,
		},
	}
	return data
}

func (repo *CacheUserRepository) entityToDomain(u user_dao.UserModel) user_domain.User {
	data := user_domain.User{
		Id:         u.Id,
		CreateTime: time.UnixMilli(u.CreateTime),
		Email:      u.Email.String,
		Phone:      u.Phone.String,
		Password:   u.Password,
		Nickname:   u.Nickname,
		Abstract:   u.Abstract,
		Avatar:     u.Avatar,
		IP:         u.IP,
		Addr:       u.Addr,
		Role:       u.Role,
		OpenID:     u.OpenID,
		UserConf: &user_domain.UserConf{
			Id:            u.UserConfModel.Id,
			CreateTime:    time.UnixMilli(u.UserConfModel.CreateTime),
			UserID:        u.Id,
			RecallMessage: u.UserConfModel.RecallMessage,
			FriendOnline:  u.UserConfModel.FriendOnline,
			Sound:         u.UserConfModel.Sound,
			SecureLink:    u.UserConfModel.SecureLink,
			SavePwd:       u.UserConfModel.SavePwd,
			SearchUser:    u.UserConfModel.SearchUser,
			Verification:  u.UserConfModel.Verification,
			Problem1:      u.UserConfModel.Problem1,
			Problem2:      u.UserConfModel.Problem2,
			Problem3:      u.UserConfModel.Problem3,
			Answer1:       u.UserConfModel.Answer1,
			Answer2:       u.UserConfModel.Answer2,
			Answer3:       u.UserConfModel.Answer3,
			Online:        u.UserConfModel.Online,
		},
	}

	return data
}
