package auth

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoToken struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// Create implements Token.
func (token *PasetoToken) Create(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return token.paseto.Encrypt(token.symmetricKey, payload, nil)
}

// Verify implements Token.
func (token *PasetoToken) Verify(tokenString string) (*Payload, error) {
	payload := &Payload{}

	err := token.paseto.Decrypt(tokenString, token.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func NewPasetoToken(symmetricKey string) (Token, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	token := &PasetoToken{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return token, nil
}
