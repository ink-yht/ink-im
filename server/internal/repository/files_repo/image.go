package files_repo

import (
	"context"
	"ink-im-server/internal/repository/cache/files_cache"
	"ink-im-server/internal/repository/dao/files_dao"
)

type AvatarRepository interface {
	UpAvatar(ctx context.Context, id uint, url string) error
}

type CacheAvatarRepository struct {
	dao   files_dao.AvatarDao
	cache files_cache.AvatarCache
}

func NewAvatarRepository(dao files_dao.AvatarDao, cache files_cache.AvatarCache) AvatarRepository {
	return &CacheAvatarRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo CacheAvatarRepository) UpAvatar(ctx context.Context, id uint, url string) error {
	return repo.dao.UpAvatar(ctx, id, url)
}
