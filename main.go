package main

import (
	"math/rand"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type gcpSeverityHook struct{}

func (h *gcpSeverityHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *gcpSeverityHook) Fire(e *logrus.Entry) error {
	severity := map[logrus.Level]string{
		logrus.PanicLevel: "CRITICAL",
		logrus.FatalLevel: "CRITICAL",
		logrus.ErrorLevel: "ERROR",
		logrus.WarnLevel:  "WARNING",
		logrus.InfoLevel:  "INFO",
		logrus.DebugLevel: "DEBUG",
		logrus.TraceLevel: "DEBUG",
	}[e.Level]

	e.Data["severity"] = severity
	return nil
}

func main() {
	// Configure logrus
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(&gcpSeverityHook{}) // 👈 attach the hook

	// Gin router without default logger (since you use logrus)
	r := gin.New()
	r.Use(gin.Recovery())

	// Optional: custom middleware to log each request with logrus
	r.Use(func(c *gin.Context) {
		c.Next()
		logrus.WithFields(logrus.Fields{
			"status": c.Writer.Status(),
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("request completed")
	})

	// Define a simple GET endpoint
	r.GET("/test-log", func(c *gin.Context) {
		ctx := c.Request.Context()

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"id": 1,
		}).Info("Test info")

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"id": 1,
		}).Warn("Test warn")

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"id": 2,
		}).Error("Test error")

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"randomNo": rand.Intn(100),
		}).Info("Test info with random number")

		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.POST("/test-body-log", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			logrus.WithError(err).Error("Failed to bind JSON")
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
