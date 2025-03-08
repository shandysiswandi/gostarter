package inbound

import (
	"net/http"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
)

type Inbound struct {
	Router    *framework.Router
	Telemetry *telemetry.Telemetry
	//
	ProfileUC        domain.Profile
	UpdateUC         domain.Update
	UpdatePasswordUC domain.UpdatePassword
	LogoutUC         domain.Logout
}

func (in Inbound) RegisterUserServiceServer() {
	he := &httpEndpoint{
		tel: in.Telemetry,
		//
		profileUC:        in.ProfileUC,
		updateUC:         in.UpdateUC,
		updatePasswordUC: in.UpdatePasswordUC,
		logoutUC:         in.LogoutUC,
	}

	in.Router.Endpoint(http.MethodGet, "/me/profile", he.Profile)
	in.Router.Endpoint(http.MethodPatch, "/me/profile", he.Update)
	in.Router.Endpoint(http.MethodPatch, "/me/password", he.UpdatePassword)
	in.Router.Endpoint(http.MethodPost, "/me/logout", he.Logout)
}
