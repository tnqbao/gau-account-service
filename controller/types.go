package controller

import (
	"time"

	"github.com/google/uuid"
)

// Basic login request structure
type ClientRequestBasicLogin struct {
	Username  *string `json:"username"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Password  *string `json:"password"`
	KeepLogin *string `json:"keepMeLogin"`
}

// Basic registration request structure
type UserBasicRegistryReq struct {
	Username    *string   `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	FullName    string    `json:"fullname,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

// User information response structure (basic info only)
type UserBasicInfoResponse struct {
	UserId      uuid.UUID  `json:"user_id"`
	FullName    string     `json:"fullname,omitempty"`
	Email       string     `json:"email,omitempty"`
	Phone       string     `json:"phone,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	AvatarURL   string     `json:"avatar_url,omitempty"`
	GithubUrl   string     `json:"github_url,omitempty"`
	FacebookUrl string     `json:"facebook_url,omitempty"`
	Username    string     `json:"username,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	Permission  string     `json:"permission,omitempty"`
}

// User security information response structure (verification & MFA only)
type UserSecurityInfoResponse struct {
	UserId          uuid.UUID              `json:"user_id"`
	IsEmailVerified bool                   `json:"is_email_verified"`
	IsPhoneVerified bool                   `json:"is_phone_verified"`
	Verifications   []UserVerificationInfo `json:"verifications,omitempty"`
	MFAs            []UserMFAInfo          `json:"mfas,omitempty"`
}

// User complete information response structure (all info)
type UserCompleteInfoResponse struct {
	UserId          uuid.UUID              `json:"user_id"`
	FullName        string                 `json:"fullname,omitempty"`
	Email           string                 `json:"email,omitempty"`
	Phone           string                 `json:"phone,omitempty"`
	DateOfBirth     *time.Time             `json:"date_of_birth,omitempty"`
	AvatarURL       string                 `json:"avatar_url,omitempty"`
	GithubUrl       string                 `json:"github_url,omitempty"`
	FacebookUrl     string                 `json:"facebook_url,omitempty"`
	Username        string                 `json:"username,omitempty"`
	Gender          string                 `json:"gender,omitempty"`
	Permission      string                 `json:"permission,omitempty"`
	IsEmailVerified bool                   `json:"is_email_verified"`
	IsPhoneVerified bool                   `json:"is_phone_verified"`
	Verifications   []UserVerificationInfo `json:"verifications,omitempty"`
	MFAs            []UserMFAInfo          `json:"mfas,omitempty"`
}

// Legacy type for backward compatibility
type UserInfoResponse = UserBasicInfoResponse
type UserDetailedInfoResponse = UserCompleteInfoResponse

// User verification information for response
type UserVerificationInfo struct {
	ID         uuid.UUID  `json:"id"`
	Method     string     `json:"method"`
	Value      string     `json:"value"`
	IsVerified bool       `json:"is_verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// User MFA information for response
type UserMFAInfo struct {
	ID         uuid.UUID  `json:"id"`
	Type       string     `json:"type"`
	Enabled    bool       `json:"enabled"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// User information update request structure
type UserInformationUpdateReq struct {
	FullName    *string    `json:"fullname,omitempty"`
	Username    *string    `json:"username,omitempty"`
	Email       *string    `json:"email,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	FacebookURL *string    `json:"facebook_url,omitempty"`
	GitHubURL   *string    `json:"github_url,omitempty"`
}

// User basic information update request structure (no security data)
type UserBasicInfoUpdateReq struct {
	FullName    *string    `json:"fullname,omitempty"`
	Username    *string    `json:"username,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	FacebookURL *string    `json:"facebook_url,omitempty"`
	GitHubURL   *string    `json:"github_url,omitempty"`
}

// User security information update request structure (email/phone changes)
type UserSecurityInfoUpdateReq struct {
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

// User complete information update request structure (all fields)
type UserCompleteInfoUpdateReq struct {
	FullName    *string    `json:"fullname,omitempty"`
	Username    *string    `json:"username,omitempty"`
	Email       *string    `json:"email,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	FacebookURL *string    `json:"facebook_url,omitempty"`
	GitHubURL   *string    `json:"github_url,omitempty"`
}

// Client request structure for google login
type ClientRequestGoogleAuthentication struct {
	Token string `json:"token" binding:"required"`
}
