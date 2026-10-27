package keys

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func PrivateFromSeed(encodedSeed string) (ed25519.PrivateKey, error) {
	seed, err := base64.StdEncoding.DecodeString(encodedSeed)
	if err != nil {
		return nil, fmt.Errorf("decode ed25519 seed: %w", err)
	}
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("ed25519 seed must be %d bytes, got %d", ed25519.SeedSize, len(seed))
	}
	return ed25519.NewKeyFromSeed(seed), nil
}

func Public(encodedKey string) (ed25519.PublicKey, error) {
	raw, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("decode ed25519 public key: %w", err)
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("ed25519 public key must be %d bytes, got %d", ed25519.PublicKeySize, len(raw))
	}
	return ed25519.PublicKey(raw), nil
}
