package user_repo

import (
	"context"
	"ink-im-server/internal/domain/user_domain"
	"ink-im-server/internal/repository/cache/user_cache"
	"ink-im-server/internal/repository/dao/user_dao"
)

type FriendRepository interface {
	FindById(ctx context.Context, uid uint) ([]user_domain.FriendsInfo, error)
}

type CacheFriendRepository struct {
	dao   user_dao.FriendDao
	cache user_cache.FriendCache
}

func NewFriendRepository(dao user_dao.FriendDao, cache user_cache.FriendCache) FriendRepository {
	return &CacheFriendRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CacheFriendRepository) FindById(ctx context.Context, uid uint) ([]user_domain.FriendsInfo, error) {
	//u, err := repo.cache.Get(ctx, uid)
	//// 缓存里面有数据
	//// 缓存里面没有数据
	//// 缓存出错了，不知道有没有数据
	//if err == nil {
	//	// 必然有数据
	//	return u, err
	//}

	users, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return []user_domain.FriendsInfo{}, err
	}

	domainFriends := make([]user_domain.FriendsInfo, len(users))
	for i, us := range users {
		domainFriends[i] = repo.entityToDomain(us)
	}

	return domainFriends, err
}

func (repo *CacheFriendRepository) domainToEntity(u user_domain.FriendsInfo) user_dao.FriendWithNotice {
	data := user_dao.FriendWithNotice{
		FriendModelID: u.FriendModelID,
		Nickname:      u.Nickname,
		Abstract:      u.Abstract,
		Avatar:        u.Avatar,
		Notice:        u.Notice,
	}
	return data
}

func (repo *CacheFriendRepository) entityToDomain(u user_dao.FriendWithNotice) user_domain.FriendsInfo {
	data := user_domain.FriendsInfo{
		FriendModelID: u.FriendModelID,
		Nickname:      u.Nickname,
		Abstract:      u.Abstract,
		Avatar:        u.Avatar,
		Notice:        u.Notice,
	}

	return data
}
