package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"job-portal-api/internal/middleware"
	"job-portal-api/internal/models"
)

func (h *handler) UserLogin(c *gin.Context) {
	ctx := c.Request.Context()
	traceid, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	var userData models.NewUser

	err := json.NewDecoder(c.Request.Body).Decode(&userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide valid email and password",
		})
		return
	}

	token, err := h.service.UserLoginService(ctx, userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *handler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()

	traceid, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	var userData models.NewUser

	err := json.NewDecoder(c.Request.Body).Decode(&userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide all details",
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide all details",
		})
		return
	}

	userDetails, err := h.service.RegisterUserService(ctx, userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userDetails)

}

func (h *handler) ForgotPasswordHandler(c *gin.Context) {
	ctx := c.Request.Context()
	traceid, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	var forgetPasswordData models.ForgetPasswordRequest

	err := json.NewDecoder(c.Request.Body).Decode(&forgetPasswordData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide a valid email",
		})
		return
	}

	// Call the service method to handle forget password logic
	_, err = h.service.ForgetPasswordService(ctx, forgetPasswordData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset instructions sent to the provided email address.",
	})
}

func (h *handler) UpdatePasswordHandler(c *gin.Context) {
	ctx := c.Request.Context()
	traceid, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	var updatePasswordData models.ResetPasswordRequest
	err := json.NewDecoder(c.Request.Body).Decode(&updatePasswordData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Verify OTP before allowing password update
	err = h.service.VerifyOTPService(ctx, updatePasswordData.Email, updatePasswordData.OTP)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}
	// Update the password in the database
	err = h.service.UpdatePasswordService(ctx, updatePasswordData.Email, updatePasswordData.NewPassword)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// ChangePasswordHandler handles requests to update the password.
func (h *handler) ChangePasswordHandler(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	var changePasswordData models.ChangePasswordRequest
	err := json.NewDecoder(c.Request.Body).Decode(&changePasswordData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Verify old password before allowing the password update
	err = h.service.VerifyOldPassword(ctx, changePasswordData.Email, changePasswordData.OldPassword)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// Update the password in the database
	err = h.service.UpdatePasswordService(ctx, changePasswordData.Email, changePasswordData.NewPassword)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
