package main

import (
	"github.com/gin-gonic/gin"
	"go-blog/config"
	"go-blog/internal/db"
	"go-blog/internal/middleware"
	"go-blog/internal/service"
)

func main() {

	// 初始化数据库
	db.InitDB()

	// 启动 Gin 路由
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", service.Register)
		api.POST("/login", service.Login)
	}

	// 需要认证的接口
	auth := r.Group("/api", middleware.JWTAuthMiddleware(config.JWT_SECRET))
	{
		auth.GET("/posts", service.ReadPost)
		auth.POST("/posts", service.CreatePost)
		auth.DELETE("/posts/:id", service.DeletePost)
		auth.PUT("/posts/:id", service.UpdatePost)
		auth.GET("/comments/:id", service.GetComments)
		auth.POST("/comments/:id", service.CreateComment)
		auth.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "已认证"})
		})
	}

	r.Run()
}
