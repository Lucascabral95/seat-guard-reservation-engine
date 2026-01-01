package middleware

import (
	"booking-service/pkg/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// func UserMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		internalKey := os.Getenv("SECRET_X_INTERNAL_SECRET")
// 		if internalKey != "" && c.GetHeader("X-Internal-Secret") == internalKey {
// 			c.Next()
// 			return
// 		}
// 		// -------------------------------

// 		tokenString, err := utils.ExtractToken(c)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			c.Abort()
// 			return
// 		}

// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, jwt.ErrSignatureInvalid
// 			}
// 			return []byte(os.Getenv("JWT_SECRET")), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		// if userID, ok := claims["sub"].(string); ok {
// 		// 	c.Set("userID", userID)
// 		// }
// 		if userID, ok := claims["id"].(string); ok {
// 			c.Set("userID", userID)
// 		}

//			c.Next()
//		}
//	}
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ✅ 0. DEJAR PASAR PREFLIGHT CORS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 1. Secreto interno
		internalKey := os.Getenv("SECRET_X_INTERNAL_SECRET")
		if internalKey != "" && c.GetHeader("X-Internal-Secret") == internalKey {
			c.Set("userID", "internal")
			c.Set("user_id", "internal")
			c.Next()
			return
		}

		// 2. Token
		tokenString, err := utils.ExtractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		var userID string
		if id, ok := claims["id"].(string); ok {
			userID = id
		} else if sub, ok := claims["sub"].(string); ok {
			userID = sub
		}

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User identity missing"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("user_id", userID)

		c.Next()
	}
}
