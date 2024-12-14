package inbound

import (
	"net/http"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"google.golang.org/grpc"
)

type Inbound struct {
	Router     *framework.Router
	GRPCServer *grpc.Server
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
