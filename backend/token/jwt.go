package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKey = 32

type JWTMaker struct {
	secretKey string
}

func NewJwtMaker(sectetKey string) (Maker, error) {
	if len(sectetKey) < minSecretKey {
		return nil, fmt.Errorf("invalid key size, must be atleast %d charecters", minSecretKey)
	}
	maker := JWTMaker{
		secretKey: sectetKey,
	}
	return &maker, nil
}

type CustomClaims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func (maker *JWTMaker) CreateToken(username string, userID int64, isAdmin bool, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, userID, isAdmin, duration)
	if err != nil {
		return "", err
	}
	// jwtClaim := jwt.RegisteredClaims{
	// 	IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
	// 	ExpiresAt: jwt.NewNumericDate(payload.ExpiresAt),
	// 	Audience:  []string{payload.Username},
	// 	ID:        payload.Id.String(),
	// }
	jwtClaim := CustomClaims{
		Username: username,
		UserID:   userID,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(payload.ExpiresAt),
			ID:        payload.Id.String(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	return jwtToken.SignedString([]byte(maker.secretKey))

}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		Id:        id,
		Username:  claims.Username,
		UserID:    claims.UserID,
		IsAdmin:   claims.IsAdmin,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	return payload, nil
}
