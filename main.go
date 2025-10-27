package main

import (
	"go-gin-auth/controllers"
	"go-gin-auth/db"
	"go-gin-auth/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	router := gin.Default()
	db.ConnectDB()
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
	}
	usersGroup := router.Group("/user")
	usersGroup.Use(utils.AuthMiddleware())
	{
		usersGroup.GET("/me", controllers.Profile)
		usersGroup.POST("/upload", controllers.UploadFile)
		usersGroup.GET("/download", controllers.ListUserFiles)
		usersGroup.GET("/download/:fileID", controllers.DownloadFile)
	}
	log.Fatal(router.Run(":6969"))
}
