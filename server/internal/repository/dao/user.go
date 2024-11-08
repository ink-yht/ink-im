package dao

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
	Id uint `gorm:"primaryKey,autoIncrement"`
	// 创建时间
	Ctime int64
	// 更新时间
	Utime    int64
	Email    sql.NullString `gorm:"unique" json:"email"`
	Password string         `gorm:"size:128" json:"password"`
	// 唯一索引冲突 改成 sql.NullString
	// 唯一索引允许有多个空值，但不能有多个 ""
	Phone    sql.NullString `gorm:"unique"`
	Nickname string         `gorm:"size:32" json:"nickname"`
	Abstract string         `gorm:"size:128" json:"abstract"` //	 简介
	Avatar   string         `gorm:"size:256" json:"avatar"`
	IP       string         `gorm:"size:32" json:"ip"`
	Addr     string         `gorm:"size:64" json:"addr"`
	Role     int8           `gorm:"size:8" json:"role"`     // 角色 1 管理员 2 普通用户
	OpenID   string         `gorm:"size:128" json:"OpenID"` // 第三方平台登录的凭证
}

type UserDao interface {
	Insert(ctx context.Context, u UserModel) error
	FindByEmail(ctx context.Context, email string) (UserModel, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &GormUserDAO{db: db}
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (UserModel, error) {
	var user UserModel
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) Insert(ctx context.Context, u UserModel) error {
	// 写入数据库

	// 毫秒
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error

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
