package user_web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"ink-im-server/internal/domain"
	"ink-im-server/internal/service/user_service"
	"ink-im-server/internal/web"
	"ink-im-server/pkg/logger"
	"net/http"
)

const (
	emailRegexPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	passwordRegexPattern = "^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[._~!@#$^&*])[A-Za-z0-9._~!@#$^&*]{8,20}$"
)

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
	ug := server.Group("/users")
	ug.POST("signup", u.Signup)
	ug.POST("login", u.Login)
	ug.POST("edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
		Nickname        string `json:"nickname"`
		Abstract        string `json:"abstract"` // 简介
		Avatar          string `json:"avatar"`
		IP              string `json:"ip"`
		Addr            string `json:"addr"`
		Role            int8   `json:"role"`   // 角色 1 管理员 2 普通用户
		OpenID          string `json:"OpenID"` // 第三方平台登录的凭证
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

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Abstract: req.Abstract,
		Avatar:   req.Avatar,
		IP:       req.IP,
		Addr:     req.Addr,
		Role:     req.Role,
		OpenID:   req.OpenID,
	})

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

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
