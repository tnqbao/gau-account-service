package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	entity "github.com/tnqbao/gau-account-service/shared/entity"
	"github.com/tnqbao/gau-account-service/shared/provider"
	"github.com/tnqbao/gau-account-service/shared/utils"
	"gorm.io/gorm"
)

func (ctrl *Controller) LoginWithGoogle(c *gin.Context) {
	ctx := c.Request.Context()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Google authentication request received")

	var req ClientRequestGoogleAuthentication
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to bind JSON request")
		utils.JSON400(c, "invalid request")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Validating Google token")

	googleUser, err := provider.GetUserInfoFromGoogle(req.Token)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Invalid Google token provided")
		utils.JSON401(c, "invalid Google token")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Google user info retrieved - Email: %s, Name: %s, Verified: %v",
		googleUser.Email, googleUser.Name, googleUser.EmailVerified)

	email := googleUser.Email
	if !ctrl.IsValidEmail(email) {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Google Login] Invalid email format from Google: %s", email)
		utils.JSON400(c, "invalid email from Google")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Checking if user exists with email: %s", email)

	user, err := ctrl.Repository.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] User not found, creating new user for email: %s", email)

			userID := uuid.New()
			newUser := &entity.User{
				UserID:     userID,
				Email:      &googleUser.Email,
				FullName:   &googleUser.Name,
				Permission: "member",
			}

			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Starting user creation transaction for new user: %s", userID.String())

			err := ctrl.ExecuteInTransaction(func(tx *gorm.DB) error {
				if googleUser.Name != "" {
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Generating username from full name for user: %s", userID.String())
					username, err := ctrl.GenerateUsernameFromFullNameWithTransaction(tx, googleUser.Name)
					if err != nil {
						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to generate username for user: %s", userID.String())
						return fmt.Errorf("failed to generate username: %w", err)
					}
					newUser.Username = &username
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Username generated for user: %s", userID.String())
				}

				if googleUser.Picture != "" {
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Uploading avatar from Google for user: %s", userID.String())
					imageURL, err := ctrl.UploadAvatarFromURL(userID, googleUser.Picture)
					if err != nil {
						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to upload avatar for user: %s", userID.String())
						return fmt.Errorf("failed to upload avatar: %w", err)
					}
					// Add CDN URL prefix
					fullImageURL := fmt.Sprintf("%s/images/%s", ctrl.Config.EnvConfig.ExternalService.CDNServiceURL, imageURL)
					newUser.AvatarURL = &fullImageURL
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Avatar uploaded successfully for user: %s", userID.String())
				}

				// Create user within transaction
				if err := ctrl.Repository.CreateUserWithTransaction(tx, newUser); err != nil {
					ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to create user in database: %s", userID.String())
					return fmt.Errorf("failed to create user: %w", err)
				}

				ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] User created successfully: %s", userID.String())

				// Create email verification record for Google user
				if googleUser.Email != "" {
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Creating email verification record for user: %s - Verified: %v", userID.String(), googleUser.EmailVerified)

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
						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to create email verification for user: %s", userID.String())
						return fmt.Errorf("failed to create email verification: %w", err)
					}

					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Email verification record created for user: %s", userID.String())
				}

				return nil
			})

			if err != nil {
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Transaction failed for new user creation")
				utils.JSON500(c, "Failed to create user: "+err.Error())
				return
			}

			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] New user creation completed successfully: %s", userID.String())
			user = newUser
		} else {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Database error while checking user existence")
			utils.JSON500(c, "database error")
			return
		}
	} else {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Existing user found: %s", user.UserID.String())

		// User exists - update verification status if Google says email is verified
		if googleUser.EmailVerified {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Updating email verification status for existing user: %s", user.UserID.String())

			err := ctrl.ExecuteInTransaction(func(tx *gorm.DB) error {
				// Get existing email verification record
				verification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(user.UserID, "email", googleUser.Email)
				if err != nil {
					ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to get email verification for user: %s", user.UserID.String())
					return fmt.Errorf("failed to get email verification: %w", err)
				}

				// Update verification status if not already verified
				if verification != nil && !verification.IsVerified {
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Updating unverified email to verified for user: %s", user.UserID.String())
					verification.IsVerified = true
					now := time.Now()
					verification.VerifiedAt = &now

					if err := ctrl.Repository.UpdateUserVerification(verification); err != nil {
						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to update email verification for user: %s", user.UserID.String())
						return fmt.Errorf("failed to update email verification: %w", err)
					}
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Email verification updated for user: %s", user.UserID.String())
				} else if verification == nil {
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Creating missing email verification record for user: %s", user.UserID.String())
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
						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to create email verification for user: %s", user.UserID.String())
						return fmt.Errorf("failed to create email verification: %w", err)
					}
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Email verification record created for user: %s", user.UserID.String())
				} else {
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Email already verified for user: %s", user.UserID.String())
				}

				return nil
			})

			if err != nil {
				// Log error but don't fail the login process
				ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Google Login] Failed to update email verification for user %s: %v", user.UserID.String(), err)
			}
		}
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Google Login] Missing device ID for user: %s", user.UserID.String())
		utils.JSON400(c, "X-Device-ID header is required")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Creating tokens for user: %s with device: %s", user.UserID.String(), deviceID)

	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(user.UserID, user.Permission, deviceID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Google Login] Failed to create tokens for user: %s", user.UserID.String())
		utils.JSON500(c, "Failed to create authentication tokens")
		return
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Tokens created successfully for user: %s", user.UserID.String())

	ctrl.SetAccessCookie(c, accessToken, expiresIn)
	ctrl.SetRefreshCookie(c, refreshToken, 30*24*60*60)

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Google Login] Google login completed successfully for user: %s", user.UserID.String())

	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}

//func (ctrl *Controller) LoginWithFacebook(c *gin.Context) {
//	ctx := c.Request.Context()
//
//	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Facebook authentication request received")
//
//	var req ClientRequestFacebookAuthentication
//	if err := c.ShouldBindJSON(&req); err != nil {
//		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Failed to bind JSON request")
//		utils.JSON400(c, "invalid request")
//		return
//	}
//
//	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Validating Facebook token")
//
//	facebookUser, err := provider.GetUserInfoFromFacebook(req.Token)
//	if err != nil {
//		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Invalid Facebook token provided")
//		utils.JSON401(c, "invalid Facebook token")
//		return
//	}
//
//	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Facebook user info retrieved - Email: %s, Name: %s", facebookUser.Email, facebookUser.Name)
//
//	email := facebookUser.Email
//	if !ctrl.IsValidEmail(email) {
//		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Facebook Login] Invalid email format from Facebook: %s", email)
//		utils.JSON400(c, "invalid email from Facebook")
//		return
//	}
//
//	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Checking if user exists with email: %s", email)
//
//	user, err := ctrl.Repository.GetUserByEmail(email)
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] User not found, creating new user for email: %s", email)
//
//			userID := uuid.New()
//			newUser := &entity.User{
//				UserID:     userID,
//				Email:      &facebookUser.Email,
//				FullName:   &facebookUser.Name,
//				Permission: "member",
//			}
//
//			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Starting user creation transaction for new user: %s", userID.String())
//
//			err := ctrl.ExecuteInTransaction(func(tx *gorm.DB) error {
//				if facebookUser.Name != "" {
//					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Generating username from full name for user: %s", userID.String())
//					username, err := ctrl.GenerateUsernameFromFullNameWithTransaction(tx, facebookUser.Name)
//					if err != nil {
//						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Failed to generate username for user: %s", userID.String())
//						return fmt.Errorf("failed to generate username: %w", err)
//					}
//					newUser.Username = &username
//					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Username generated for user: %s", userID.String())
//				}
//
//				if facebookUser.Picture.Data.URL != "" {
//					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Uploading avatar from Facebook for user: %s", userID.String())
//					imageURL, err := ctrl.UploadAvatarFromURL(userID, facebookUser.Picture.Data.URL)
//					if err != nil {
//						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Failed to upload avatar for user: %s", userID.String())
//						return fmt.Errorf("failed to upload avatar: %w", err)
//					}
//					// Add CDN URL prefix
//					fullImageURL := fmt.Sprintf("%s/images/%s", ctrl.Config.EnvConfig.ExternalService.CDNServiceURL, imageURL)
//					newUser.AvatarURL = &fullImageURL
//					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Avatar uploaded successfully for user: %s", userID.String())
//				}
//
//				// Create user within transaction
//				if err := ctrl.Repository.CreateUserWithTransaction(tx, newUser); err != nil {
//					ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Failed to create user in database: %s", userID.String())
//					return fmt.Errorf("failed to create user: %w", err)
//				}
//
//				ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] User created successfully: %s", userID.String())
//
//				// Create email verification record for Facebook user
//				if facebookUser.Email != "" {
//					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Creating email verification record for user: %s", userID.String())
//
//					emailVerification := entity.UserVerification{
//						ID:         uuid.New(),
//						UserID:     userID,
//						Method:     "email",
//						Value:      facebookUser.Email,
//						IsVerified: true, // Facebook emails are considered verified
//					}
//					now := time.Now()
//					emailVerification.VerifiedAt = &now
//
//					if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
//						ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Failed to create email verification for user: %s", userID.String())
//						return fmt.Errorf("failed to create email verification: %w", err)
//					}
//
//					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Email verification record created for user: %s", userID.String())
//				}
//
//				return nil
//			})
//
//			if err != nil {
//				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Transaction failed for new user creation")
//				utils.JSON500(c, "Failed to create user: "+err.Error())
//				return
//			}
//
//			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] New user creation completed successfully: %s", userID.String())
//			user = newUser
//		} else {
//			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Database error while checking user existence")
//			utils.JSON500(c, "database error")
//			return
//		}
//	} else {
//		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Existing user found: %s", user.UserID.String())
//	}
//
//	deviceID := c.GetHeader("X-Device-ID")
//	if deviceID == "" {
//		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Facebook Login] Missing device ID for user: %s", user.UserID.String())
//		utils.JSON400(c, "X-Device-ID header is required")
//		return
//	}
//
//	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Creating tokens for user: %s with device: %s", user.UserID.String(), deviceID)
//
//	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(user.UserID.String(), user.Permission, deviceID)
//	if err != nil {
//		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Facebook Login] Failed to create tokens for user: %s", user.UserID.String())
//		utils.JSON500(c, "Failed to create authentication tokens")
//		return
//	}
//
//	expiresIn := int(time.Until(expiresAt).Seconds())
//
//	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Tokens created successfully for user: %s", user.UserID.String())
//
//	ctrl.SetAccessCookie(c, accessToken, expiresIn)
//	ctrl.SetRefreshCookie(c, refreshToken, 30*24*60*60)
//
//	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Facebook Login] Facebook login completed successfully for user: %s", user.UserID.String())
//
//	utils.JSON200(c, gin.H{
//		"access_token":  accessToken,
//		"refresh_token": refreshToken,
//		"expires_in":    expiresIn,
//	})
//}
