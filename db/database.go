package db

import (
	"go-gin-auth/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	line := "+-------------------------------+"
	if err != nil {
		println(line)
		println("| ❌ Failed to connect database |")
		println(line)
		panic(err)
	}
	println("+--------------------------+")
	println("| ✅ Connected to database |")
	println("+--------------------------+")
	err = db.AutoMigrate(&models.User{}, &models.File{})
	if err != nil {
		println(line)
		println("| ❌ Failed to migrate database |")
		println(line)
		panic(err)
	}
	DB = db
}
