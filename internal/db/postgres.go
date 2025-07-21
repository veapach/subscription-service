package db

import (
	"fmt"
	"subscription-service/internal/config"
	"subscription-service/internal/model"
	"subscription-service/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *gorm.DB
	var err error

	// Попытки подключения (нужна задержка, так как база еще не готова, когда запускается контейнер)
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		logger.Log.Println("Ожидание подключения к БД...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Log.Fatal("Не удалось подключиться к БД после нескольких попыток: ", err)
	}

	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		logger.Log.Fatal("Не удалось выполнить миграции: ", err)
	}

	return db
}
