package middleware

import (
	"net/http"

	"prima_cookbook/auth"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Request need access token",
				"status":  http.StatusUnauthorized,
			})

			c.Abort()
			return
		}

		email, _, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
				"status":  http.StatusUnauthorized,
			})

			c.Abort()
			return
		}

		c.Set("x-email", email)

		c.Next()
	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Request need access token",
				"status":  http.StatusUnauthorized,
			})

			c.Abort()
			return
		}

		email, role, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
				"status":  http.StatusUnauthorized,
			})

			c.Abort()
			return
		}

		// role 1 == admin
		if role != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Your role forbids you to access this page",
				"status":  http.StatusUnauthorized,
			})

			c.Abort()
			return
		}

		c.Set("x-email", email)
		c.Next()
	}
}
