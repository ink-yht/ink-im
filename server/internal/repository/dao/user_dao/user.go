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
	Id         uint `gorm:"primaryKey,autoIncrement"`
	CreateTime int64
	UpdateTime int64
	Email      sql.NullString `gorm:"unique" json:"email"`
	Password   string         `gorm:"size:128" json:"password"`
	// 唯一索引冲突 改成 sql.NullString
	// 唯一索引允许有多个空值，但不能有多个 ""
	Phone         sql.NullString `gorm:"unique"`
	Nickname      string         `gorm:"size:32" json:"nickname"`
	Abstract      string         `gorm:"size:128" json:"abstract"` //	 简介
	Avatar        string         `gorm:"size:256" json:"avatar"`
	IP            string         `gorm:"size:32" json:"ip"`
	Addr          string         `gorm:"size:64" json:"addr"`
	Role          int8           `gorm:"size:8" json:"role"`     // 角色 1 管理员 2 普通用户
	OpenID        string         `gorm:"size:128" json:"OpenID"` // 第三方平台登录的凭证
	UserConfModel *UserConfModel `gorm:"foreignKey:UserID" json:"UserConfModel"`
}

type UserDao interface {
	Insert(ctx context.Context, u UserModel) error
	FindByEmail(ctx context.Context, email string) (UserModel, error)
	FindById(ctx context.Context, uid uint) (UserModel, error)
	UpdateNonZeroFields(ctx context.Context, u UserModel) (UserModel, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &GormUserDAO{db: db}
}

func (dao *GormUserDAO) UpdateNonZeroFields(ctx context.Context, u UserModel) (UserModel, error) {

	tx := dao.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//var u UserModel
	//
	//user := &UserModel{
	//	// 设置要更新的用户模型字段值
	//	Id:       u.Id,
	//	Email:    u.Email,
	//	Phone:    u.Phone,
	//	Nickname: u.Nickname,
	//	Abstract: u.Abstract,
	//	Avatar:   u.Avatar,
	//	IP:       u.IP,
	//	Addr:     u.Addr,
	//	Role:     u.Role,
	//}

	//userConf := &UserConfModel{
	//	// 设置要更新的用户配置模型字段值
	//	RecallMessage: u.UserConfModel.RecallMessage,
	//	FriendOnline:  u.UserConfModel.FriendOnline,
	//	Sound:         u.UserConfModel.Sound,
	//	SecureLink:    u.UserConfModel.SecureLink,
	//	SavePwd:       u.UserConfModel.SavePwd,
	//	SearchUser:    u.UserConfModel.SearchUser,
	//	Verification:  u.UserConfModel.Verification,
	//	Problem1:      u.UserConfModel.Problem1,
	//	Problem2:      u.UserConfModel.Problem2,
	//	Problem3:      u.UserConfModel.Problem3,
	//	Answer1:       u.UserConfModel.Answer1,
	//	Answer2:       u.UserConfModel.Answer2,
	//	Answer3:       u.UserConfModel.Answer3,
	//}
	//var ucf UserConfModel
	//
	//if err := tx.Model(userConf).Where("user_id =?", u.Id).Updates(UserConfModel{}).Error; err != nil {
	//	tx.Rollback()
	//	return ucf, err
	//}
	//
	//return user, err

	if err := tx.Model(&UserModel{}).Where("id =?", u.Id).Updates(u).Error; err != nil {
		tx.Rollback()
		return UserModel{}, err
	}

	// 获取更新后的用户数据
	updatedUser := UserModel{}
	if err := tx.Where("id =?", u.Id).First(&updatedUser).Error; err != nil {
		tx.Rollback()
		return UserModel{}, err
	}

	// 更新用户配置表
	userConf := u.UserConfModel
	if userConf != nil {
		if err := tx.Model(&UserConfModel{}).Where("user_id =?", u.Id).Updates(userConf).Error; err != nil {
			tx.Rollback()
			return UserModel{}, err
		}

		// 获取更新后的用户配置数据
		updatedUserConf := UserConfModel{}
		if err := tx.Where("user_id =?", u.Id).First(&updatedUserConf).Error; err != nil {
			tx.Rollback()
			return UserModel{}, err
		}

		// 将更新后的用户配置关联到用户数据
		updatedUser.UserConfModel = &updatedUserConf
	}

	if err := tx.Commit().Error; err != nil {
		return UserModel{}, err
	}

	return updatedUser, nil

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
