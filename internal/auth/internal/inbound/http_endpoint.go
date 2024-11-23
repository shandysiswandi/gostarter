package inbound

import (
	"encoding/json"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
)

var errInvalidBody = goerror.NewInvalidFormat("invalid request body")

type httpEndpoint struct {
	loginUC          domain.Login
	registerUC       domain.Register
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func (e *httpEndpoint) Login(c framework.Context) (any, error) {
	var req LoginRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.loginUC.Call(c.Context(), domain.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		return nil, err
	}

	return LoginResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresIn:  resp.AccessExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
	}, nil
}

func (e *httpEndpoint) Register(c framework.Context) (any, error) {
	var req RegisterRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.registerUC.Call(c.Context(), domain.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return RegisterResponse{Email: resp.Email}, nil
}

func (e *httpEndpoint) RefreshToken(c framework.Context) (any, error) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.refreshTokenUC.Call(c.Context(), domain.RefreshTokenInput{RefreshToken: req.RefreshToken})
	if err != nil {
		return nil, err
	}

	return RefreshTokenResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresIn:  resp.AccessExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
	}, nil
}

func (e *httpEndpoint) ForgotPassword(c framework.Context) (any, error) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.forgotPasswordUC.Call(c.Context(), domain.ForgotPasswordInput{Email: req.Email})
	if err != nil {
		return nil, err
	}

	return ForgotPasswordResponse{Email: resp.Email, Message: resp.Message}, nil
}

func (e *httpEndpoint) ResetPassword(c framework.Context) (any, error) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.resetPasswordUC.Call(c.Context(), domain.ResetPasswordInput{
		Token:    req.Token,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return ResetPasswordResponse{Message: resp.Message}, nil
}
