package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/username/otp-service/internal/config"
	"github.com/username/otp-service/internal/handlers"
	"github.com/username/otp-service/internal/services"
)

func main() {
	// 1. Initialize Redis
	if err := config.InitRedis(); err != nil {
		log.Fatalf("Critical error: %v", err)
	}

	// 2. Initialize Services
	otpService := services.NewOTPService()
	totpService := services.NewTOTPService()

	// 3. Initialize Handlers
	otpHandler := handlers.NewOTPHandler(otpService)
	totpHandler := handlers.NewTOTPHandler(totpService)

	// 4. Setup Router
	r := gin.Default()

	// 5. Define Routes
	api := r.Group("/api")
	{
		otp := api.Group("/otp")
		{
			otp.POST("/generate", otpHandler.GenerateOTP)
			otp.POST("/verify", otpHandler.VerifyOTP)
		}

		totp := api.Group("/totp")
		{
			totp.POST("/setup", totpHandler.SetupTOTP)
			totp.POST("/verify", totpHandler.VerifyTOTP)
		}
	}

	// 6. Start Server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
