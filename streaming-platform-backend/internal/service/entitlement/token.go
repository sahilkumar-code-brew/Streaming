package entitlement

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	domain "github.com/code-brew-ai/streaming-platform-backend/internal/domain/entitlement"
)

type PlaybackClaims struct {
	StreamID     string `json:"sid"`
	DeviceIDHash string `json:"did"`
	jwt.RegisteredClaims
}

type TokenSigner interface {
	Sign(request domain.Request, issuedAt time.Time) (string, time.Time, error)
}

type PlaybackSigner struct {
	privateKey ed25519.PrivateKey
	keyID      string
	ttl        time.Duration
}

func NewPlaybackSigner(privateKey ed25519.PrivateKey, keyID string, ttl time.Duration) *PlaybackSigner {
	return &PlaybackSigner{privateKey: privateKey, keyID: keyID, ttl: ttl}
}

func (s *PlaybackSigner) Sign(request domain.Request, issuedAt time.Time) (string, time.Time, error) {
	expiresAt := issuedAt.Add(s.ttl)
	claims := PlaybackClaims{
		StreamID:     request.StreamID,
		DeviceIDHash: hashDeviceID(request.DeviceID),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   request.UserID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = s.keyID
	signed, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign playback token: %w", err)
	}
	return signed, expiresAt, nil
}

func hashDeviceID(deviceID string) string {
	digest := sha256.Sum256([]byte(deviceID))
	return hex.EncodeToString(digest[:])
}
