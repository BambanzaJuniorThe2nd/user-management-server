package security

import (
	"errors"
	"fmt"
	"os"
	"server/models"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
)

var (
	JwtSecretKey     = []byte(os.Getenv("JWT_SECRET_KEY"))
	JwtSigningMethod = jwt.SigningMethodHS256.Name
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmptyPassword      = errors.New("password can't be empty")
	ErrInvalidAuthToken   = errors.New("invalid auth-token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("Unauthorized")
)

type MyCustomClaims struct {
	jwt.StandardClaims
	IsAdmin bool `json:"isAdmin"`
}

func NewToken(user *models.User) (string, error) {
	// Create the Claims
	claims := MyCustomClaims{
		jwt.StandardClaims{
			Id:        user.ID,
			Issuer:    user.ID,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
		user.IsAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func validateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return JwtSecretKey, nil
}

func ParseToken(tokenString string) (*MyCustomClaims, error) {
	claims := new(MyCustomClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, validateSignedMethod)
	if err != nil {
		return nil, err
	}
	var ok bool
	claims, ok = token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidAuthToken
	}
	return claims, nil
}
