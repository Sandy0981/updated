package models

import "gorm.io/gorm"

type NewUser struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique"`
	Email        string `json:"email" gorm:"unique"`
	PasswordHash string `json:"-"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
}

type ForgetPasswordResponse struct {
	Message string `json:"message"`
}

// ResetPasswordRequest represents the structure of the request for OTP-based password reset.
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	OTP         string `json:"otp"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
	// You can include additional fields if needed, such as user information or error details.
}
