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
	if err := r.db.Save(user).Error; err != nil {
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
