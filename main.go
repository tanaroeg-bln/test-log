package main

import (
	"os"
	"poc-log/handlers"

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
	// Initialize Zap logger
	zapProductionConfig := zap.NewProductionConfig()

	//  add getGCPEncoderConfig
	zapProductionConfig.EncoderConfig = getGCPEncoderConfig()

	zapProductionLogger, _ := zapProductionConfig.Build()

	zap.ReplaceGlobals(zapProductionLogger)

	defer zapProductionLogger.Sync()

	// Gin router without default logger (since you use zap)
	r := gin.New()
	r.Use(gin.Recovery())

	// Optional: custom middleware to log each request with zap
	r.Use(func(c *gin.Context) {
		c.Next()

		zap.S().With(
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		).Info("request completed")
	})

	// Define a simple GET endpoint
	r.GET("/test-log", handlers.TestLog)

	r.POST("/test-body-log", handlers.TestBodyLog)
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
