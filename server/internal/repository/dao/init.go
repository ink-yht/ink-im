package dao

import (
	"gorm.io/gorm"
	"ink-im-server/internal/repository/dao/user_dao"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&user_dao.UserModel{},         // 用户表
		&user_dao.FriendModel{},       // 好友表
		&user_dao.FriendVerifyModel{}, // 好友验证表
		&user_dao.UserConfModel{},     // 用户配置表
	)
}
