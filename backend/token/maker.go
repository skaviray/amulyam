package token

import (
	"time"
)

type Maker interface {
	CreateToken(uusername string, userID int64, isAdmin bool, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
