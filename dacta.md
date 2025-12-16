# API Documentation for Gau Account Service

## Overview
This document provides a detailed description of the APIs available in the Gau Account Service. The service is responsible for managing user accounts, authentication, and related operations. It supports integration with RabbitMQ for messaging and provides endpoints for user registration, login, logout, profile management, and security features.

---

## Base URL
- **Development**: `http://localhost:8080`
- **Production**: `https://<DOMAIN_NAME>`

---

## Endpoints

### 1. User Registration
**POST** `/api/v2/account/basic/register`

#### Description
Registers a new user with the service using basic authentication (email and password).

#### Request Body
```json
{
  "email": "string",
  "password": "string",
  "name": "string"
}
```

#### Response
```json
{
  "message": "User registered successfully."
}
```

---

### 2. User Login
**POST** `/api/v2/account/basic/login`

#### Description
Logs in a user using basic authentication (email and password) and returns a JWT token.

#### Request Body
```json
{
  "email": "string",
  "password": "string"
}
```

#### Response
```json
{
  "token": "string"
}
```

---

### 3. Email Verification
**GET** `/api/v2/account/verify-email/:token`

#### Description
Verifies a user's email address using a token sent via email.

#### Path Parameters
- `token` (required): The verification token.

#### Response
```json
{
  "message": "Email verified successfully."
}
```

---

### 4. Send Email Verification
**POST** `/api/v2/account/send-verification/:user_id`

#### Description
Sends an email verification link to the specified user.

#### Path Parameters
- `user_id` (required): The ID of the user to send the verification email to.

#### Response
```json
{
  "message": "Verification email sent successfully."
}
```

---

### 5. Profile Management
**GET** `/api/v2/account/profile/basic`

#### Description
Fetches the basic profile details of the authenticated user.

#### Headers
- `Authorization`: Bearer `<JWT_TOKEN>`

#### Response
```json
{
  "email": "string",
  "name": "string"
}
```

---

**PUT** `/api/v2/account/profile/basic`

#### Description
Updates the basic profile details of the authenticated user.

#### Headers
- `Authorization`: Bearer `<JWT_TOKEN>`

#### Request Body
```json
{
  "name": "string"
}
```

#### Response
```json
{
  "message": "Profile updated successfully."
}
```

---

### 6. Security Features
**GET** `/api/v2/account/profile/security`

#### Description
Fetches security-related settings and activity logs for the authenticated user.

#### Headers
- `Authorization`: Bearer `<JWT_TOKEN>`

#### Response
```json
{
  "mfa_enabled": true,
  "recent_logins": [
    {
      "ip": "192.168.1.1",
      "location": "New York, USA",
      "timestamp": "2025-12-15T10:00:00Z"
    }
  ]
}
```

---

### 7. OAuth 2.0 Authentication

**POST** `/api/v2/account/sso/google`

#### Description
Logs in a user using Google OAuth 2.0 authentication.

#### Request Body
```json
{
  "idToken": "string" // Google ID token
}
```

#### Response
```json
{
  "token": "string",
  "provider": "google"
}
```

---

### 8. Logout
**POST** `/api/v2/account/logout`

#### Description
Logs out the authenticated user.

#### Headers
- `Authorization`: Bearer `<JWT_TOKEN>`

#### Response
```json
{
  "message": "Logged out successfully."
}
```

---
