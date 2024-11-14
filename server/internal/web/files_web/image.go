package files_web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ink-im-server/internal/service/files_service"
	"ink-im-server/internal/web"
	"ink-im-server/internal/web/user_web"
	"ink-im-server/pkg/logger"
	"net/http"
	"os"
	"path"
)

type AvatarHandler struct {
	svc files_service.AvatarService
	l   logger.Logger
}

func NewAvatarHandler(svc files_service.AvatarService, l logger.Logger) *AvatarHandler {
	return &AvatarHandler{
		svc: svc,
		l:   l,
	}
}

func (i *AvatarHandler) AvatarRegisterRouters(server *gin.Engine) {

	ig := server.Group("/api/files")
	ig.POST("/avatar", i.Avatar)
	ig.GET("/uploads/:imageType/:imageName", i.ImageShow)
}

func (i *AvatarHandler) Avatar(ctx *gin.Context) {

	userClaims := ctx.MustGet("claims").(user_web.UserClaims)
	fmt.Println(userClaims.Uid)
	image, err := i.svc.Avatar(ctx, userClaims.Uid)
	ctx.Header("Content-Type", "image/jpeg")
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "上传失败",
			Data: image,
		})
		i.l.Error("图片上传失败", logger.String("image", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "上传成功",
		Data: image,
	})
}

func (i *AvatarHandler) ImageShow(ctx *gin.Context) {

	// 直接从URL路径获取参数
	imageType := ctx.Param("imageType")
	imageName := ctx.Param("imageName")
	filePath := path.Join("uploads", imageType, imageName)
	fmt.Println(filePath)
	bateDate, err := os.ReadFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "图片获取失败",
			Data: nil,
		})
		i.l.Error("图片获取失败", logger.String("image", err.Error()))
		return
	}
	// 直接写入图片数据
	// 设置Content-Type为image/jpeg
	ctx.Header("Content-Type", "image/jpeg")

	// 直接返回图片数据
	ctx.Data(http.StatusOK, "image/jpeg", bateDate)
}
