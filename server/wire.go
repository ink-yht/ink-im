//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"ink-im-server/internal/repository/dao/user_dao"
	"ink-im-server/internal/repository/user_repo"
	"ink-im-server/internal/service/user_service"
	"ink-im-server/internal/web/user_web"
	"ink-im-server/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitLogger,

		// DAO 部分
		user_dao.NewUserDAO,

		// cache 部分

		// repository 部分
		user_repo.NewUserRepository,

		// service 部分
		user_service.NewUserService,

		// Handler 部分
		user_web.NewUserHandler,

		// 中间件
		ioc.InitWebServer,
		ioc.InitMiddleWares,
	)
	return new(gin.Engine)
}
