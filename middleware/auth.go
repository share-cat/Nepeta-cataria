package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"share-cat/Nepeta-cataria/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		log.Printf("Received token: %s", token)
		token = strings.TrimPrefix(token, "Bearer ")

		parsedToken, err := utils.ParseToken(token)
		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userint, err := strconv.Atoi(claims["userID"].(string))
		if err != nil {
			log.Printf("Error parsing userID: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set("userID", userint)
		c.Next()
	}
}
