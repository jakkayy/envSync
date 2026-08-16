package middleware

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var (
	TotalRequests uint64
	TotalErrors   uint64
)

// MetricsMiddleware tracks HTTP request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		atomic.AddUint64(&TotalRequests, 1)
		c.Next()
		if c.Writer.Status() >= 400 {
			atomic.AddUint64(&TotalErrors, 1)
		}
	}
}

// MetricsHandler serves Prometheus-format metrics
func MetricsHandler(c *gin.Context) {
	metrics := fmt.Sprintf(`# HELP envsync_http_requests_total Total number of HTTP requests
# TYPE envsync_http_requests_total counter
envsync_http_requests_total %d

# HELP envsync_http_errors_total Total number of HTTP 4xx/5xx errors
# TYPE envsync_http_errors_total counter
envsync_http_errors_total %d
`, atomic.LoadUint64(&TotalRequests), atomic.LoadUint64(&TotalErrors))

	c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(metrics))
}
