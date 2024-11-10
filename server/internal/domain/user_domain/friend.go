package user_domain

// Friend 好友表
type Friend struct {
	Id            uint
	CreateTime    int64
	UpdateTime    int64
	SendUserID    uint   // 发起验证方
	SendUserModel User   // 发起验证方
	RevUserID     uint   // 接受验证方
	RevUserModel  User   // 接受验证方
	Notice        string // 备注
}
