package controller

import (
	"github.com/MijailStell/go-with-gin/entity"
	"github.com/MijailStell/go-with-gin/service"
	"github.com/gin-gonic/gin"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(ctx *gin.Context) entity.Video
}

type videoController struct {
	service service.VideoService
}

func New(service service.VideoService) VideoController {
	return &videoController{
		service: service,
	}
}

func (controller *videoController) FindAll() []entity.Video {
	return controller.service.FindAll()
}

func (controller *videoController) Save(context *gin.Context) entity.Video {
	var video entity.Video
	context.BindJSON(&video)
	controller.service.Save(video)
	return video
}
