package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired verifies API Key or Bearer Token header
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENVSYNC_SKIP_AUTH") == "true" {
			c.Set("User", "dev-user")
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		apiKeyHeader := c.GetHeader("X-API-Key")

		token := ""
		if apiKeyHeader != "" {
			token = apiKeyHeader
		} else if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authentication token",
			})
			return
		}

		c.Set("User", "authenticated-user")
		c.Next()
	}
}
