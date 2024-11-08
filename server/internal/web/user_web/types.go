package user_web

import "github.com/gin-gonic/gin"

type handler interface {
	UserRegisterRouters(server *gin.Engine)
}
