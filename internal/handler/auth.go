package handler

import (
	"errors"
	"log"
	"net/http"

	"Test/internal/domain"
	"Test/internal/response"
	"Test/internal/service"
)

// AuthHandler menangani endpoint autentikasi.
type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: authService}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Login menangani POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		log.Printf("auth handler: login gagal: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.MessageData(w, http.StatusOK, "login successful", loginResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresIn:   result.ExpiresIn,
	})
}
