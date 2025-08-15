package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/entity"
	"gorm.io/gorm"
	"strings"
)

func (r *Repository) CreateUser(user *entity.User) error {
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
func (r *Repository) CreateUserWithTransaction(tx *gorm.DB, user *entity.User) error {
	if user.AvatarURL == nil || *user.AvatarURL == "" {
		defaultAvatar := "https://cdn.gauas.online/images/avatar/default_image.jpg"
		user.AvatarURL = &defaultAvatar
	}
	if err := tx.Create(user).Error; err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (r *Repository) UpdateUser(user *entity.User) (*entity.User, error) {
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
	err := r.Db.Model(&entity.User{}).Where("user_id = ?", user.UserID).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserWithTransaction updates a user within a transaction
func (r *Repository) UpdateUserWithTransaction(tx *gorm.DB, user *entity.User) (*entity.User, error) {
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
	err := tx.Model(&entity.User{}).Where("user_id = ?", user.UserID).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(id uuid.UUID) error {
	var user entity.User
	if err := r.Db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	if err := r.Db.Delete(&user).Error; err != nil {
		return fmt.Errorf("error deleting user with id %s: %v", id, err)
	}
	return nil
}

func (r *Repository) GetUserById(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.Db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CountUsersByFullName counts users with the same fullname
func (r *Repository) CountUsersByFullName(fullName string) (int64, error) {
	var count int64
	if err := r.Db.Model(&entity.User{}).Where("full_name = ?", fullName).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting users by fullname: %v", err)
	}
	return count, nil
}

// CountUsersByFullNameWithTransaction counts users with the same fullname within a transaction
func (r *Repository) CountUsersByFullNameWithTransaction(tx *gorm.DB, fullName string) (int64, error) {
	var count int64
	if err := tx.Model(&entity.User{}).Where("full_name = ?", fullName).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting users by fullname: %v", err)
	}
	return count, nil
}

func (r *Repository) GetUserByIdentifierAndPassword(identifierType, identifier, hashedPassword string) (*entity.User, error) {
	var user entity.User

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
