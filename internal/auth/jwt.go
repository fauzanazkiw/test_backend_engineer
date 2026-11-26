package auth

import (
	"errors"
	"strconv"
	"time"

	"Test/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

// TokenIssuer adalah nilai issuer yang ditanam pada setiap token.
const TokenIssuer = "cinema-api"

// ErrInvalidToken dikembalikan ketika token tidak valid, rusak, atau sudah expired.
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims adalah payload JWT yang dipakai API.
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken membuat access token HS256 untuk user yang berhasil login.
func GenerateToken(secret []byte, expiry time.Duration, user *model.User, now time.Time) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			Issuer:    TokenIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken memvalidasi signature, algorithm, issuer, dan masa berlaku token.
func ParseToken(secret []byte, tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(TokenIssuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
