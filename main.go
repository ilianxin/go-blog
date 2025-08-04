package main

import (
	"github.com/gin-gonic/gin"
	"go-blog/config"
	"go-blog/internal/db"
	"go-blog/internal/middleware"
	"go-blog/internal/model"
)

func main() {

	// 初始化数据库
	db.InitDB()

	// 启动 Gin 路由
	r := gin.Default()
	auth := r.Group("/api", middleware.JWTAuthMiddleware(config.JWT_SECRET))
	auth.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "已认证"})
	})

}
