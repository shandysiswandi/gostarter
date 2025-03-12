package inbound

import (
	"encoding/json"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
)

var errInvalidBody = goerror.NewInvalidFormat("Request payload malformed")

type httpEndpoint struct {
	telemetry *telemetry.Telemetry

	loginUC          domain.Login
	registerUC       domain.Register
	verifyUC         domain.Verify
	refreshTokenUC   domain.RefreshToken
	forgotPasswordUC domain.ForgotPassword
	resetPasswordUC  domain.ResetPassword
}

func (h *httpEndpoint) Login(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.http.Login")
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
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.http.Register")
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

func (h *httpEndpoint) Verify(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.http.Verify")
	defer span.End()

	var req VerifyRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.verifyUC.Call(ctx, domain.VerifyInput{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		return nil, err
	}

	return VerifyResponse{
		Email:    resp.Email,
		VerifyAt: resp.VerifyAt.Format(time.RFC3339),
	}, nil
}

func (h *httpEndpoint) RefreshToken(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.http.RefreshToken")
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
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.http.RefreshToken")
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
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "auth.inbound.http.RefreshToken")
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
