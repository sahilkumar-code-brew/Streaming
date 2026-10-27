package session

import (
	"context"
	"crypto/ed25519"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/apperror"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/httpx"
)

const bearerPrefix = "Bearer "

type userIDKey struct{}

type Verifier struct {
	publicKey ed25519.PublicKey
}

func NewVerifier(publicKey ed25519.PublicKey) *Verifier {
	return &Verifier{publicKey: publicKey}
}

func (v *Verifier) UserID(rawToken string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(rawToken, claims, v.key,
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}), jwt.WithExpirationRequired())
	if errors.Is(err, jwt.ErrTokenExpired) {
		return "", apperror.Wrap(apperror.CodeSessionExpired, "session expired", err)
	}
	if err != nil {
		return "", apperror.Wrap(apperror.CodeUnauthorized, "invalid session", err)
	}
	if _, parseErr := uuid.Parse(claims.Subject); parseErr != nil {
		return "", apperror.Wrap(apperror.CodeUnauthorized, "invalid session subject", parseErr)
	}
	return claims.Subject, nil
}

func (v *Verifier) key(_ *jwt.Token) (any, error) {
	return v.publicKey, nil
}

func Require(verifier *Verifier, logger *slog.Logger) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := verifier.UserID(strings.TrimPrefix(r.Header.Get("Authorization"), bearerPrefix))
			if err != nil {
				httpx.WriteError(w, r, logger, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey{}, userID)))
		})
	}
}

func UserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey{}).(string)
	return userID
}
