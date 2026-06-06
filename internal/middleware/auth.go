package middleware

import (
	"net/http"
	"strings"

	"fashion-store/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Periksa token di Header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 2. Jika tidak ada di header, periksa token di Cookie (untuk rendering HTML)
		if tokenString == "" {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenString = cookie
			}
		}

		// Jika token tidak ditemukan
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi habis atau tidak valid, silakan masuk terlebih dahulu"})
			c.Abort()
			return
		}

		// Validasi token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token kedaluwarsa atau tidak valid, silakan masuk kembali"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Data token tidak dapat diuraikan"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial pengguna di dalam token tidak valid"})
			c.Abort()
			return
		}

		role, _ := claims["role"].(string)

		// Simpan informasi User ID dan Role di context Gin agar bisa dipakai di Handlers
		c.Set("userID", uint(userIDFloat))
		c.Set("role", role)

		c.Next()
	}
}