package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/provider"
	"github.com/tnqbao/gau-account-service/utils"
	"gorm.io/gorm"
)

func (ctrl *Controller) LoginWithGoogle(c *gin.Context) {
	var req ClientRequestGoogleAuthentication
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "invalid request")
		return
	}

	googleUser, err := provider.GetUserInfoFromGoogle(req.Token)
	if err != nil {
		utils.JSON401(c, "invalid Google token")
		return
	}

	email := googleUser.Email
	if !ctrl.IsValidEmail(email) {
		utils.JSON400(c, "invalid email from Google")
		return
	}

	user, err := ctrl.Repository.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userID := uuid.New()
			newUser := &entity.User{
				UserID:     userID,
				Email:      &googleUser.Email,
				FullName:   &googleUser.Name,
				Permission: "member",
			}

			err := ctrl.ExecuteInTransaction(func(tx *gorm.DB) error {
				if googleUser.Name != "" {
					username, err := ctrl.GenerateUsernameFromFullNameWithTransaction(tx, googleUser.Name)
					if err != nil {
						return fmt.Errorf("failed to generate username: %w", err)
					}
					newUser.Username = &username
				}

				if googleUser.Picture != "" {
					imageURL, err := ctrl.UploadAvatarFromURL(userID, googleUser.Picture)
					if err != nil {
						return fmt.Errorf("failed to upload avatar: %w", err)
					}
					// Add CDN URL prefix
					fullImageURL := fmt.Sprintf("%s/images/%s", ctrl.Config.EnvConfig.ExternalService.CDNServiceURL, imageURL)
					newUser.AvatarURL = &fullImageURL
				}

				// Create user within transaction
				if err := ctrl.Repository.CreateUserWithTransaction(tx, newUser); err != nil {
					return fmt.Errorf("failed to create user: %w", err)
				}

				// Create email verification record for Google user
				if googleUser.Email != "" {
					emailVerification := entity.UserVerification{
						ID:         uuid.New(),
						UserID:     userID,
						Method:     "email",
						Value:      googleUser.Email,
						IsVerified: googleUser.EmailVerified,
					}

					// Set verified timestamp if email is already verified by Google
					if emailVerification.IsVerified {
						now := time.Now()
						emailVerification.VerifiedAt = &now
					}

					if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
						return fmt.Errorf("failed to create email verification: %w", err)
					}
				}

				return nil
			})

			if err != nil {
				utils.JSON500(c, "Failed to create user: "+err.Error())
				return
			}

			user = newUser
		} else {
			utils.JSON500(c, "database error")
			return
		}
	} else {
		// User exists - update verification status if Google says email is verified
		if googleUser.EmailVerified {
			err := ctrl.ExecuteInTransaction(func(tx *gorm.DB) error {
				// Get existing email verification record
				verification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(user.UserID, "email", googleUser.Email)
				if err != nil {
					return fmt.Errorf("failed to get email verification: %w", err)
				}

				// Update verification status if not already verified
				if verification != nil && !verification.IsVerified {
					verification.IsVerified = true
					now := time.Now()
					verification.VerifiedAt = &now

					if err := ctrl.Repository.UpdateUserVerification(verification); err != nil {
						return fmt.Errorf("failed to update email verification: %w", err)
					}
				} else if verification == nil {
					// Create verification record if it doesn't exist
					emailVerification := entity.UserVerification{
						ID:         uuid.New(),
						UserID:     user.UserID,
						Method:     "email",
						Value:      googleUser.Email,
						IsVerified: true,
					}
					now := time.Now()
					emailVerification.VerifiedAt = &now

					if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
						return fmt.Errorf("failed to create email verification: %w", err)
					}
				}

				return nil
			})

			if err != nil {
				// Log error but don't fail the login process
				fmt.Printf("Warning: Failed to update email verification for user %s: %v\n", user.UserID, err)
			}
		}
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		utils.JSON400(c, "X-Device-ID header is required")
		return
	}

	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(
		user.UserID,
		user.Permission,
		deviceID,
	)
	if err != nil {
		utils.JSON500(c, "failed to generate token")
		return
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	ctrl.SetAccessCookie(c, accessToken, expiresIn)
	ctrl.SetRefreshCookie(c, refreshToken, 30*24*60*60)

	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}

func (ctrl *Controller) LoginWithFacebook(c *gin.Context) {
	utils.JSON200(c, gin.H{
		"message": "Login with Facebook is not implemented yet",
	})
}
