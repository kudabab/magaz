package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kudabab/market-s/entity"
	"github.com/kudabab/market-s/service"
)

var jwtKey = []byte("my_secret_key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &entity.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		user, err := service.GetUserByUsername(context.Background(), claims.Username)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "user not fouind"})
			c.Abort()
			return
		}

		fmt.Println("Auth user: ", user)

		c.Set("user", user)
		c.Next()
	}
}
