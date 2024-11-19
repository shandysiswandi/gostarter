package inbound

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
)

const msgInvalidBody = "invalid request body"

type httpEndpoint struct {
	loginUC          domain.Login
	registerUC       domain.Register
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func (e *httpEndpoint) Login(ctx context.Context, r *http.Request) (any, error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat(msgInvalidBody)
	}

	resp, err := e.loginUC.Call(ctx, domain.LoginInput{Email: req.Email, Password: req.Password})
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

func (e *httpEndpoint) Register(ctx context.Context, r *http.Request) (any, error) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat(msgInvalidBody)
	}

	resp, err := e.registerUC.Call(ctx, domain.RegisterInput{Email: req.Email, Password: req.Password})
	if err != nil {
		return nil, err
	}

	return RegisterResponse{Email: resp.Email}, nil
}

func (e *httpEndpoint) RefreshToken(ctx context.Context, r *http.Request) (any, error) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat(msgInvalidBody)
	}

	resp, err := e.refreshTokenUC.Call(ctx, domain.RefreshTokenInput{RefreshToken: req.RefreshToken})
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

func (e *httpEndpoint) ForgotPassword(ctx context.Context, r *http.Request) (any, error) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat(msgInvalidBody)
	}

	resp, err := e.forgotPasswordUC.Call(ctx, domain.ForgotPasswordInput{Email: req.Email})
	if err != nil {
		return nil, err
	}

	return ForgotPasswordResponse{Email: resp.Email, Message: resp.Message}, nil
}

func (e *httpEndpoint) ResetPassword(ctx context.Context, r *http.Request) (any, error) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat(msgInvalidBody)
	}

	resp, err := e.resetPasswordUC.Call(ctx, domain.ResetPasswordInput{
		Token: req.Token, Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return ResetPasswordResponse{Message: resp.Message}, nil
}
