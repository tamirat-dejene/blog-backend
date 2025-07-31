package middleware

import (
	"net/http"
	"strings"

	"g6/blog-api/Delivery/bootstrap"
	domain "g6/blog-api/Domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(env bootstrap.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(env.ATS), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["userId"].(string))
		if role, ok := claims["role"]; ok {
			c.Set("role", role.(string))
		}
		c.Next()
	}
}

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != string(domain.RoleSuperAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "super_admin only"})
			return
		}
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != string(domain.RoleAdmin) && role != string(domain.RoleSuperAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin or super_admin only"})
			return
		}
		c.Next()
	}
}
