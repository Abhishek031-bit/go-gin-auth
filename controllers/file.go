package controllers

import (
	"fmt"
	"go-gin-auth/db"
	"go-gin-auth/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, &gin.H{"error": "Unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": fmt.Sprintf("File upload failed: %v", err)})
		return
	}

	safeEmail := strings.ReplaceAll(userEmail.(string), "@", "_at_")

	uploadDir := filepath.Join("uploads", safeEmail)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": fmt.Sprintf("Could not create upload directory: %v", err)})
		return
	}

	fullFilePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, fullFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": fmt.Sprintf("Could not save uploaded file: %v", err)})
		return
	}

	var user models.User
	if result := db.DB.Where("email = ?", userEmail).First(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": fmt.Sprintf("Could not find user: %v", result.Error)})
		return
	}

	fileRecord := models.File{
		UserID:       user.ID,
		OriginalName: file.Filename,
		FilePath:     fullFilePath,
	}

	if result := db.DB.Create(&fileRecord); result.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": fmt.Sprintf("Could not save file record: %v", result.Error)})
		return
	}

	c.JSON(http.StatusOK, &gin.H{"message": "File uploaded successfully", "file_name": file.Filename, "path": fullFilePath})
}

func ListUserFiles(c *gin.Context) {
	userEmail, _ := c.Get("userEmail")
	var user models.User
	if result := db.DB.Where("email = ?", userEmail).First(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "Could not find user"})
		return
	}
	var files []models.File
	if result := db.DB.Where("user_id = ?", user.ID).Find(&files); result.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "Could not find files"})
		return
	}
	c.JSON(http.StatusOK, &gin.H{"total": len(files), "files": files})
}

func DownloadFile(c *gin.Context) {
	fileID := c.Param("fileID")

	userEmail, _ := c.Get("userEmail")
	var file models.File

	if result := db.DB.First(&file, fileID); result.Error != nil {
		c.JSON(http.StatusNotFound, &gin.H{"error": "File not found"})
		return
	}

	var owner models.User

	db.DB.Select("email").First(&owner, file.UserID)
	if owner.Email != userEmail {
		c.JSON(http.StatusForbidden, &gin.H{"error": "You do not own this file"})
		return
	}
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, &gin.H{"error": "File not found on server"})
		return
	}
	c.File(file.FilePath)
}
