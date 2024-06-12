package app

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	DB, err := gorm.Open(sqlite.Open("database/data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	DB.AutoMigrate(&Media{}, &MonitorDir{}, &NotificationTask{})

	return DB
}
