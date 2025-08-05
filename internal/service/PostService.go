package service

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/db"
	"go-blog/internal/model"
	"net/http"
)

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreatePost(c *gin.Context) {
	// 这里应该调用数据库操作来创建文章
	// 假设我们使用 GORM 进行数据库操作

	userID, exists := c.Get("userID")

	if !exists {
		RespondError(c, "User not authenticated")
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, "参数错误")
		return
	}

	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint), // 假设 userID 是 uint 类型
	}

	if err := db.GetDB().Create(&post).Error; err != nil {
		RespondError(c, "内部错误")
		return
	}

	RespondSuccess(c, post)

}

func ReadPost(c *gin.Context) {

	postID, postExists := c.Get("postID")
	userID, userExists := c.Get("userID")

	if !userExists {
		RespondError(c, "用户未认证")
		return
	}

	if postExists {
		var post model.Post

		if err := db.GetDB().First(&post, postID).Error; err != nil {
			RespondError(c, "文章未找到")
			return
		}

		c.JSON(http.StatusOK, post)

		return
	}

	if userExists {
		var posts []model.Post
		if err := db.GetDB().Where("user_id = ?", userID).Find(&posts).Error; err != nil {
			RespondError(c, "获取文章失败")
			return
		}
		c.JSON(http.StatusOK, posts)
	}

}

func UpdatePost(c *gin.Context) {

	postID := c.Param("id")

	if postID == "" {
		RespondError(c, "文章ID无效")
		return
	}

	userId, userExists := c.Get("userID")

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, "参数错误")
		return
	}

	var post model.Post
	if err := db.GetDB().First(&post, postID).Error; err != nil {
		RespondError(c, "文章未找到")
		return
	}

	if userExists {
		if post.UserID != userId.(uint) {
			RespondError(c, "无权更新此文章")
			return
		}
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := db.GetDB().Save(&post).Error; err != nil {
		RespondError(c, "更新文章失败")
		return
	}

	RespondSuccess(c, post)
}

func DeletePost(c *gin.Context) {
	postID := c.Param("postID")

	userId, userExists := c.Get("userID")

	if postID == "" {
		RespondError(c, "文章ID无效")
		return
	}

	var post model.Post
	if err := db.GetDB().First(&post, postID).Error; err != nil {
		RespondError(c, "文章未找到")
		return
	}

	if userExists {
		if post.UserID != userId.(uint) {
			RespondError(c, "无权删除此文章")
			return
		}
	}

	if err := db.GetDB().Delete(&post).Error; err != nil {
		RespondError(c, "删除文章失败")
		return
	}

	RespondSuccess(c, "文章已删除")
}
