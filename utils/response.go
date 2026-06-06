package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSONResponse menyajikan format data seragam untuk semua endpoint API
func JSONResponse(c *gin.Context, statusCode int, statusMsg string, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  statusMsg,
		Message: message,
		Data:    data,
	})
}