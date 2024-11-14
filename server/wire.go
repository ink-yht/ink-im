//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"ink-im-server/internal/repository/cache/files_cache"
	"ink-im-server/internal/repository/cache/user_cache"
	"ink-im-server/internal/repository/dao/files_dao"
	"ink-im-server/internal/repository/dao/user_dao"
	"ink-im-server/internal/repository/files_repo"
	"ink-im-server/internal/repository/user_repo"
	"ink-im-server/internal/service/files_service"
	"ink-im-server/internal/service/user_service"
	"ink-im-server/internal/web/files_web"
	"ink-im-server/internal/web/user_web"
	"ink-im-server/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis, ioc.InitLogger,

		// DAO 部分
		user_dao.NewUserDAO,
		user_dao.NewFriendDAO,
		files_dao.NewAvatarDAO,

		// cache 部分
		user_cache.NewUserCache,
		user_cache.NewFriendCache,
		files_cache.NewAvatarCache,

		// repository 部分
		user_repo.NewUserRepository,
		user_repo.NewFriendRepository,
		files_repo.NewAvatarRepository,

		// service 部分
		user_service.NewUserService,
		user_service.NewFriendService,
		files_service.NewAvatarService,

		// Handler 部分
		user_web.NewUserHandler,
		user_web.NewFriendHandler,
		files_web.NewAvatarHandler,

		// 中间件
		ioc.InitWebServer,
		ioc.InitMiddleWares,
	)
	return new(gin.Engine)
}
