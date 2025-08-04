package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 统一错误响应
func RespondErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

// 统一错误响应
func RespondError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}
