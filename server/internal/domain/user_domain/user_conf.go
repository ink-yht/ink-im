package user_domain

import (
	"time"
)

type UserConf struct {
	Id            uint
	CreateTime    time.Time
	UserID        uint
	UserModel     User
	RecallMessage *string // 撤回消息的提示内容
	FriendOnline  bool    // 好友上线提醒
	Sound         bool    // 声音
	SecureLink    bool    // 安全链接
	SavePwd       bool    // 保存密码
	SearchUser    int8    // 别人查找到你的方式 0 不允许别人查找到我， 1  通过用户号找到我 2 可以通过昵称搜索到我
	Verification  int8    // 好友验证 0 不允许任何人添加  1 允许任何人添加  2 需要验证消息 3 需要回答问题  4  需要正确回答问题
	// 验证问题  为3和4的时候需要
	Problem1 string `json:"problem1"`
	Problem2 string `json:"problem2"`
	Problem3 string `json:"problem3"`
	Answer1  string `json:"answer1"`
	Answer2  string `json:"answer2"`
	Answer3  string `json:"answer3"`
	Online   bool   // 是否在线
}
