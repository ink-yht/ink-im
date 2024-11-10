package user_domain

import "ink-im-server/internal/domain/commit"

// FriendVerify 好友验证表
type FriendVerify struct {
	Id                   uint
	CreateTime           int64
	UpdateTime           int64
	SendUserID           uint        // 发起验证方
	SendUserModel        User        // 发起验证方
	RevUserID            uint        // 接受验证方
	RevUserModel         User        // 接受验证方
	Status               int8        // 状态 0 未操作 1 同意 2 拒绝 3 忽略
	AdditionalMessages   string      // 附加消息
	VerificationQuestion *commit.VFQ // 验证问题  为3和4的时候需要`
}
