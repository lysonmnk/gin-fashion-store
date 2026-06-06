package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination menyimpan data konfigurasi pembagian halaman
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// GetPagination membaca parameter 'page' dan 'limit' dari URL Query
func GetPagination(c *gin.Context) Pagination {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "12")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 12
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}