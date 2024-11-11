package user_domain

// FriendsInfo 自定义结构体用于存储查询结果
type FriendsInfo struct {
	FriendModelID uint   `json:"friend_model_id"`
	Nickname      string `json:"nickname"`
	Abstract      string `json:"abstract"`
	Avatar        string `json:"avatar"`
	Notice        string `json:"notice"`
}
