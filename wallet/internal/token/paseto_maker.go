package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/namnv2496/go-wallet/config"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(
	config config.Config,
) (Maker, error) {
	if len(config.TokenSymmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(config.TokenSymmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userId int64, username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, username, role, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
