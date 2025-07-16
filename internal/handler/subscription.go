package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"subscription-service/internal/model"
	"time"
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

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req createSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date должен быть в формате MM-YYYY"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		ed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_date должен быть в формате MM-YYYY"})
			return
		}
		endDate = &ed
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить подписку"})
		return
	}

	c.JSON(http.StatusCreated, sub)
}
