package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/schemas"
	"strings"
)

func (r *Repository) CreateUser(user *schemas.User) error {
	if err := r.db.Omit("image_url").Create(user).Error; err != nil {
		return fmt.Errorf("error creating user credential: %v", err)
	}
	return nil
}

func (r *Repository) UpdateUser(user *schemas.User) (*schemas.User, error) {
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
	if user.ImageURL != nil {
		data["image_url"] = user.ImageURL
	}

	err := r.db.Model(&schemas.User{}).Where("user_id = ?", user.UserID).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(id uuid.UUID) error {
	var user schemas.User
	if err := r.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	if err := r.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("error deleting user with id %s: %v", id, err)
	}
	return nil
}

func (r *Repository) GetUserById(id uuid.UUID) (*schemas.User, error) {
	var user schemas.User
	if err := r.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*schemas.User, error) {
	var user schemas.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByIdentifierAndPassword(identifierType, identifier, hashedPassword string) (*schemas.User, error) {
	var user schemas.User

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

	if err := r.db.Where(fmt.Sprintf("%s = ? AND password = ?", queryField), identifier, hashedPassword).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found with %s and password: %v", queryField, err)
	}

	return &user, nil
}
