package service

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/model"
	"gorm.io/gorm"
	"net/http"
)

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func createPost(c *gin.Context, db *gorm.DB) (uint, error) {
	// 这里应该调用数据库操作来创建文章
	// 假设我们使用 GORM 进行数据库操作

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return 0, err
	}

	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint), // 假设 userID 是 uint 类型
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "内部错误"})
		return 0, err
	}

	return post.ID, nil
}

func readPost(c *gin.Context, db *gorm.DB) {
	postID, postExists := c.Get("postID")

	userID, userExists := c.Get("userID")

	if !userExists && !postExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户未认证或文章ID无效"})
		return
	}

	if postExists {
		var post model.Post

		if err := db.First(&post, postID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
			return
		}

		c.JSON(http.StatusOK, post)
	}

	if userExists {
		var posts []model.Post
		if err := db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
			return
		}
		c.JSON(http.StatusOK, posts)
	}

}

func updatePost(c *gin.Context, db *gorm.DB) {
	postID, exists := c.Get("postID")

	userId, userExists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID无效"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var post model.Post
	if err := db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	if userExists {
		if post.UserID != userId.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权更新此文章"})
			return
		}
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func deletePost(c *gin.Context, db *gorm.DB) {
	postID, exists := c.Get("postID")

	userId, userExists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID无效"})
		return
	}

	var post model.Post
	if err := db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	if userExists {
		if post.UserID != userId.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文章"})
			return
		}
	}

	if err := db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}
