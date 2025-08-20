package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateHashedPassword(password string) (hashedPassword string, err error) {
	hashedPassord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to create hashed password %e", err)
	}
	return string(hashedPassord), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
