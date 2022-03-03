package main

import (
	"io"
	"os"

	"github.com/MijailStell/go-with-gin/controller"
	"github.com/MijailStell/go-with-gin/middleware"
	"github.com/MijailStell/go-with-gin/service"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func setupLogOutPut() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	setupLogOutPut()
	server := gin.New()
	server.Use(gin.Recovery(),
		middleware.Logger(),
		middleware.BasicAuth(),
		gindump.Dump())

	server.GET("/videos", func(ctx *gin.Context) {
		ctx.JSON(200, videoController.FindAll())
	})

	server.POST("/videos", func(ctx *gin.Context) {
		ctx.JSON(200, videoController.Save(ctx))
	})

	server.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ok!!",
		})
	})

	server.Run(":8080")
}
