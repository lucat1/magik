package jwt

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWToken struct {
	secret string
}

type TokenClaims struct {
	Email string
	jwt.StandardClaims
}

func (t JWToken) Generate(email string, length time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(length).Unix(),
			Issuer:    "magik-jwt",
		}})
	return token.SignedString(token)
}

func (t JWToken) Validate(raw string) (string, error) {
	claims := TokenClaims{}
	token, err := jwt.ParseWithClaims(raw, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.Email, nil
	} else {
		return "", errors.New("wrong jwt claims type or token invalide")
	}

}

func NewGenerator(secret string) JWToken {
	return JWToken{secret: secret}
}
