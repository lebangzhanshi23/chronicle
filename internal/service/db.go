package service

import (
	"log"

	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.Task{}, &model.TaskLog{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	DB = db
}
