package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/otp-service/internal/models"
	"github.com/username/otp-service/internal/services"
)

type OTPHandler struct {
	otpService *services.OTPService
}

func NewOTPHandler(s *services.OTPService) *OTPHandler {
	return &OTPHandler{otpService: s}
}

// GenerateOTP handles the request to generate and send an OTP
func (h *OTPHandler) GenerateOTP(c *gin.Context) {
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	_, err := h.otpService.GenerateOTP(req.Email)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, models.Response{
			Status:  http.StatusTooManyRequests,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "OTP sent successfully (mocked)",
	})
}

// VerifyOTP handles the OTP verification request
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req models.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	success, err := h.otpService.VerifyOTP(req.Email, req.Code)
	if err != nil || !success {
		message := "Verification failed"
		if err != nil {
			message = err.Error()
		}
		c.JSON(http.StatusUnauthorized, models.Response{
			Status:  http.StatusUnauthorized,
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "OTP verified successfully",
	})
}
