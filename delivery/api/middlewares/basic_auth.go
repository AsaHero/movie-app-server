package middlewares

import (
	"encoding/base64"
	"strings"

	"github.com/AsaHero/movie-app-server/delivery/api/outerr"
	"github.com/AsaHero/movie-app-server/pkg/config"
	"github.com/gin-gonic/gin"
)

func BasicAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the authorization header.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			outerr.Unauthorized(c, "Authorization header is required")
			return
		}

		// Decode the provided credentials
		prefix := "Basic "
		if !strings.HasPrefix(authHeader, prefix) {
			outerr.Unauthorized(c, "Invalid authorization header format")
			return
		}

		encodedCredentials := authHeader[len(prefix):]
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			outerr.Unauthorized(c, "Invalid base64 credentials")
			return
		}

		// Check if the decoded credentials are in "username:password" format
		parts := strings.SplitN(string(decodedCredentials), ":", 2)
		if len(parts) != 2 || parts[0] != cfg.Admin.Username || parts[1] != cfg.Admin.Password {
			outerr.Forbidden(c, "Invalid credentials")
			return
		}

		// Proceed to the next middleware/handler.
		c.Next()
	}
}
