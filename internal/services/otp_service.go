package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/username/otp-service/internal/config"
)

const (
	OTP_TTL      = 5 * time.Minute
	COOLDOWN_TTL = 1 * time.Minute
)

// OTPService handles OTP generation, storage, and verification
type OTPService struct{}

// NewOTPService creates a new OTPService instance
func NewOTPService() *OTPService {
	return &OTPService{}
}

// GenerateOTP generates a 6-digit OTP and stores it in Redis with a TTL
func (s *OTPService) GenerateOTP(email string) (string, error) {
	// 1. Check for cooldown
	cooldownKey := fmt.Sprintf("otp_cooldown:%s", email)
	exists, err := config.DB.Exists(config.Ctx, cooldownKey).Result()
	if err != nil {
		return "", err
	}
	if exists > 0 {
		return "", errors.New("please wait before requesting a new OTP")
	}

	// 2. Generate 6-digit OTP
	otp := encodeToString(6)

	// 3. Store in Redis
	otpKey := fmt.Sprintf("otp:%s", email)
	err = config.DB.Set(config.Ctx, otpKey, otp, OTP_TTL).Err()
	if err != nil {
		return "", err
	}

	// 4. Set cooldown
	err = config.DB.Set(config.Ctx, cooldownKey, "1", COOLDOWN_TTL).Err()
	if err != nil {
		return "", err
	}

	// In a real app, you would send this via email/SMS here
	fmt.Printf("[MOCK SEND] OTP for %s is: %s\n", email, otp)

	return otp, nil
}

// VerifyOTP checks if the provided code matches the one in Redis
func (s *OTPService) VerifyOTP(email string, code string) (bool, error) {
	otpKey := fmt.Sprintf("otp:%s", email)
	storedOTP, err := config.DB.Get(config.Ctx, otpKey).Result()
	if err != nil {
		return false, fmt.Errorf("invalid or expired OTP")
	}

	if storedOTP != code {
		return false, errors.New("incorrect OTP")
	}

	// Delete OTP after successful verification to prevent reuse
	config.DB.Del(config.Ctx, otpKey)
	return true, nil
}

// Helper to generate a random n-digit string
func encodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max || err != nil {
		// Fallback for safety, though crypto/rand is robust
		return "123456"
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
