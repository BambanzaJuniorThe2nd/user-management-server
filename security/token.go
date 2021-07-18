package security

import (
	"errors"
	"fmt"
	"server/models"
	"time"
	"server/config"

	jwt "github.com/form3tech-oss/jwt-go"
)

var (
	JwtSecretKey     = []byte(config.GetConfig().Token)
	JwtSigningMethod = jwt.SigningMethodHS256.Name
)

var ErrInvalidAuthToken   = errors.New("invalid auth-token")

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
