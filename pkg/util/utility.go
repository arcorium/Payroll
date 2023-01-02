package util

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// ---------------- START CONSTANT
const TIMEOUT = 10

// ---------------- END CONSTANT

func CreateTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*TIMEOUT)
}

func Hash(str_ string) (string, error) {
	passBytes := []byte(str_)
	hashed, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
	return string(hashed), err
}

func GetValue[T any](val_ T, err error) T {
	return val_
}
