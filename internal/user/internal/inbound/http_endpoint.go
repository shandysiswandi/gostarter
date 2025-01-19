package inbound

import (
	"encoding/json"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type httpEndpoint struct {
	tel *telemetry.Telemetry

	profileUC        domain.Profile
	updateUC         domain.Update
	updatePasswordUC domain.UpdatePassword
	logoutUC         domain.Logout
}

func (h *httpEndpoint) Profile(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "user.inbound.httpEndpoint.Profile")
	defer span.End()

	resp, err := h.profileUC.Call(ctx, domain.ProfileInput{})
	if err != nil {
		return nil, err
	}

	return User{
		ID:    resp.ID,
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}

func (h *httpEndpoint) Update(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "user.inbound.httpEndpoint.Update")
	defer span.End()

	var req UpdateRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat("invalid request body")
	}

	resp, err := h.updateUC.Call(ctx, domain.UpdateInput{Name: req.Name})
	if err != nil {
		return nil, err
	}

	return User{
		ID:    resp.ID,
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}

func (h *httpEndpoint) UpdatePassword(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "user.inbound.httpEndpoint.UpdatePassword")
	defer span.End()

	var req UpdatePasswordRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat("invalid request body")
	}

	resp, err := h.updatePasswordUC.Call(ctx, domain.UpdatePasswordInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return User{
		ID:    resp.ID,
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}

func (h *httpEndpoint) Logout(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "user.inbound.httpEndpoint.Logout")
	defer span.End()

	resp, err := h.logoutUC.Call(ctx, domain.LogoutInput{AccessToken: c.Header().Get("Authorization")})
	if err != nil {
		return nil, err
	}

	return LogoutResponse{Message: resp.Message}, nil
}
