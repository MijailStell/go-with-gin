package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"company/system/microservices/controller"
	"company/system/microservices/middleware"
	"company/system/microservices/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	gindump "github.com/tpkeeper/gin-dump"
)

var (
	videoService    service.VideoService       = service.NewVideoService()
	loginService    service.LoginService       = service.NewLoginService()
	jwtService      service.JWTService         = service.NewJWTService()
	videoController controller.VideoController = controller.NewVideoController(videoService)
	loginController controller.LoginController = controller.NewLoginController(loginService, jwtService)
)

const VIDEOS_ROOT = "/videos"

func setupLogOutPut() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func setupConfig() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	mode, ok := viper.Get("SERVICE.MODE").(string)

	if ok {
		gin.SetMode(mode)
	}
	viper.SetDefault("SERVICE.PORT", "5000")
}

func main() {
	setupLogOutPut()
	setupConfig()
	server := gin.New()

	server.Use(gin.Recovery(),
		middleware.Logger(),
		gindump.Dump())

	server.Static("/css", "./templates/css")

	server.LoadHTMLGlob("templates/*.html")

	// Login Endpoint: Authentication + Token creation
	server.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	// JWT Authorization Middleware applies to "/api" only.
	apiRoutes := server.Group("/api", middleware.AuthorizeJWT())
	{
		apiRoutes.GET(VIDEOS_ROOT, func(ctx *gin.Context) {
			ctx.JSON(200, videoController.FindAll())
		})

		apiRoutes.POST(VIDEOS_ROOT, func(ctx *gin.Context) {
			err := videoController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video es válido."})
			}
		})
	}

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET(VIDEOS_ROOT, videoController.ShowAll)
	}

	servicePort, ok := viper.Get("SERVICE.PORT").(string)
	if os.Getenv("ASPNETCORE_PORT") != "" {
		servicePort = os.Getenv("ASPNETCORE_PORT")
		fmt.Println(servicePort)
	}

	// if type assert is not valid it will throw an error
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	if servicePort == "" {
		servicePort = "5000"
	}

	server.Run(":" + servicePort)
}
