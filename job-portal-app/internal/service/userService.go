package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"job-portal-api/internal/models"
	"strconv"
	"time"
)

func (s *Service) UserLoginService(ctx context.Context, userData models.NewUser) (string, error) {
	// Checking the email in the db
	var userDetails models.User
	userDetails, err := s.UserRepo.VerifyUserCredentials(ctx, userData.Email)
	if err != nil {
		log.Info().Err(err).Msg("Failed to verify user credentials")
		return "", fmt.Errorf("failed to verify user credentials: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDetails.PasswordHash), []byte(userData.Password))
	if err != nil {
		log.Info().Err(err).Msg("Invalid password provided")
		return "", errors.New("invalid password provided")
	}
	claims := jwt.RegisteredClaims{
		Issuer:    "job portal project",
		Subject:   strconv.FormatUint(uint64(userDetails.ID), 10),
		Audience:  jwt.ClaimStrings{"users"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token, err := s.auth.GenerateAuthToken(claims)
	if err != nil {
		log.Info().Err(err).Msg("Failed to generate authentication token")
		return "", fmt.Errorf("failed to generate authentication token: %w", err)
	}
	return token, nil
}

func (s *Service) RegisterUserService(ctx context.Context, userData models.NewUser) (models.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("generating password hash: %w", err)
	}
	userDetails := models.User{
		Username:     userData.Username,
		Email:        userData.Email,
		PasswordHash: string(hashedPass),
	}
	userDetails, err = s.UserRepo.InsertUser(ctx, userDetails)
	if err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}
