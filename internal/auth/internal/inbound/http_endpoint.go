package inbound

import (
	"encoding/json"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

var errInvalidBody = goerror.NewInvalidFormat("invalid request body")

type httpEndpoint struct {
	telemetry *telemetry.Telemetry

	loginUC          domain.Login
	registerUC       domain.Register
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func (h *httpEndpoint) Login(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.httpEndpoint.Login")
	defer span.End()

	var req LoginRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.loginUC.Call(ctx, domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
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

func (h *httpEndpoint) Register(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.httpEndpoint.Register")
	defer span.End()

	var req RegisterRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.registerUC.Call(ctx, domain.RegisterInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return RegisterResponse{Email: resp.Email}, nil
}

func (h *httpEndpoint) RefreshToken(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.httpEndpoint.RefreshToken")
	defer span.End()

	var req RefreshTokenRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.refreshTokenUC.Call(ctx, domain.RefreshTokenInput{RefreshToken: req.RefreshToken})
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

func (h *httpEndpoint) ForgotPassword(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.httpEndpoint.RefreshToken")
	defer span.End()

	var req ForgotPasswordRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.forgotPasswordUC.Call(ctx, domain.ForgotPasswordInput{Email: req.Email})
	if err != nil {
		return nil, err
	}

	return ForgotPasswordResponse{
		Email:   resp.Email,
		Message: resp.Message,
	}, nil
}

func (h *httpEndpoint) ResetPassword(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.httpEndpoint.RefreshToken")
	defer span.End()

	var req ResetPasswordRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.resetPasswordUC.Call(ctx, domain.ResetPasswordInput{
		Token:    req.Token,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return ResetPasswordResponse{Message: resp.Message}, nil
}
