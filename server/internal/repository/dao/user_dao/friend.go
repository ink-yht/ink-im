package user_dao

// FriendModel 好友表
type FriendModel struct {
	Id            uint `gorm:"primaryKey,autoIncrement"`
	CreateTime    int64
	UpdateTime    int64
	SendUserID    uint      `json:"sendUserID"`                     // 发起验证方
	SendUserModel UserModel `gorm:"foreignKey:SendUserID" json:"-"` // 发起验证方
	RevUserID     uint      `json:"revUserID"`                      // 接受验证方
	RevUserModel  UserModel `gorm:"foreignKey:RevUserID" json:"-"`  // 接受验证方
	Notice        string    `gorm:"size:128" json:"notice"`         // 备注
}
