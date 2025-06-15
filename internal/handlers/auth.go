package handlers

import (
	"errors"
	"net/http"

	"github.com/dinis/musync/internal/config"
	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/dto"
	"github.com/dinis/musync/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authCfg config.AuthConfig, emailCfg config.EmailConfig) *AuthHandler {
	authService := services.NewAuthService(database.GlobalDB, authCfg, emailCfg)
	return &AuthHandler{
		authService: authService,
	}
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ConfirmResetRequest struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type SetPasswordRequest struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	err := h.authService.SignUp(c.Request.Context(), req.Email, req.Username)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, dto.NewErrorResponse("User already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to create user"))
		return
	}

	c.JSON(http.StatusCreated, dto.NewAuthResponse("User created successfully. Please check your email for verification and to set your password."))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("Invalid credentials"))
		case errors.Is(err, services.ErrEmailNotVerified):
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("Please verify your email first"))
		case errors.Is(err, services.ErrPasswordNotSet):
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("Please set your password first"))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to login"))
		}
		return
	}

	c.JSON(http.StatusOK, dto.NewTokenResponse(token))
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Verification code is required"))
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), code)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrVerificationCodeRequired):
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Verification code is required"))
		case errors.Is(err, services.ErrInvalidVerificationCode):
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid verification code"))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to verify email"))
		}
		return
	}

	c.JSON(http.StatusOK, dto.NewAuthResponse("Email verified successfully"))
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		// Don't reveal if the error is because the user doesn't exist
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to process reset request"))
		return
	}

	c.JSON(http.StatusOK, dto.NewAuthResponse("If your email is registered, you will receive a password reset link"))
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ConfirmResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Code, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidResetCode):
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid or expired reset code"))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to reset password"))
		}
		return
	}

	c.JSON(http.StatusOK, dto.NewAuthResponse("Password reset successfully"))
}

// SetPassword handles setting a password during email verification
func (h *AuthHandler) SetPassword(c *gin.Context) {
	var req SetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	err := h.authService.SetPassword(c.Request.Context(), req.Code, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidVerificationCode):
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid verification code"))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to set password"))
		}
		return
	}

	c.JSON(http.StatusOK, dto.NewAuthResponse("Password set successfully. You can now log in."))
}
