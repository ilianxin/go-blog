package service

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go-blog/config"
	"go-blog/internal/db"
	"go-blog/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		RespondError(c, "参数错误: "+err.Error())
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondError(c, "Failed to hash password")
		return
	}
	user.Password = string(hashedPassword)

	if err := db.GetDB().Create(&user).Error; err != nil {
		RespondError(c, "Failed to create user")
		return
	}

	RespondSuccess(c, "User registered successfully")
}

func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		RespondError(c, "参数错误: "+err.Error())
		return
	}

	var storedUser model.User
	if err := db.GetDB().Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		RespondError(c, "Invalid username or password")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		RespondError(c, "Invalid username or password")
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		RespondError(c, "Failed to generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

	return
}
