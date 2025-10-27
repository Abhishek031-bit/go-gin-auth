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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.AuthRequest true "User registration details"
// @Success 201 {object} map[string]string "{"message": "user created"}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 409 {object} map[string]string "{"error": "email already exists"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /auth/register [post]
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

// Login godoc
// @Summary Log in a user
// @Description Log in a user with email and password to get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.AuthRequest true "User login details"
// @Success 200 {object} map[string]string "{"token": "jwt_token_string"}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 401 {object} map[string]string "{"error": "invalid email or password"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /auth/login [post]
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
