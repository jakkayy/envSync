package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"

// LoggerAndRequestID returns a Gin middleware for structured logging and correlation tracking
func LoggerAndRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(HeaderXRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header(HeaderXRequestID, requestID)
		c.Set("RequestID", requestID)

		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		fmt.Printf(`{"time":"%s","request_id":"%s","status":%d,"latency":"%s","ip":"%s","method":"%s","path":"%s"}`+"\n",
			time.Now().Format(time.RFC3339),
			requestID,
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
