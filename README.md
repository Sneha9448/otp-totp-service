# Go OTP & TOTP Service

A production-ready authentication microservice built with **Golang** and **Gin**, providing robust One-Time Password (OTP) and Time-based One-Time Password (TOTP) functionalities. It uses **Redis** for efficient caching and rate limiting.

##  Features

- **OTP Management**: Secure generation and verification of 6-digit numeric codes.
- **TOTP Support**: Google Authenticator compatible multi-factor authentication.
- **QR Code Generation**: Instant QR code delivery for easy TOTP app setup.
- **Redis Integration**: High-performance storage for OTP codes and TOTP secrets.
- **Rate Limiting & Security**: Built-in TTL for OTPs and protection against reuse.
- **Clean Architecture**: Modular structure for easy maintenance and scalability.

##  Technology Stack

- **Language**: Go (v1.25+)
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Caching**: [Redis](https://redis.io/)
- **MFA Library**: [pquerna/otp](https://github.com/pquerna/otp)
- **QR Generator**: [go-qrcode](https://github.com/skip2/go-qrcode)

##  Prerequisites

Before running this project, ensure you have the following installed:
- Go (latest version)
- Redis Server (running on `localhost:6379`)

##  Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Sneha9448/otp-totp-service.git
   cd auth
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```
   The server will start at `http://localhost:8080`.

##  Project Structure

```text
auth/
├── internal/
│   ├── config/       # Redis & App configurations
│   ├── handlers/     # HTTP Request handlers (Gin)
│   ├── models/       # Data structures & JSON models
│   └── services/     # Core business logic (OTP/TOTP)
├── main.go           # Application entry point
└── go.mod            # Project dependencies
```

##  API Documentation

### OTP (One-Time Password)

#### 1. Generate OTP
- **Endpoint**: `POST /api/otp/generate`
- **Payload**:
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Description**: Generates a 6-digit OTP, stores it in Redis with a 5-minute TTL, and returns a success message.

#### 2. Verify OTP
- **Endpoint**: `POST /api/otp/verify`
- **Payload**:
  ```json
  {
    "email": "user@example.com",
    "code": "123456"
  }
  ```

---

### TOTP (Time-based OTP / MFA)

#### 1. Setup TOTP
- **Endpoint**: `POST /api/totp/setup`
- **Payload**:
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Response**: Returns a `secret` and a `qr_code_base64` string to be scanned by Google Authenticator or Authy.

#### 2. Verify TOTP
- **Endpoint**: `POST /api/totp/verify`
- **Payload**:
  ```json
  {
    "email": "user@example.com",
    "code": "654321"
  }
  ```

##  Security Note
This implementation uses Redis to store TOTP secrets for demonstration. In a production environment, TOTP secrets should be stored in a persistent, encrypted database associated with the user profile.


