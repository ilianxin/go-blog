package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 统一错误响应
func RespondErrorWithCode(c *gin.Context, code int, msg string) {
	logrus.Warn(msg)
	c.JSON(code, gin.H{"error": msg})
}

// 统一错误响应
func RespondError(c *gin.Context, msg string) {
	logrus.Warn(msg)
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}
