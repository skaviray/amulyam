package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20/chacha"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type Paseto struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d charecters", chacha20poly1305.KeySize)
	}
	maker := &Paseto{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *Paseto) CreateToken(username string, userID int64, isAdmin bool, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, userID, isAdmin, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *Paseto) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
