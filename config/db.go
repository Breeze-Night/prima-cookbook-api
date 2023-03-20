package config

import (
	"prima_cookbook/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "root@tcp(127.0.0.1:3306)/prima_cookbook?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	DB.AutoMigrate(&models.User{}, &models.Recipe{}, &models.Ingredient{})

}
