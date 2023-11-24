package repository

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"job-portal-api/internal/models"
)

func (r *Repo) InsertUser(ctx context.Context, UserDetails models.User) (models.User, error) {
	result := r.DB.Create(&UserDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.User{}, errors.New("failed to create the user")
	}
	return UserDetails, nil
}

func (r *Repo) VerifyUserCredentials(ctx context.Context, email string) (models.User, error) {
	var userDetails models.User
	result := r.DB.Where("email = ?", email).First(&userDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.User{}, errors.New("authentication failed")
	}
	return userDetails, nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var userDetails models.User
	result := r.DB.Where("email = ?", email).First(&userDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.User{}, errors.New("user not found")
	}
	return userDetails, nil
}

// UpdatePassword updates the user's password in the database.
func (r *Repo) UpdatePassword(ctx context.Context, email, hashedPassword string) error {
	// Assuming you have a User model with an Email field
	user := &models.User{}
	err := r.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		log.Printf("Error finding user with email %s: %v", email, err)
		return errors.New("user not found")
	}
	// Update the user's password
	user.PasswordHash = hashedPassword
	err = r.DB.Save(user).Error
	if err != nil {
		log.Printf("Error updating password for user with email %s: %v", email, err)
		return errors.New("failed to update password")
	}

	return nil
}
