package handlers

import (
	"fmt "
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/otp-service/internal/config"
	"github.com/username/otp-service/internal/models"
	"github.com/username/otp-service/internal/services"
)

type TOTPHandler struct {
	totpService *services.TOTPService
}

func NewTOTPHandler(s *services.TOTPService) *TOTPHandler {
	return &TOTPHandler{totpService: s}
}

// SetupTOTP handles the TOTP registration flow
func (h *TOTPHandler) SetupTOTP(c *gin.Context) {
	var req models.OTPRequest // Just reuse the email request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	secret, qrCode, err := h.totpService.SetupTOTP(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate TOTP setup",
		})
		return
	}

	// Store secret in Redis for testing (In production, this belongs in a persistent DB)
	secretKey := fmt.Sprintf("totp_secret:%s", req.Email)
	err = config.DB.Set(config.Ctx, secretKey, secret, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to save TOTP secret",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "TOTP setup generated",
		Data: models.TOTPSetupResponse{
			Secret: secret,
			QRCode: qrCode,
		},
	})
}

// VerifyTOTP handles the TOTP verification flow
func (h *TOTPHandler) VerifyTOTP(c *gin.Context) {
	var req models.TOTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	// Retrieve secret from Redis
	secretKey := fmt.Sprintf("totp_secret:%s", req.Email)
	secret, err := config.DB.Get(config.Ctx, secretKey).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  http.StatusNotFound,
			Message: "TOTP not set up for this email",
		})
		return
	}

	success, err := h.totpService.VerifyTOTP(req.Email, secret, req.Code)
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
		Message: "TOTP verified successfully",
	})
}
