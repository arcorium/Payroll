package util

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// ---------------- START CONSTANT
const (
	CONTEXT_TIMEOUT            = time.Second * 10
	JWT_REFRESH_TIMEOUT        = time.Hour * 24
	JWT_ACCESS_TIMEOUT         = time.Minute * 15
	JWT_COOKIE_REFRESH_TIMEOUT = JWT_REFRESH_TIMEOUT
	JWT_COOKIE_ACCESS_TIMEOUT  = JWT_ACCESS_TIMEOUT
)

const (
	JWT_SIGNING_METHOD      = "HS256"
	JWT_COOKIE_REFRESH_NAME = "rtoken"
	JWT_COOKIE_ACCESS_NAME  = "atoken"
)

// ---------------- END CONSTANT

func CreateTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
}

func Hash(str_ string) (string, error) {
	passBytes := []byte(str_)
	hashed, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
	return string(hashed), err
}

func GetValue[T any](val_ T, _ error) T {
	return val_
}

func GetError[T any](_ T, err error) error {
	return err
}

func IsEmpty(data_ string) bool {
	return len(data_) < 1
}

func GenerateRefreshToken(claims_ jwt.Claims, secretKey_ []byte) (string, error) {
	// Add extra
	claims := claims_.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(JWT_REFRESH_TIMEOUT).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := token.SignedString(secretKey_)
	return res, err
}

func GenerateAccessToken(claims_ jwt.Claims, secretKey_ []byte) (string, error) {
	// Add extra
	claims := claims_.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(JWT_ACCESS_TIMEOUT).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := token.SignedString(secretKey_)
	return res, err
}

func GenerateBsonObject(data_ any) (bson.M, error) {
	update := bson.M{}
	marshalled, err := bson.Marshal(data_)
	if err != nil {
		return update, err
	}

	err = bson.Unmarshal(marshalled, &update)
	if err != nil {
		return update, err
	}

	return update, nil
}
