package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"Test/internal/auth"
	"Test/internal/domain"
	"Test/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// AuthService menangani proses autentikasi user.
type AuthService struct {
	users  *repository.UserRepository
	secret []byte
	expiry time.Duration
	now    func() time.Time
}

func NewAuthService(users *repository.UserRepository, secret string, expiry time.Duration) *AuthService {
	return &AuthService{
		users:  users,
		secret: []byte(secret),
		expiry: expiry,
		now:    time.Now,
	}
}

// LoginResult adalah credential yang dikembalikan setelah login berhasil.
type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
}

// Login memvalidasi credential dan menghasilkan access token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := s.users.GetByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(s.secret, s.expiry, user, s.now())
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.expiry.Seconds()),
	}, nil
}
