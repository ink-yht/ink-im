package files_dao

import (
	"context"
	"gorm.io/gorm"
	"ink-im-server/internal/repository/dao/user_dao"
)

type AvatarDao interface {
	UpAvatar(ctx context.Context, id uint, url string) error
}

type GormAvatarDAO struct {
	db *gorm.DB
}

func NewAvatarDAO(db *gorm.DB) AvatarDao {
	return &GormAvatarDAO{db: db}
}

func (dao GormAvatarDAO) UpAvatar(ctx context.Context, id uint, url string) error {
	var user user_dao.UserModel
	return dao.db.WithContext(ctx).Model(&user).Where("id = ?", id).Update("avatar", url).Error
}
