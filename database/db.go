package database

import (
	"ambassador/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := "ambas:ampass3210@tcp(127.0.0.1:3306)/ambassador_db?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{})
}
