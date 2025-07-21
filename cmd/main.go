package main

import (
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/db"
	"subscription-service/internal/handler"
	"subscription-service/pkg/logger"

	_ "subscription-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Subscription Service API
// @version 1.0
// @description API сервис для управления онлайн-подписками
// @host localhost:8080
// @BasePath /
func main() {

	logger.Init()
	logger.Log.Info("Запуск сервиса управления подписками")

	cfg := config.LoadConfig()

	database := db.InitDB()
	if database == nil {
		log.Fatal("БД не инициализирована")
	}

	r := gin.Default()

	subHandler := handler.SubscriptionHandler{DB: database}
	r.POST("/subscriptions", subHandler.CreateSubscription)
	r.GET("/subscriptions", subHandler.GetAllSubscriptions)
	r.GET("/subscriptions/:id", subHandler.GetSubscription)
	r.PUT("/subscriptions/:id", subHandler.UpdateSubscription)
	r.DELETE("/subscriptions/:id", subHandler.DeleteSubscription)

	r.GET("/subscriptions/sum", subHandler.GetSubscriptionsSum)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	err := r.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal("Ошибка при запуска сервера:", err)
	}

}
