package service

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/db"
	"go-blog/internal/model"
)

func GetComments(c *gin.Context) {
	var comments []*model.Comment

	postID, postExists := c.Get("postID")

	if !postExists {
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
	var req model.Comment

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, "参数错误")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		RespondError(c, "用户未认证")
		return
	}

	req.UserID = userID.(uint)

	if err := db.GetDB().Create(&req).Error; err != nil {
		RespondError(c, "内部错误")
		return
	}

	RespondSuccess(c, "评论创建成功")
}
