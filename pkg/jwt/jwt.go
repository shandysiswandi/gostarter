package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrParsePrivateKey      = errors.New("failed to parse private key")
	ErrParsePublicKey       = errors.New("failed to parse public key")
	ErrInvalidRSAPrivateKey = errors.New("invalid rsa private key")
	ErrInvalidRSAPublicKey  = errors.New("invalid rsa public key")
	ErrTokenExpired         = errors.New("token is expired")
)

type JWT interface {
	Generate(c *Claim) (string, error)
	Verify(token string) (*Claim, error)
}

type JSONWebToken struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJSONWebToken(private, public string) (*JSONWebToken, error) {
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

	return &JSONWebToken{
		privateKey: rsaPrivateKey,
		publicKey:  rsaPublicKey,
	}, nil
}

func (m *JSONWebToken) Generate(c *Claim) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodRS256, c).SignedString(m.privateKey)
}

func (m *JSONWebToken) Verify(token string) (*Claim, error) {
	tkn, err := jwt.ParseWithClaims(token, &Claim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return m.publicKey, nil
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
