package middleware

import (
	"context"
	"net/http"
	"strings"

	"Test/internal/auth"
	"Test/internal/response"
)

type contextKey string

const claimsContextKey contextKey = "auth_claims"

// Auth memvalidasi header Authorization: Bearer <token> sebelum request
// Request tanpa token valid akan ditolak dengan 401.
func Auth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if header == "" {
				response.Error(w, http.StatusUnauthorized, "authorization header is required")
				return
			}

			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			claims, err := auth.ParseToken(secret, parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext mengambil data user yang sudah divalidasi dari context.
func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*auth.Claims)
	return claims, ok
}
