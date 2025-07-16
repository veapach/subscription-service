package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/db"
	"subscription-service/internal/handler"
)

func main() {

	cfg := config.LoadConfig()

	database := db.InitDB()
	if database == nil {
		log.Fatal("БД не инициализирована")
	}

	r := gin.Default()

	subHandler := handler.SubscriptionHandler{DB: database}
	r.POST("/subscriptions", subHandler.CreateSubscription)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	err := r.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal("Ошибка при запуска сервера:", err)
	}

}
