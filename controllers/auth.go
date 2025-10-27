package controllers

import (
	"go-gin-auth/db"
	"go-gin-auth/models"
	"go-gin-auth/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input models.AuthRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "error hashing password"})
		return
	}
	user := models.User{Email: input.Email, Password: string(hashedPassword)}

	if result := db.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusConflict, &gin.H{"error": "email already exists"})
		return
	}
	c.JSON(http.StatusCreated, &gin.H{"message": "user created"})
}

func Login(c *gin.Context) {
	var input models.AuthRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := db.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, &gin.H{"error": "invalid email or password"})
		return
	}

	if os.Getenv("JWT_SECRET") == "" {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "JWT secret not set"})
		return
	}

	tokenString, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "error generating token"})
		return
	}

	c.JSON(http.StatusOK, &gin.H{"token": tokenString})
}
