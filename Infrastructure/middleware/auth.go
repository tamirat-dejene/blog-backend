package middleware

import (
	"net/http"

	"g6/blog-api/Delivery/bootstrap"
	domain "g6/blog-api/Domain"
	utils "g6/blog-api/Utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware checks if the user is authenticated by verifying the JWT token
func AuthMiddleware(env bootstrap.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := utils.GetCookie(c, "access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found in cookies, please login again"})
			c.Abort()
			return
		}

		// Parse and validate the JWT
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(env.ATS), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		sub, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims: missing user id"})
			return
		}
		c.Set("user_id", sub)

		// Set role if present
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
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
