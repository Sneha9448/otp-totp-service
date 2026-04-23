package models

// Response is a generic API response structure
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OTPRequest represents the input for OTP generation
type OTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// OTPVerifyRequest represents the input for OTP verification
type OTPVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// TOTPSetupResponse represents the data returned after TOTP setup
type TOTPSetupResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code_base64"` // Base64 encoded QR code image
}

// TOTPVerifyRequest represents the input for TOTP verification
type TOTPVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}
