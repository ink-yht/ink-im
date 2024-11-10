package user_domain

import "time"

type User struct {
	Id         uint
	CreateTime time.Time
	Email      string
	Password   string
	Phone      string
	Nickname   string
	Abstract   string // 简介
	Avatar     string
	IP         string
	Addr       string
	Role       int8   // 角色 1 管理员 2 普通用户
	OpenID     string // 第三方平台登录的凭证
	UserConf   *UserConf
}
