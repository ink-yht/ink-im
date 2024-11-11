package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	middlelware "ink-im-server/internal/web/middleware"
	"ink-im-server/internal/web/user_web"
	"ink-im-server/pkg/logger"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc,

	userHdl *user_web.UserHandler,
	friendHdl *user_web.FriendHandler,
) *gin.Engine {

	server := gin.Default()
	server.Use(mdls...)
	userHdl.UserRegisterRouters(server)
	friendHdl.FriendRegisterRouters(server)
	return server
}

func InitMiddleWares(l logger.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{

		corsHdl(),

		//log.NewMiddlewaresLoggerBuilder(func(ctx context.Context, al *log.AccessLog) {
		//	l.Debug("HTTP请求", logger.Field{Key: "al", Value: al})
		//}).AllowReqBody().AllowRespBody().Build(),

		middlelware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").Build(),
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
