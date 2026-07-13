package middleware

import (
	"strings"

	"cms-blog-api/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the JWT token sent in the Authorization header.
// Header format expected: "Authorization: Bearer <jwt_token>"
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondUnauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.RespondUnauthorized(c, "Authorization header must be formatted as Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.RespondUnauthorized(c, "Invalid or expired token: "+err.Error())
			c.Abort()
			return
		}

		// Attach user identity and role to context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
