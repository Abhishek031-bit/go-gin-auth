package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Show the status of the API.
// @Description get the status of API.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{"message": "API is healthy"})
}
