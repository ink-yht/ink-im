package user_dao

import (
	"context"
	"gorm.io/gorm"
)

// FriendModel 好友表
type FriendModel struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SendUserID     uint   `gorm:"index" json:"sendUserID"`
	RevUserID      uint   `gorm:"index" json:"revUserID"`
	SendUserRemark string `gorm:"size:128" json:"sendUserRemark"` // 当前用户给好友的备注
	RevUserRemark  string `gorm:"size:128" json:"revUserRemark"`  // 好友给当前用户的备注
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
}

// FriendWithNotice 自定义结构体用于存储查询结果
type FriendWithNotice struct {
	Nickname      string `json:"nickname"`
	Abstract      string `json:"abstract"`
	Avatar        string `json:"avatar"`
	FriendModelID uint   `json:"friend_model_id"`
	Notice        string `json:"notice"`
}

type FriendDao interface {
	FindById(ctx context.Context, id uint) ([]FriendWithNotice, error)
}

type GormFriendDAO struct {
	db *gorm.DB
}

func NewFriendDAO(db *gorm.DB) FriendDao {
	return &GormFriendDAO{db: db}
}

func (dao *GormFriendDAO) FindById(ctx context.Context, id uint) ([]FriendWithNotice, error) {
	friends, err := dao.GetUserFriends(ctx, id)
	return friends, err
}

// GetUserFriends 查询函数
func (dao *GormFriendDAO) GetUserFriends(ctx context.Context, userID uint) ([]FriendWithNotice, error) {
	var friendWithNotices []FriendWithNotice

	// 开始事务
	tx := dao.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 使用子查询获取好友的信息及对应的Notice和FriendModel的id
	if err := tx.WithContext(ctx).
		Table("user_models").
		Select("user_models.id, user_models.nickname, user_models.abstract, user_models.avatar, friend_models.id AS friend_model_id, "+
			"CASE "+
			"WHEN friend_models.send_user_id = ? THEN friend_models.send_user_remark "+
			"ELSE friend_models.rev_user_remark "+
			"END AS notice", userID).
		Joins("INNER JOIN friend_models ON "+
			"(friend_models.send_user_id = ? AND friend_models.rev_user_id = user_models.id) "+
			"OR (friend_models.send_user_id = user_models.id AND friend_models.rev_user_id = ?)", userID, userID).
		Scan(&friendWithNotices).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return friendWithNotices, nil
}
