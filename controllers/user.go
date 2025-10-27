package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Profile godoc
// @Security ApiKeyAuth
// @Summary Get user profile
// @Description Get the profile information of the authenticated user
// @Tags user
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]string "{"email": "user@example.com"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /user/me [get]
func Profile(c *gin.Context) {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "could not retrieve user email"})
		return
	}
	c.JSON(http.StatusOK, &gin.H{"email": userEmail})
}
