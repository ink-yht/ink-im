package user_web

import (
	"github.com/gin-gonic/gin"
	"ink-im-server/internal/domain/user_domain"
	"ink-im-server/internal/service/user_service"
	"ink-im-server/internal/web"
	"ink-im-server/pkg/logger"
	"net/http"
)

type FriendHandler struct {
	svc user_service.FriendService
	l   logger.Logger
}

func NewFriendHandler(svc user_service.FriendService, l logger.Logger) *FriendHandler {
	return &FriendHandler{
		svc: svc,
		l:   l,
	}
}

func (f *FriendHandler) FriendRegisterRouters(server *gin.Engine) {
	fg := server.Group("/friends")
	fg.GET("/info", f.Info)
}

func (f *FriendHandler) Info(ctx *gin.Context) {
	userClaims := ctx.MustGet("claims").(UserClaims)
	users, err := f.svc.Info(ctx, userClaims.Uid)
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。

		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}

	var respData []user_domain.FriendsInfo
	for _, user := range users {
		respData = append(respData, user_domain.FriendsInfo{
			FriendModelID: user.FriendModelID,
			Nickname:      user.Nickname,
			Abstract:      user.Abstract,
			Avatar:        user.Avatar,
			Notice:        user.Notice,
		})
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "好友信息获取成功",
		Data: respData,
	})
}
