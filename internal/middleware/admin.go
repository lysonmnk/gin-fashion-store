package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		// BUG FIX: Bedakan response JSON (API) vs redirect (halaman web)
		// Sebelumnya selalu return JSON 403, meski diakses dari browser
		isAPIRequest := strings.HasPrefix(c.Request.URL.Path, "/api")

		if !exists || role != "admin" {
			if isAPIRequest {
				c.JSON(http.StatusForbidden, gin.H{"error": "Hak akses ditolak, rute ini hanya khusus untuk Administrator"})
			} else {
				c.Redirect(http.StatusSeeOther, "/")
			}
			c.Abort()
			return
		}
		c.Next()
	}
}