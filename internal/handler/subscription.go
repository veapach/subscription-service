package handler

import (
	"net/http"
	"subscription-service/internal/model"
	"subscription-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	DB *gorm.DB
}

type createSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       uint    `json:"price" binding:"required,gt=0"`
	UserID      string  `json:"user_id" binding:"required,uuid4"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

// @Summary Создать подписку
// @Description Создание новой подписки для пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body createSubscriptionRequest true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req createSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Ошибка при валидации запроса: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		logger.Log.Warn("Ошибка при парсинге start_date: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date должен быть в формате MM-YYYY"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		ed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			logger.Log.Warn("Ошибка при парсинге end_date: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_date должен быть в формате MM-YYYY"})
			return
		}
		endDate = &ed
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Log.Warn("Ошибка при парсинге user_id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id должен быть UUID"})
		return
	}

	sub := model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userUUID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.DB.Create(&sub).Error; err != nil {
		logger.Log.Error("Ошибка при сохранении подписки: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить подписку"})
		return
	}

	logger.Log.Info("Создана новая подписка: ", sub.ID)
	c.JSON(http.StatusCreated, sub)
}

// @Summary Получить все подписки
// @Description Возвращает список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	var subscriptions []model.Subscription
	if err := h.DB.Find(&subscriptions).Error; err != nil {
		logger.Log.Error("Ошибка при получении подписок: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить подписки"})
		return
	}

	logger.Log.Info("Получено подписок: ", len(subscriptions))
	c.JSON(http.StatusOK, subscriptions)
}

// @Summary Получить подписку по ID
// @Description Возвращает подписку по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	var subscription model.Subscription
	if err := h.DB.First(&subscription, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Log.Warn("Подписка не найдена: ", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Подписка не найдена"})
		} else {
			logger.Log.Error("Не удалось получить подписку: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить подписку"})
		}
		return
	}

	logger.Log.Info("Получена подписка: ", subscription.ID)
	c.JSON(http.StatusOK, subscription)
}

// @Summary Обновить подписку
// @Description Обновление данных существующей подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body createSubscriptionRequest true "Обновлённые данные подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")

	var sub model.Subscription
	if err := h.DB.First(&sub, id).Error; err != nil {
		logger.Log.Warn("Подписка не найдена: ", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "подписка не найдена"})
		return
	}

	var req createSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Ошибка при валидации запроса: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		logger.Log.Warn("Ошибка при парсинге start_date: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date должен быть в формате MM-YYYY"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		ed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			logger.Log.Warn("Ошибка при парсинге end_date: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_date должен быть в формате MM-YYYY"})
			return
		}
		endDate = &ed
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Log.Warn("Ошибка при парсинге user_id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id должен быть UUID"})
		return
	}

	sub.ServiceName = req.ServiceName
	sub.Price = req.Price
	sub.UserID = userUUID
	sub.StartDate = startDate
	sub.EndDate = endDate

	if err := h.DB.Save(&sub).Error; err != nil {
		logger.Log.Error("Ошибка при обновлении подписки: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить подписку"})
		return
	}

	logger.Log.Info("Обновлена подписка: ", sub.ID)
	c.JSON(http.StatusOK, sub)
}

// @Summary Удалить подписку
// @Description Удаление подписки по ID
// @Tags subscriptions
// @Param id path int true "ID подписки"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&model.Subscription{}, id).Error; err != nil {
		logger.Log.Error("Ошибка при удалении подписки: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить подписку"})
		return
	}

	logger.Log.Info("Удалена подписка: ", id)
	c.Status(http.StatusNoContent)
}

// @Summary Подсчитать сумму подписок
// @Description Возвращает суммарную стоимость подписок за период с фильтрацией по user_id и service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string true "Дата начала периода (MM-YYYY)"
// @Param end_date query string true "Дата окончания периода (MM-YYYY)"
// @Success 200 {object} map[string]uint "sum"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/sum [get]
func (h *SubscriptionHandler) GetSubscriptionsSum(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date должен быть в формате MM-YYYY"})
		return
	}

	endDate, err := time.Parse("01-2006", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date должен быть в формате MM-YYYY"})
		return
	}

	query := h.DB.Model(&model.Subscription{}).Where("start_date >= ? AND start_date <= ?", startDate, endDate)

	if userID != "" {
		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id должен быть UUID"})
			return
		}
		query = query.Where("user_id = ?", uid)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	var total uint64
	err = query.Select("SUM(price)").Scan(&total).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подсчете суммы подписок"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_price": total})
}
