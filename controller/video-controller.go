package controller

import (
	"net/http"

	"company/system/microservices/entity"
	"company/system/microservices/service"
	customValidators "company/system/microservices/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(ctx *gin.Context) error
	ShowAll(ctx *gin.Context)
}

type videoController struct {
	service service.VideoService
}

var validate *validator.Validate

func NewVideoController(service service.VideoService) VideoController {
	validate = validator.New()
	validate.RegisterValidation("is-gmail", customValidators.ValidateIsGmail)
	return &videoController{
		service: service,
	}
}

func (controller *videoController) FindAll() []entity.Video {
	return controller.service.FindAll()
}

func (controller *videoController) Save(context *gin.Context) error {
	var video entity.Video
	err := context.ShouldBindJSON(&video)
	if err != nil {
		return err
	}
	err = validate.Struct(video)
	if err != nil {
		return err
	}
	controller.service.Save(video)
	return nil
}

func (controller *videoController) ShowAll(context *gin.Context) {
	videos := controller.service.FindAll()
	data := gin.H{
		"title":  "Video Page",
		"videos": videos,
	}
	context.HTML(http.StatusOK, "index.html", data)
}
