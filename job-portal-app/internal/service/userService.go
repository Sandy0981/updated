package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"job-portal-api/internal/config"
	"job-portal-api/internal/models"
	"math/rand"
	"net/smtp"
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

func (s *Service) ForgetPasswordService(ctx context.Context, data models.ForgetPasswordRequest) (models.ForgetPasswordResponse, error) {
	// Check if the user exists
	_, err := s.UserRepo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return models.ForgetPasswordResponse{}, errors.New("user not found")
	}

	// Generate OTP
	otp := generateOTP()

	// Save OTP and expiration in Redis cache
	err = s.rdb.SaveOTPInCache(data.Email, otp)
	if err != nil {
		log.Printf("Error saving OTP in cache: %v", err)
		return models.ForgetPasswordResponse{}, errors.New("failed to generate OTP")
	}

	// Send OTP via email
	err = sendOTPEmail(data.Email, otp)
	if err != nil {
		log.Printf("Error sending OTP email: %v", err)
		return models.ForgetPasswordResponse{}, errors.New("failed to send OTP via email")
	}

	return models.ForgetPasswordResponse{Message: "OTP sent to the registered email address"}, nil
}

func generateOTP() string {
	return strconv.Itoa(rand.Intn(999999))
}

func sendOTPEmail(email, otp string) error {
	cfg := config.GetConfig()
	from := "sandeepsinghs321@gmail.com"
	password := "latc ymgz ksxv zuzc"

	// Recipient's email address
	to := email

	// SMTP server details
	smtpServer := "smtp.gmail.com"
	smtpPort := cfg.MailConfig.Port

	// Message content
	message := []byte(fmt.Sprintf("Subject: Forget Password OTP\n\nYour OTP for forget password: %s", otp))

	// Authentication information
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// SMTP connection
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Print("Email sent successfully!")
	return nil
}

func (s *Service) VerifyOTPService(ctx context.Context, email string, otp string) error {
	storedOTP, err := s.rdb.GetOTP(email)
	if err != nil {
		return errors.New("failed to retrieve OTP from the database")
	}

	// Check if the OTP matches the stored OTP
	if otp != storedOTP {
		return errors.New("invalid OTP")
	}

	// Delete the OTP from the database (optional, depending on your use case)
	err = s.rdb.DeleteOTP(email)
	if err != nil {
		return errors.New("failed to delete OTP from the database")
	}

	// OTP verification successful
	return nil
}

func (s *Service) UpdatePasswordService(ctx context.Context, email, newPassword string) error {
	// Hash the new password before updating
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing the new password for email %s: %v", email, err)
		return errors.New("failed to hash the new password")
	}

	// Update the password in the database
	err = s.UserRepo.UpdatePassword(ctx, email, string(hashedPassword))
	if err != nil {
		log.Printf("Error updating the password for email %s: %v", email, err)
		return fmt.Errorf("failed to update the password: %w", err)
	}

	log.Printf("Password updated successfully for email %s", email)
	return nil
}

func (s *Service) VerifyOldPassword(ctx context.Context, email string, oldPass string) error {
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("failed to retrieve user by email")
		return err
	}

	// Compare the oldPass with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPass))
	if err != nil {
		// Passwords do not match
		log.Warn().Str("email", email).Msg("invalid old password")
		return errors.New("invalid old password")
	}

	// Old password is valid
	log.Info().Str("email", email).Msg("old password verification successful")
	return nil
}
