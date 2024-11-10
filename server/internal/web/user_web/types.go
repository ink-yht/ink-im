package user_web

import "github.com/gin-gonic/gin"

type Handler interface {
	UserRegisterRouters(server *gin.Engine)
}
