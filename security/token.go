package security

import (
	"fmt"
	"os"
	"server/models"
	"server/util"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
)

var (
	JwtSecretKey     = []byte(os.Getenv("JWT_SECRET_KEY"))
	JwtSigningMethod = jwt.SigningMethodHS256.Name
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
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
		user.IsAdmin,
	}
	// claims := jwt.StandardClaims{
	// 	Id:        user.Id,
	// 	Issuer:    user.Id,
	// 	IssuedAt:  time.Now().Unix(),
	// 	ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
	// }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func validateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return JwtSecretKey, nil
}

func ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, validateSignedMethod)
	if err != nil {
		return nil, err
	}
	var ok bool
	claims, ok = token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, util.ErrInvalidAuthToken
	}
	return claims, nil
}
