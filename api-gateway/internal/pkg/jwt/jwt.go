package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	JWTSecret = "abduazim"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Println("Received Token:", tokenString)

		if tokenString == "" {
			c.JSON(401, gin.H{"error": "authorization token is required"})
			c.Abort()
			return
		}
		tokenString = "Bearer " + tokenString

		tokenParts := strings.Split(tokenString, " ")
		fmt.Println("Token Parts:", tokenParts)

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		})

		if err != nil {
			fmt.Println("Token Parse Error:", err)
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims["exp"].(float64))
			fmt.Println("Token Expiration Time:", time.Unix(exp, 0))
			fmt.Println("Current Time:", time.Now())
			if time.Now().Unix() > exp {
				c.JSON(401, gin.H{"error": "token expired"})
				c.Abort()
				return
			}
			c.Set("email", claims["email"])
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
		}
	}
}
