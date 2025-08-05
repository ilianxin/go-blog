package service

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/db"
	"go-blog/internal/model"
	"strconv"
)

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func GetComments(c *gin.Context) {
	var comments []*model.Comment

	postID := c.Param("id")

	if postID == "" {
		RespondError(c, "文章ID无效")
		return
	}

	if err := db.GetDB().Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		RespondError(c, "内部错误")
		return
	}

	RespondSuccess(c, comments)
}

func CreateComment(c *gin.Context) {

	idStr := c.Param("id")
	postID64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		RespondError(c, "文章ID无效")
		return
	}
	postID := uint(postID64)

	var req CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, "参数错误")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		RespondError(c, "用户未认证")
		return
	}

	comment := model.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  postID,
	}

	if err := db.GetDB().Create(&comment).Error; err != nil {
		RespondError(c, "内部错误")
		return
	}

	RespondSuccess(c, comment)
}
