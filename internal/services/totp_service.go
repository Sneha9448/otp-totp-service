package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"github.com/username/otp-service/internal/config"
)

// TOTPService handles TOTP generation and verification
type TOTPService struct{}

// NewTOTPService creates a new TOTPService instance
func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

// SetupTOTP generates a secret and a QR code for a user
func (s *TOTPService) SetupTOTP(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GoOTPApp",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}

	secret := key.Secret()

	// Generate QR code as PNG bytes
	var png []byte
	png, err = qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return "", "", err
	}

	// Double encode: png -> base64 string
	qrBase64 := base64.StdEncoding.EncodeToString(png)

	return secret, qrBase64, nil
}

// VerifyTOTP validates a TOTP code and prevents reuse in the same time window
func (s *TOTPService) VerifyTOTP(email string, secret string, code string) (bool, error) {
	// 1. Basic validation
	valid := totp.Validate(code, secret)
	if !valid {
		return false, fmt.Errorf("invalid TOTP code")
	}

	// 2. Reuse prevention
	// Store the used code in Redis for 30 seconds (default TOTP window)
	reuseKey := fmt.Sprintf("totp_used:%s:%s", email, code)
	exists, err := config.DB.Exists(config.Ctx, reuseKey).Result()
	if err != nil {
		return false, err
	}
	if exists > 0 {
		return false, fmt.Errorf("this TOTP code has already been used")
	}

	// Mark as used
	err = config.DB.Set(config.Ctx, reuseKey, "1", 30*time.Second).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}
