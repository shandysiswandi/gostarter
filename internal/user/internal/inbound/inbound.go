package inbound

import (
	"net/http"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
)

type Inbound struct {
	Router *framework.Router
	//
	ProfileUC domain.Profile
	UpdateUC  domain.Update
	LogoutUC  domain.Logout
}

func (in Inbound) RegisterUserServiceServer() {
	he := &httpEndpoint{
		profileUC: in.ProfileUC,
		updateUC:  in.UpdateUC,
		logoutUC:  in.LogoutUC,
	}

	in.Router.Endpoint(http.MethodGet, "/me/profile", he.Profile)
	in.Router.Endpoint(http.MethodPatch, "/me/update", he.Update)
	in.Router.Endpoint(http.MethodPost, "/me/logout", he.Logout)
}
