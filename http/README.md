# HTTP Module

## Mô tả | Description

**Tiếng Việt:** Module HTTP xử lý REST API requests cho Gau Account Service.

**English:** HTTP module handling REST API requests for Gau Account Service.

## Cấu trúc | Structure

```
http/
├── controller/     # Request handlers
├── middlewares/    # HTTP middlewares
└── routes/         # Route definitions
```

## Components

### controller/
- `main.go` - Controller initialization
- `register.go` - User registration
- `login.go` - User authentication
- `profile.go` - Profile management
- `mfa.go` - Multi-factor authentication
- `dto.go` - Data transfer objects
- `helper.go` - Helper functions

### middlewares/
- `main.go` - Middleware setup
- `jwt.go` - JWT authentication
- `cors.go` - CORS configuration

### routes/
- `routes.go` - API route definitions

## API Endpoints

### Authentication
```
POST /api/v2/account/basic/register    # Register
POST /api/v2/account/basic/login       # Login
```

### Profile
```
GET  /api/v2/account/profile/basic     # Get basic info
PUT  /api/v2/account/profile/basic     # Update basic info
GET  /api/v2/account/profile/security  # Get security info
PUT  /api/v2/account/profile/security  # Update security info
```

### MFA
```
GET  /api/v2/account/mfa/totp/qr       # Generate TOTP QR
POST /api/v2/account/mfa/totp/enable   # Enable TOTP
POST /api/v2/account/mfa/totp/verify   # Verify TOTP
```

## Usage

```go
// Initialize controller
ctrl := controller.NewController(config, infra)

// Setup routes
router := routes.SetupRouter(config)
router.Run(":8080")
```
