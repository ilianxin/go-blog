package service

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/db"
	"go-blog/internal/model"
	"gorm.io/gorm"
)

func getComments(c *gin.Context, db *gorm.DB) {
	var comments []*model.Comment

	postID, postExists := c.Get("postID")

	if !postExists {
		c.JSON(400, gin.H{"error": "文章ID无效"})
		return
	}

	if err := db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}

	c.JSON(200, comments)
}

func createComment(c *gin.Context, db *gorm.DB) {
	var req model.Comment

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "用户未认证"})
		return
	}

	req.UserID = userID.(uint)

	if err := db.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{"error": "内部错误"})
		return
	}

	c.JSON(201, gin.H{"message": "评论创建成功", "comment_id": req.ID})
}
