package inbound

import (
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
)

type httpEndpoint struct {
	profileUC domain.Profile
}

func (e *httpEndpoint) Profile(c framework.Context) (any, error) {
	var email string
	clm := jwt.GetClaim(c.Context())
	if clm != nil {
		email = clm.Email
	}

	resp, err := e.profileUC.Call(c.Context(), domain.ProfileInput{Email: email})
	if err != nil {
		return nil, err
	}

	return ProfileResponse{
		ID:    resp.ID,
		Email: resp.Email,
	}, nil
}
