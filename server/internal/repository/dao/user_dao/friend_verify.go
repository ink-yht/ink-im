package user_dao

// FriendVerifyModel 好友验证表
type FriendVerifyModel struct {
	Id                 uint `gorm:"primaryKey,autoIncrement"`
	CreateTime         int64
	UpdateTime         int64
	SendUserID         uint      `json:"sendUserID"`                         // 发起验证方
	SendUserModel      UserModel `gorm:"foreignKey:SendUserID" json:"-"`     // 发起验证方
	RevUserID          uint      `json:"revUserID"`                          // 接受验证方
	RevUserModel       UserModel `gorm:"foreignKey:RevUserID" json:"-"`      // 接受验证方
	Status             int8      `json:"status"`                             // 状态 0 未操作 1 同意 2 拒绝 3 忽略
	AdditionalMessages string    `gorm:"size:128" json:"additionalMessages"` // 附加消息
	// 验证问题  为3和4的时候需要V // 验证问题  为3和4的时候需要
	Problem1 string `json:"problem1"`
	Problem2 string `json:"problem2"`
	Problem3 string `json:"problem3"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
}
