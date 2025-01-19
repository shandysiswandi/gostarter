package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type Asymmetric struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTAsymmetric(private, public string) (*Asymmetric, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(private)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, ErrParsePrivateKey
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidRSAPrivateKey
	}

	// ---
	publicKeyBytes, err := base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, err
	}

	block, _ = pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, ErrParsePublicKey
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidRSAPublicKey
	}

	return &Asymmetric{
		privateKey: rsaPrivateKey,
		publicKey:  rsaPublicKey,
	}, nil
}

func (ja *Asymmetric) Generate(c *Claim) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodRS256, c).SignedString(ja.privateKey)
}

func (ja *Asymmetric) Verify(token string) (*Claim, error) {
	tkn, err := jwt.ParseWithClaims(token, &Claim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return ja.publicKey, nil
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
