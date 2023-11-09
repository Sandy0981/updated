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
