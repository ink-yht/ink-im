package user_web

import (
	"fmt"
	"ink-im-server/internal/domain/user_domain"
	"ink-im-server/internal/service/user_service"
	"ink-im-server/internal/web"
	"ink-im-server/pkg/logger"
	"net/http"
	"os"
	"path"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	passwordRegexPattern = "^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[._~!@#$^&*])[A-Za-z0-9._~!@#$^&*]{8,20}$"
)

//// 确保 UserHandler 上实现了 handler 接口
//var _ Handler = &UserHandler{}
//
//// 这个更优雅
//var _ Handler = (*UserHandler)(nil)

type UserHandler struct {
	emailRexExp *regexp.Regexp
	passwordRex *regexp.Regexp
	svc         user_service.UserService
	l           logger.Logger
}

func NewUserHandler(svc user_service.UserService, l logger.Logger) *UserHandler {
	return &UserHandler{
		emailRexExp: regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:         svc,
		l:           l,
	}
}

func (u *UserHandler) UserRegisterRouters(server *gin.Engine) {
	//server.GET("/uploads/:imageType/:imageName", u.AvatarShow)
	ug := server.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/info", u.Info)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
		Phone           string `json:"phone"`
		Nickname        string `json:"nickname"`
		Abstract        string `json:"abstract"` // 简介
		Avatar          string `json:"avatar"`
		IP              string `json:"ip"`
		Addr            string `json:"addr"`
		Role            int8   `json:"role"`   // 角色 1 管理员 2 普通用户
		OpenID          string `json:"OpenID"` // 第三方平台登录的凭证
		UserConf        struct {
			RecallMessage *string `json:"recallMessage"` // 撤回消息的提示内容
			FriendOnline  bool    `json:"friendOnline"`  // 好友上线提醒
			SecureLink    bool    `json:"secureLink"`    // 安全链接
			SavePwd       bool    `json:"savePwd"`       // 保存密码
			SearchUser    int8    `json:"searchUser"`    // 别人查找到你的方式
			Verification  int8    `json:"verification"`  // 好友验证方式
			// 验证问题  为3和4的时候需要
			Problem1 string `json:"problem1"`
			Problem2 string `json:"problem2"`
			Problem3 string `json:"problem3"`
			Answer1  string `json:"answer1"`
			Answer2  string `json:"answer2"`
			Answer3  string `json:"answer3"`
			Online   bool   `json:"online"` // 是否在线
		} `json:"userConf"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := u.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("邮箱校验失败", logger.String("email", err.Error()))
		return
	}

	if !isEmail {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "非法邮箱格式",
			Data: nil,
		})
		u.l.Warn("非法邮箱格式", logger.Field{Key: "err", Value: err})
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "两次输入密码不一致",
			Data: nil,
		})
		u.l.Warn("两次输入密码不一致", logger.Field{Key: "err", Value: err})
		return
	}

	isPassword, err := u.passwordRex.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("密码校验失败", logger.String("password", err.Error()))
		return
	}

	if !isPassword {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "密码必须包含数字、字母、特殊字符，且不少于八位",
			Data: nil,
		})
		u.l.Warn("非法密码格式", logger.Field{Key: "err", Value: err})
		return
	}

	user := user_domain.User{
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Abstract: req.Abstract,
		Avatar:   req.Avatar,
		IP:       req.IP,
		Addr:     req.Addr,
		Role:     req.Role,
		OpenID:   req.OpenID,
		UserConf: &user_domain.UserConf{
			RecallMessage: nil,
			FriendOnline:  false,
			Sound:         true,
			SecureLink:    false,
			SavePwd:       false,
			SearchUser:    2,
			Verification:  2,
			Problem1:      "",
			Problem2:      "",
			Problem3:      "",
			Answer1:       "",
			Answer2:       "",
			Answer3:       "",
			Online:        true,
		},
	}

	err = u.svc.SignUp(ctx, user)

	// err 有两种情况
	// 1.系统错误
	// 2.邮箱已注册
	if err == user_service.ErrDuplicate {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "邮箱已被注册",
			Data: nil,
		})
		u.l.Warn("邮箱已被注册", logger.Field{Key: "err", Value: err})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("注册时系统错误", logger.String("zc", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "注册成功",
		Data: nil,
	})
	u.l.Info("注册成功", logger.String("email", req.Email))
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Abstract string `json:"abstract"` // 简介
		Avatar   string `json:"avatar"`
		IP       string `json:"ip"`
		Addr     string `json:"addr"`
		Role     int8   `json:"role"`   // 角色 1 管理员 2 普通用户
		OpenID   string `json:"OpenID"` // 第三方平台登录的凭证
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := u.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("邮箱校验失败", logger.String("email", err.Error()))
		return
	}

	if !isEmail {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "非法邮箱格式",
			Data: nil,
		})
		u.l.Warn("非法邮箱格式", logger.Field{Key: "err", Value: err})
		return
	}

	isPassword, err := u.passwordRex.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("密码校验失败", logger.String("password", err.Error()))
		return
	}

	if !isPassword {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "密码必须包含数字、字母、特殊字符，且不少于八位",
			Data: nil,
		})
		u.l.Warn("非法密码格式", logger.Field{Key: "err", Value: err})
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == user_service.ErrInvalidUserOrPassword {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "用户不存在或密码不对",
			Data: nil,
		})
		u.l.Warn("用户不存在或密码不对", logger.Field{Key: "err", Value: err})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("登录时系统错误", logger.String("dl", err.Error()))
		return
	}

	if err = u.setJWT(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 2,
			"msg":  "系统异常",
		})
		return
	}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "登录成功",
		Data: nil,
	})
	u.l.Info("登录成功", logger.String("email", req.Email))
}

func (u *UserHandler) Edit(ctx *gin.Context) {

	type Req struct {
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
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, web.Result{Code: 1, Msg: "昵称不能为空"})
		u.l.Warn("昵称不能为空")
		return
	}

	if len(req.Abstract) > 128 {
		ctx.JSON(http.StatusOK, web.Result{Code: 1, Msg: "简介过长"})
		u.l.Warn("简介过长")
		return
	}

	userClaims := ctx.MustGet("claims").(UserClaims)

	data := user_domain.User{
		Id:       userClaims.Uid,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Abstract: req.Abstract,
		Avatar:   req.Avatar,
		UserConf: &user_domain.UserConf{
			RecallMessage: req.RecallMessage,
			FriendOnline:  req.FriendOnline,
			Sound:         req.Sound,
			SecureLink:    req.SecureLink,
			SavePwd:       req.SavePwd,
			SearchUser:    req.SearchUser,
			Verification:  req.Verification,
			Problem1:      req.Problem1,
			Problem2:      req.Problem2,
			Problem3:      req.Problem3,
			Answer1:       req.Answer1,
			Answer2:       req.Answer2,
			Answer3:       req.Answer3,
		},
	}
	fmt.Println("data:", data)
	err := u.svc.Edit(ctx, data)
	if err != nil {
		u.l.Error("修改个人信息失败", logger.String("email", req.Email))
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "修改个人信息失败",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "个人信息修改成功",
		Data: nil,
	})
	u.l.Info("修改个人信息失败", logger.String("email", req.Email))
}

func (u *UserHandler) Info(ctx *gin.Context) {

	userClaims := ctx.MustGet("claims").(UserClaims)

	user, err := u.svc.Info(ctx, userClaims.Uid)
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("找不到 id")
		return
	}

	data := user_domain.Resp{
		Id:            user.Id,
		Email:         user.Email,
		Phone:         user.Phone,
		Nickname:      user.Nickname,
		Abstract:      user.Abstract,
		Avatar:        user.Avatar,
		RecallMessage: user.UserConf.RecallMessage,
		FriendOnline:  user.UserConf.FriendOnline,
		Sound:         user.UserConf.Sound,
		SecureLink:    user.UserConf.SecureLink,
		SavePwd:       user.UserConf.SavePwd,
		SearchUser:    user.UserConf.SearchUser,
		Verification:  user.UserConf.Verification,
		Problem1:      user.UserConf.Problem1,
		Problem2:      user.UserConf.Problem2,
		Problem3:      user.UserConf.Problem3,
		Answer1:       user.UserConf.Answer1,
		Answer2:       user.UserConf.Answer2,
		Answer3:       user.UserConf.Answer3,
		Online:        user.UserConf.Online,
	}
	fmt.Println("avatar:", data.Avatar)

	//// 使用 strings.Split 函数按 '/' 分割字符串
	//parts := strings.Split(data.Avatar, "/")
	//
	//// 输出所有部分，以便查看
	//fmt.Println("All parts:", parts)
	//
	//// 获取第二部分和第三部分
	//// 注意：索引从0开始，所以第二部分是索引1，第三部分是索引2
	//if len(parts) > 3 {
	//	secondPart := parts[2]
	//	thirdPart := parts[3]
	//	show, err := u.AvatarShow(ctx, secondPart, thirdPart)
	//	if err != nil {
	//		return
	//	}
	//	data.Avatar = show
	//} else {
	//	fmt.Println("The address does not have enough parts.")
	//}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "个人信息获取成功",
		Data: data,
	})
	u.l.Info("个人信息获取成功", logger.String("email", user.Email))
}

func (u *UserHandler) setJWT(ctx *gin.Context, uid uint) error {
	// 设置 JWT 登录态
	claims := UserClaims{
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {

		return err
	}

	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *UserHandler) AvatarShow(ctx *gin.Context, imageType, imageName string) (avatar string, err error) {
	filePath := path.Join("uploads", imageType, imageName)
	fmt.Println(filePath)
	bateDate, err := os.ReadFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "图片获取失败",
			Data: nil,
		})
		u.l.Error("图片获取失败", logger.String("image", err.Error()))
		return
	}
	// 直接写入图片数据
	// 设置Content-Type为image/jpeg
	ctx.Header("Content-Type", "image/jpeg")

	// 直接返回图片数据
	ctx.Data(http.StatusOK, "image/jpeg", bateDate)
	return
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明自己要放进去 token 里的数据
	Uid       uint
	UserAgent string
}

var JWTKey = []byte("3vnkm3RPr55524y0uuG2PeEUPAT1t3PI")
