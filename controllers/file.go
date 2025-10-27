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

// UploadFile godoc
// @Security ApiKeyAuth
// @Summary Upload a file
// @Description Uploads a file for the authenticated user
// @Tags user
// @Accept mpfd
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string "{"message": "File uploaded successfully", "file_name": "example.txt", "path": "uploads/user_at_example.com/example.txt"}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /user/upload [post]
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

	c.JSON(http.StatusCreated, &gin.H{"message": "File uploaded successfully", "file_name": file.Filename, "path": fullFilePath})
}

// ListUserFiles godoc
// @Security ApiKeyAuth
// @Summary List user files
// @Description Lists all files uploaded by the authenticated user
// @Tags user
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{} "{"total": 1, "files": [{"ID": 1, "CreatedAt": "2023-01-01T00:00:00Z", "UpdatedAt": "2023-01-01T00:00:00Z", "DeletedAt": null, "FilePath": "uploads/user_at_example.com/file.txt", "OriginalName": "file.txt", "UserID": 1}]}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /user/download [get]
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

// DownloadFile godoc
// @Security ApiKeyAuth
// @Summary Download a user file
// @Description Downloads a specific file owned by the authenticated user
// @Tags user
// @Accept */*
// @Produce octet-stream
// @Param fileID path int true "File ID"
// @Success 200 {file} byte "File content"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 403 {object} map[string]string "{"error": "Forbidden"}"
// @Failure 404 {object} map[string]string "{"error": "File not found"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /user/download/{fileID} [get]
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
