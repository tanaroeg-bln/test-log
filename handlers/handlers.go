package handlers

import (
	"errors"
	"math/rand"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestLog(c *gin.Context) {
	zap.S().With(
		"id", 1,
	).Info("Test info")

	zap.S().With(
		"id", 1,
	).Warn("Test warn")

	err := errors.New("mock error")
	zap.S().With(
		"id", 2,
		"error", err,
	).Error("Test error")

	zap.S().With(
		"randomNo", rand.Intn(100),
	).Info("Test info with random number")

	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

func TestBodyLog(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		zap.S().Error("Failed to bind JSON", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	c.JSON(200, gin.H{"message": "Request body logged"})
}
