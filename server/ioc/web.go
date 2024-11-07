package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"ink-im-server/internal/web/user_web"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *user_web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.UserRegisterRouters(server)
	return server
}

func InitMiddleWares() []gin.HandlerFunc {
	return []gin.HandlerFunc{

		corsHdl(),

		//middlelware.NewLoginJWTMiddlewareBuilder().
		//	IgnorePaths("/users/signup").
		//	IgnorePaths("/users/login_sms/code/send").
		//	IgnorePaths("/users/login_sms").
		//	IgnorePaths("/users/login").Build(),
		//
		//ratelimit.NewBuilder(redisClient, time.Minute, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		//AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your.com")
		},
		MaxAge: 12 * time.Hour,
	})
}