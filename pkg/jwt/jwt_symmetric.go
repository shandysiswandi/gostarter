package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type JWTSymetric struct {
	secret []byte
}

func NewJWTSymetric(secret []byte) *JWTSymetric {
	return &JWTSymetric{secret: secret}
}

func (js *JWTSymetric) Generate(c *Claim) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(js.secret)
}

func (js *JWTSymetric) Verify(token string) (*Claim, error) {
	tkn, err := jwt.ParseWithClaims(token, &Claim{}, func(t *jwt.Token) (interface{}, error) {
		return js.secret, nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ErrTokenExpired
	}

	if err != nil {
		return nil, err
	}

	if claims, ok := tkn.Claims.(*Claim); ok {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
