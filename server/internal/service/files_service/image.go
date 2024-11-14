package files_service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"ink-im-server/internal/domain/file_domain"
	"ink-im-server/internal/repository/files_repo"
	"ink-im-server/pkg/logger"
	"os"
	"path"
)

type AvatarService interface {
	Avatar(ctx *gin.Context, id uint) ([]file_domain.FileResponse, error)
}

type avatarService struct {
	repo files_repo.AvatarRepository
	l    logger.Logger
}

func NewAvatarService(repo files_repo.AvatarRepository, l logger.Logger) AvatarService {
	return &avatarService{
		repo: repo,
		l:    l,
	}
}

func (svc avatarService) Avatar(ctx *gin.Context, id uint) ([]file_domain.FileResponse, error) {

	//var WhiteImageList = []string{
	//	"jpg", "png", "jpeg", "ioc", "gif", "svg", "webp",
	//}

	type Config struct {
		Size int    `yaml:"size"`
		Path string `yaml:"path"`
	}
	var c Config
	err := viper.UnmarshalKey("uploads", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败: %s \n", err))
	}

	imageType := ctx.Request.FormValue("imageType")
	if imageType == "" {
		return nil, errors.New("imageType 不能为空")
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		svc.l.Error("上传图片错误")
		return nil, errors.New("上传图片错误")
	}
	fileList, ok := form.File["image"]
	if !ok {
		return nil, errors.New("上传图片错误")
	}

	// 判断图片路径是否存在，没有就创建
	basePath := c.Path + imageType
	_, err = os.ReadDir(basePath)
	if err != nil {
		err = os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			svc.l.Error("图片路径创建失败", logger.String("err", err.Error()))
			return nil, errors.New("图片路径创建失败")
		}
	}

	var resList []file_domain.FileResponse
	for _, image := range fileList {
		//nameList := strings.Split(image.Filename, ".")
		//suffix := strings.ToLower(nameList[len(nameList)-1])
		filePath := path.Join(c.Path, imageType, image.Filename)
		// 判断大小
		size := float64(image.Size) / float64(1024*1024)
		if size >= float64(c.Size) {
			resList = append(resList, file_domain.FileResponse{
				Filename:  image.Filename,
				IsSuccess: false,
				Msg:       fmt.Sprintf("图片大小超过设定大小，当前大小为: %.2fMb,设定大小为: %dMb ", size, c.Size),
			})
			continue
		}

		err := ctx.SaveUploadedFile(image, filePath)
		if err != nil {
			svc.l.Error("图片路径错误")
			resList = append(resList, file_domain.FileResponse{
				Filename:  image.Filename,
				IsSuccess: false,
				Msg:       "上传失败",
			})
			continue
		}

		url := "/" + filePath
		err = svc.repo.UpAvatar(ctx, id, url)

		resList = append(resList, file_domain.FileResponse{
			Filename:  image.Filename,
			IsSuccess: true,
			Msg:       "上传成功",
		})
	}

	return resList, err
}
