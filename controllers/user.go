package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "could not retrieve user email"})
		return
	}
	c.JSON(http.StatusOK, &gin.H{"email": userEmail})
}
