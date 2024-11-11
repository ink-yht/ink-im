package user_dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicate      = errors.New("账号冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserModel struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreateTime    int64          `json:"createTime"`
	UpdateTime    int64          `json:"updateTime"`
	Email         sql.NullString `gorm:"unique" json:"email"`
	Password      string         `gorm:"size:128" json:"password"`
	Phone         sql.NullString `gorm:"unique" json:"phone"`
	Nickname      string         `gorm:"size:32" json:"nickname"`
	Abstract      string         `gorm:"size:128" json:"abstract"`
	Avatar        string         `gorm:"size:256" json:"avatar"`
	IP            string         `gorm:"size:32" json:"ip"`
	Addr          string         `gorm:"size:64" json:"addr"`
	Role          int8           `gorm:"size:8" json:"role"`
	OpenID        string         `gorm:"size:128" json:"openID"`
	UserConfModel *UserConfModel `gorm:"foreignKey:UserID" json:"userConfModel"`
	Friends       []FriendModel  `gorm:"foreignKey:SendUserID" json:"friends"` // 定义好友关系
}
type UserDao interface {
	Insert(ctx context.Context, u UserModel) error
	FindByEmail(ctx context.Context, email string) (UserModel, error)
	FindById(ctx context.Context, uid uint) (UserModel, error)
	UpdateNonZeroFields(ctx context.Context, u UserModel) error
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &GormUserDAO{db: db}
}

func (dao *GormUserDAO) UpdateNonZeroFields(ctx context.Context, u UserModel) error {

	tx := dao.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&UserModel{}).Where("id =?", u.ID).Updates(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新用户配置表
	userConf := u.UserConfModel
	if userConf != nil {
		if err := tx.Model(&UserConfModel{}).Where("user_id =?", u.ID).Updates(userConf).Error; err != nil {
			tx.Rollback()
			return err
		}

	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil

}

func (dao *GormUserDAO) FindById(ctx context.Context, uid uint) (UserModel, error) {
	var user UserModel
	err := dao.db.WithContext(ctx).Preload("UserConfModel").Take(&user, uid).Error
	return user, err
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (UserModel, error) {
	var user UserModel
	err := dao.db.WithContext(ctx).Preload("UserConfModel").Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) Insert(ctx context.Context, u UserModel) error {
	// 写入数据库

	avatarPath := "./logo.png"

	// 毫秒
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	u.Avatar = avatarPath

	err := dao.db.WithContext(ctx).Preload("UserConfModel").Create(&u).Error

	// 如果错误是MySQL错误类型
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 邮箱冲突
			return ErrDuplicate
		}
	}

	// 系统错误
	return err
}
