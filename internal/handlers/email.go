package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"email-service/internal/email"
)

type EmailHandler struct {
	emailService *email.EmailService
}

func NewEmailHandler(emailService *email.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

type SendAuthCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=update registration"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (h *EmailHandler) SendAuthCode(c *gin.Context) {
	var req SendAuthCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	if err := h.emailService.SendAuthCode(req.Email, req.Code, req.Type); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to send email",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Authentication code sent successfully",
	})
}

func (h *EmailHandler) HealthCheck(c *gin.Context) {
	if err := h.emailService.TestConnection(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"service": "email-service",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "email-service",
	})
}
