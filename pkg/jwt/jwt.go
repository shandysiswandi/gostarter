package jwt

import (
	"errors"
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
