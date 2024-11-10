package user_domain

type Resp struct {
	Id            uint    `json:"id"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Nickname      string  `json:"nickname"`
	Abstract      string  `json:"abstract"`
	Avatar        string  `json:"avatar"`
	RecallMessage *string `json:"recallMessage"`
	FriendOnline  bool    `json:"friendOnline"`
	Sound         bool    `json:"sound"`
	SecureLink    bool    `json:"secureLink"`
	SavePwd       bool    `json:"savePwd"`
	SearchUser    int8    `json:"searchUser"`
	Verification  int8    `json:"verification"`
	Problem1      string  `json:"problem1"`
	Problem2      string  `json:"problem2"`
	Problem3      string  `json:"problem3"`
	Answer1       string  `json:"answer1"`
	Answer2       string  `json:"answer2"`
	Answer3       string  `json:"answer3"`
	Online        bool    `json:"online"`
}
