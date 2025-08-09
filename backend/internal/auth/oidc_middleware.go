package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
)

// Middleware è una funzione che implementa un middleware HTTP
type Middleware func(next http.Handler) http.Handler

// OIDC middleware that validates access tokens. If issuer is empty, returns a noop middleware.

type contextKey string

const ProfessionalIDKey contextKey = "professional_id"

// FromContextProfessionalID estrae il ProfessionalID dal context, se presente
func FromContextProfessionalID(ctx context.Context) (uuid.UUID, bool) {
	val := ctx.Value(ProfessionalIDKey)
	if val == nil {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

func NewMiddleware(ctx context.Context, issuer, audience string, skipVerify bool) (Middleware, error) {
	if issuer == "" {
		// No-op middleware (dev mode)
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
		}, nil
	}
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("oidc provider: %w", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: audience})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing auth", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, "invalid auth", http.StatusUnauthorized)
				return
			}
			idToken := parts[1]
			ctx := r.Context()
			if skipVerify {
				// skip verification (dev)
				_ = idToken
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// verify
			token, err := verifier.Verify(ctx, idToken)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			// estrai claim professional_id
			var claims map[string]interface{}
			if err := token.Claims(&claims); err == nil {
				if pid, ok := claims["professional_id"].(string); ok {
					if uuidVal, err := uuid.Parse(pid); err == nil {
						ctx = context.WithValue(ctx, ProfessionalIDKey, uuidVal)
					}
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}
