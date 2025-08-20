package main

import (
	"math/rand"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GCP Severity Level mapping for Zap
func gcpSeverityLevel(level zapcore.Level) string {
	switch level {
	case zapcore.PanicLevel, zapcore.FatalLevel:
		return "CRITICAL"
	case zapcore.ErrorLevel:
		return "ERROR"
	case zapcore.WarnLevel:
		return "WARNING"
	case zapcore.InfoLevel:
		return "INFO"
	case zapcore.DebugLevel:
		return "DEBUG"
	default:
		return "INFO"
	}
}

// Custom encoder config for GCP logging
func getGCPEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.LevelKey = "severity"
	config.MessageKey = "message"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(gcpSeverityLevel(level))
	}
	return config
}

func main() {
	// Configure Zap logger with GCP-compatible format
	config := getGCPEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()

	// Gin router without default logger (since you use zap)
	r := gin.New()
	r.Use(gin.Recovery())

	// Optional: custom middleware to log each request with zap
	r.Use(func(c *gin.Context) {
		c.Next()
		logger.Info("request completed",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
	})

	// Define a simple GET endpoint
	r.GET("/test-log", func(c *gin.Context) {
		logger.Info("Test info",
			zap.Int("id", 1),
		)

		logger.Warn("Test warn",
			zap.Int("id", 1),
		)

		logger.Error("Test error",
			zap.Int("id", 2),
		)

		logger.Info("Test info with random number",
			zap.Int("randomNo", rand.Intn(100)),
		)

		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.POST("/test-body-log", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		c.JSON(200, gin.H{"message": "Request body logged"})
	})
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
