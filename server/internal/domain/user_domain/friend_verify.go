package user_domain

// FriendVerify 好友验证表
type FriendVerify struct {
	Id                 uint
	CreateTime         int64
	UpdateTime         int64
	SendUserID         uint   // 发起验证方
	SendUserModel      User   // 发起验证方
	RevUserID          uint   // 接受验证方
	RevUserModel       User   // 接受验证方
	Status             int8   // 状态 0 未操作 1 同意 2 拒绝 3 忽略
	AdditionalMessages string // 附加消息
	// 验证问题  为3和4的时候需要
	Problem1 string `json:"problem1"`
	Problem2 string `json:"problem2"`
	Problem3 string `json:"problem3"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
}
