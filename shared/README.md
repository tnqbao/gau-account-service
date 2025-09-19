# Shared Module

## Mô tả | Description

**Tiếng Việt:** Module chung chứa các thành phần được chia sẻ giữa HTTP và Consumer services.

**English:** Shared module containing common components used by HTTP and Consumer services.

## Cấu trúc | Structure

```
shared/
├── config/         # Application configuration
├── entity/         # Database entities
├── infra/          # Infrastructure (DB, Logger)
├── provider/       # External service clients
├── repository/     # Data access layer
└── utils/          # Helper functions
```

## Components

### config/
- `main.go` - Config initialization
- `env_config.go` - Environment variables
- `cors.json` - CORS settings

### entity/
- `user.go` - User model
- `user_mfa.go` - MFA settings
- `user_verification.go` - Verification records

### infra/
- `main.go` - Infrastructure setup
- `postgres.go` - Database connection
- `logger.go` - Logging with OpenTelemetry

### provider/
- `authorization_service.go` - Auth service client
- `upload_service.go` - Upload service client
- `sso.go` - OAuth providers

### repository/
- `user.go` - User data operations
- `verification_code.go` - Verification operations

### utils/
- `jwt.go` - JWT utilities
- `response.go` - HTTP response helpers
- `coalesce.go` - Null handling

## Configuration

### Environment Variables
```bash
# Database
PGPOOL_HOST=localhost
PGPOOL_DB=gau_account
PGPOOL_USER=postgres
PGPOOL_PASSWORD=password

# JWT
JWT_SECRET_KEY=your-secret-key
JWT_EXPIRE=604800

# External Services
AUTHORIZATION_SERVICE_URL=http://localhost:8080
UPLOAD_SERVICE_URL=http://localhost:8081
```

## Usage

```go
// Initialize
cfg := config.NewConfig()
infra := infra.InitInfra(cfg)
repo := repository.InitRepository(infra)

// Database operations
user := &entity.User{Username: &username}
err := repo.User.Create(user)

// JWT operations
token, err := utils.GenerateJWT(userID, cfg.JWT.SecretKey)
```
