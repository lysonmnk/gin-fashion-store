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

		// 2. Jika tidak ada di header, periksa token di Cookie
		if tokenString == "" {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenString = cookie
			}
		}

		// Deteksi apakah request dari browser atau API
		isAPIRequest := strings.HasPrefix(c.Request.URL.Path, "/api")

		// Jika token tidak ditemukan
		if tokenString == "" {
			if isAPIRequest {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi habis atau tidak valid, silakan masuk terlebih dahulu"})
			} else {
				c.Redirect(http.StatusSeeOther, "/login")
			}
			c.Abort()
			return
		}

		// Validasi token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			if isAPIRequest {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token kedaluwarsa atau tidak valid, silakan masuk kembali"})
			} else {
				c.SetCookie("token", "", -1, "/", "", false, true)
				c.Redirect(http.StatusSeeOther, "/login")
			}
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if isAPIRequest {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Data token tidak dapat diuraikan"})
			} else {
				c.Redirect(http.StatusSeeOther, "/login")
			}
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			if isAPIRequest {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial pengguna di dalam token tidak valid"})
			} else {
				c.Redirect(http.StatusSeeOther, "/login")
			}
			c.Abort()
			return
		}

		role, _ := claims["role"].(string)

		c.Set("userID", uint(userIDFloat))
		c.Set("role", role)

		c.Next()
	}
}

// GetNavbarData membaca cookie token dan mengembalikan data isLoggedIn & isAdmin untuk navbar
func GetNavbarData(c *gin.Context) gin.H {
	cookie, err := c.Cookie("token")
	if err != nil || cookie == "" {
		return gin.H{
			"isLoggedIn": false,
			"isAdmin":    false,
		}
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return gin.H{
			"isLoggedIn": false,
			"isAdmin":    false,
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return gin.H{
			"isLoggedIn": false,
			"isAdmin":    false,
		}
	}

	role, _ := claims["role"].(string)

	return gin.H{
		"isLoggedIn": true,
		"isAdmin":    role == "admin",
	}
}