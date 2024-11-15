package inbound

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

type Endpoint struct {
	loginUC          domain.Login
	registerUC       domain.Register
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func NewHTTPEndpoint(
	loginUC domain.Login,
	registerUC domain.Register,
	refreshTokenUC domain.RefreshToken,
	forgotPasswordUC domain.ForgotPassword,
	resetPasswordUC domain.ResetPassword,
) *Endpoint {
	return &Endpoint{
		loginUC:          loginUC,
		registerUC:       registerUC,
		refreshTokenUC:   refreshTokenUC,
		forgotPasswordUC: forgotPasswordUC,
		resetPasswordUC:  resetPasswordUC,
	}
}

func (e *Endpoint) Login(ctx context.Context, r *http.Request) (any, error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
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

func (e *Endpoint) Register(ctx context.Context, r *http.Request) (any, error) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	_, err := e.registerUC.Call(ctx, domain.RegisterInput{Email: req.Email, Password: req.Password})
	if err != nil {
		return nil, err
	}

	return RegisterResponse{Email: req.Email}, nil
}

func (e *Endpoint) RefreshToken(ctx context.Context, r *http.Request) (any, error) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
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

func (e *Endpoint) ForgotPassword(ctx context.Context, r *http.Request) (any, error) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	_, err := e.forgotPasswordUC.Call(ctx, domain.ForgotPasswordInput{Email: req.Email})
	if err != nil {
		return nil, err
	}

	return ForgotPasswordResponse{
		Email:   req.Email,
		Message: "If an account with this email exists, you'll receive a password reset email shortly.",
	}, nil
}

func (e *Endpoint) ResetPassword(ctx context.Context, r *http.Request) (any, error) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	_, err := e.resetPasswordUC.Call(ctx, domain.ResetPasswordInput{
		Token: req.Token, Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return ResetPasswordResponse{Message: "Your password has been successfully reset."}, nil
}
