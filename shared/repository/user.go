package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	entity2 "github.com/tnqbao/gau-account-service/shared/entity"
	"gorm.io/gorm"
)

func (r *Repository) CreateUser(user *entity2.User) error {
	if user.AvatarURL == nil || *user.AvatarURL == "" {
		defaultAvatar := "https://cdn.gauas.online/images/avatar/default_image.jpg"
		user.AvatarURL = &defaultAvatar
	}
	if err := r.Db.Create(user).Error; err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

// CreateUserWithTransaction creates a user within a transaction
func (r *Repository) CreateUserWithTransaction(tx *gorm.DB, user *entity2.User) error {
	if user.AvatarURL == nil || *user.AvatarURL == "" {
		defaultAvatar := "https://cdn.gauas.online/images/avatar/default_image.jpg"
		user.AvatarURL = &defaultAvatar
	}
	if err := tx.Create(user).Error; err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (r *Repository) UpdateUser(user *entity2.User) (*entity2.User, error) {
	data := map[string]interface{}{}

	if user.Username != nil {
		data["username"] = user.Username
	}
	if user.FullName != nil {
		data["full_name"] = user.FullName
	}
	if user.Email != nil {
		data["email"] = user.Email
	}
	if user.Phone != nil {
		data["phone"] = user.Phone
	}
	if user.DateOfBirth != nil {
		data["date_of_birth"] = user.DateOfBirth
	}
	if user.Gender != nil {
		data["gender"] = user.Gender
	}
	if user.FacebookURL != nil {
		data["facebook_url"] = user.FacebookURL
	}
	if user.GithubURL != nil {
		data["github_url"] = user.GithubURL
	}
	if user.AvatarURL != nil {
		data["avatar_url"] = user.AvatarURL
	}

	// Update the user in the database
	err := r.Db.Model(&entity2.User{}).Where("user_id = ?", user.UserID).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserWithTransaction updates a user within a transaction
func (r *Repository) UpdateUserWithTransaction(tx *gorm.DB, user *entity2.User) (*entity2.User, error) {
	data := map[string]interface{}{}

	if user.Username != nil {
		data["username"] = user.Username
	}
	if user.FullName != nil {
		data["full_name"] = user.FullName
	}
	if user.Email != nil {
		data["email"] = user.Email
	}
	if user.Phone != nil {
		data["phone"] = user.Phone
	}
	if user.DateOfBirth != nil {
		data["date_of_birth"] = user.DateOfBirth
	}
	if user.Gender != nil {
		data["gender"] = user.Gender
	}
	if user.FacebookURL != nil {
		data["facebook_url"] = user.FacebookURL
	}
	if user.GithubURL != nil {
		data["github_url"] = user.GithubURL
	}
	if user.AvatarURL != nil {
		data["avatar_url"] = user.AvatarURL
	}

	// Update the user in the database
	err := tx.Model(&entity2.User{}).Where("user_id = ?", user.UserID).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(id uuid.UUID) error {
	var user entity2.User
	if err := r.Db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	if err := r.Db.Delete(&user).Error; err != nil {
		return fmt.Errorf("error deleting user with id %s: %v", id, err)
	}
	return nil
}

func (r *Repository) GetUserById(id uuid.UUID) (*entity2.User, error) {
	var user entity2.User
	if err := r.Db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*entity2.User, error) {
	var user entity2.User
	if err := r.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CountUsersByFullName counts users with the same fullname
func (r *Repository) CountUsersByUsername(fullName string) (int64, error) {
	var count int64
	if err := r.Db.Model(&entity2.User{}).Where("username = ?", fullName).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting users by fullname: %v", err)
	}
	return count, nil
}

// CountUsersByUsernameWithTransaction counts users with usernames starting with the base username and returns all matching usernames
func (r *Repository) CountUsersByUsernameWithTransaction(tx *gorm.DB, baseUsername string) (int64, error) {
	var count int64
	// Count all usernames that start with baseUsername (e.g., "nguyenvana", "nguyenvana1", "nguyenvana2")
	if err := tx.Model(&entity2.User{}).Where("username LIKE ?", baseUsername+"%").Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting users by username: %v", err)
	}
	return count, nil
}

// GetUsernamesStartingWithTransaction gets all usernames starting with base username within a transaction
func (r *Repository) GetUsernamesStartingWithTransaction(tx *gorm.DB, baseUsername string) ([]string, error) {
	var usernames []string
	if err := tx.Model(&entity2.User{}).Where("username LIKE ?", baseUsername+"%").Pluck("username", &usernames).Error; err != nil {
		return nil, fmt.Errorf("error getting usernames: %v", err)
	}
	return usernames, nil
}

func (r *Repository) GetUserByIdentifierAndPassword(identifierType, identifier, hashedPassword string) (*entity2.User, error) {
	var user entity2.User

	var queryField string
	switch strings.ToLower(identifierType) {
	case "email":
		queryField = "email"
	case "phone":
		queryField = "phone"
	case "username":
		queryField = "username"
	default:
		return nil, fmt.Errorf("invalid identifier type: %s", identifierType)
	}

	if err := r.Db.Where(fmt.Sprintf("%s = ? AND password = ?", queryField), identifier, hashedPassword).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found with %s and password: %v", queryField, err)
	}

	return &user, nil
}

// CreateUserVerification creates a new user verification record
func (r *Repository) CreateUserVerification(verification *entity2.UserVerification) error {
	if err := r.Db.Create(verification).Error; err != nil {
		return fmt.Errorf("error creating user verification: %v", err)
	}
	return nil
}

// CreateUserVerificationWithTransaction creates a user verification within a transaction
func (r *Repository) CreateUserVerificationWithTransaction(tx *gorm.DB, verification *entity2.UserVerification) error {
	if err := tx.Create(verification).Error; err != nil {
		return fmt.Errorf("error creating user verification: %v", err)
	}
	return nil
}

// GetUserVerifications gets all verification records for a user
func (r *Repository) GetUserVerifications(userID uuid.UUID) ([]entity2.UserVerification, error) {
	var verifications []entity2.UserVerification
	if err := r.Db.Where("user_id = ?", userID).Find(&verifications).Error; err != nil {
		return nil, fmt.Errorf("error getting user verifications: %v", err)
	}
	return verifications, nil
}

// GetUserVerificationByMethodAndValue gets a specific verification record
func (r *Repository) GetUserVerificationByMethodAndValue(userID uuid.UUID, method, value string) (*entity2.UserVerification, error) {
	var verification entity2.UserVerification
	if err := r.Db.Where("user_id = ? AND method = ? AND value = ?", userID, method, value).First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user verification: %v", err)
	}
	return &verification, nil
}

// UpdateUserVerification updates a verification record
func (r *Repository) UpdateUserVerification(verification *entity2.UserVerification) error {
	if err := r.Db.Save(verification).Error; err != nil {
		return fmt.Errorf("error updating user verification: %v", err)
	}
	return nil
}

// UpdateEmailVerificationStatus updates the email verification status to verified
func (r *Repository) UpdateEmailVerificationStatus(userID uuid.UUID, email string) error {
	now := time.Now()
	if err := r.Db.Model(&entity2.UserVerification{}).
		Where("user_id = ? AND method = ? AND value = ?", userID, "email", email).
		Updates(map[string]interface{}{
			"is_verified": true,
			"verified_at": now,
		}).Error; err != nil {
		return fmt.Errorf("error updating email verification status: %v", err)
	}
	return nil
}

// CreateUserMFA creates a new MFA record
func (r *Repository) CreateUserMFA(mfa *entity2.UserMFA) error {
	if err := r.Db.Create(mfa).Error; err != nil {
		return fmt.Errorf("error creating user MFA: %v", err)
	}
	return nil
}

// GetUserMFAs gets all MFA records for a user
func (r *Repository) GetUserMFAs(userID uuid.UUID) ([]entity2.UserMFA, error) {
	var mfas []entity2.UserMFA
	if err := r.Db.Where("user_id = ?", userID).Find(&mfas).Error; err != nil {
		return nil, fmt.Errorf("error getting user MFAs: %v", err)
	}
	return mfas, nil
}

// GetUserMFAByType gets a specific MFA record by type
func (r *Repository) GetUserMFAByType(userID uuid.UUID, mfaType string) (*entity2.UserMFA, error) {
	var mfa entity2.UserMFA
	if err := r.Db.Where("user_id = ? AND type = ?", userID, mfaType).First(&mfa).Error; err != nil {
		return nil, err
	}
	return &mfa, nil
}

// UpdateUserMFA updates an MFA record
func (r *Repository) UpdateUserMFA(mfa *entity2.UserMFA) error {
	if err := r.Db.Save(mfa).Error; err != nil {
		return fmt.Errorf("error updating user MFA: %v", err)
	}
	return nil
}
