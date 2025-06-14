package dto

// AuthResponse represents a generic response for authentication operations
type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// TokenResponse represents a response containing an authentication token
type TokenResponse struct {
	Token string `json:"token"`
}

// NewAuthResponse creates a new AuthResponse with a message
func NewAuthResponse(message string) AuthResponse {
	return AuthResponse{
		Message: message,
	}
}

// NewErrorResponse creates a new AuthResponse with an error message
func NewErrorResponse(errorMsg string) AuthResponse {
	return AuthResponse{
		Error: errorMsg,
	}
}

// NewTokenResponse creates a new TokenResponse with a token
func NewTokenResponse(token string) TokenResponse {
	return TokenResponse{
		Token: token,
	}
}