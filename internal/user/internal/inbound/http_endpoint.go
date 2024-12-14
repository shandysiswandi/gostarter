package inbound

import (
	"encoding/json"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
)

type httpEndpoint struct {
	profileUC domain.Profile
	updateUC  domain.Update
	logoutUC  domain.Logout
}

func (e *httpEndpoint) Profile(c framework.Context) (any, error) {
	var email string
	clm := jwt.GetClaim(c.Context())
	if clm != nil {
		email = clm.Subject
	}

	resp, err := e.profileUC.Call(c.Context(), domain.ProfileInput{Email: email})
	if err != nil {
		return nil, err
	}

	return User{
		ID:    resp.ID,
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}

func (e *httpEndpoint) Update(c framework.Context) (any, error) {
	var uid uint64
	clm := jwt.GetClaim(c.Context())
	if clm != nil {
		uid = clm.AuthID
	}

	var req UpdateRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, goerror.NewInvalidFormat("invalid request body")
	}

	resp, err := e.updateUC.Call(c.Context(), domain.UpdateInput{ID: uid, Name: req.Name})
	if err != nil {
		return nil, err
	}

	return User{
		ID:    resp.ID,
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}

func (e *httpEndpoint) Logout(c framework.Context) (any, error) {
	ac := c.Header().Get("Authorization")
	resp, err := e.logoutUC.Call(c.Context(), domain.LogoutInput{AccessToken: ac})
	if err != nil {
		return nil, err
	}

	return LogoutResponse{Message: resp.Message}, nil
}
