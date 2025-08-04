package db

import (
	"go-blog/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	user := "rwuser"
	pass := "Rwpass@123"
	dsn := user + ":" + pass + "@tcp(47.111.78.104:3306)/go?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
}

func GetDB() *gorm.DB {
	return db
}
