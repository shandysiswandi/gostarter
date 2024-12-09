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
	LogoutUC  domain.Logout
}

func (in Inbound) RegisterUserServiceServer() {
	he := &httpEndpoint{
		profileUC: in.ProfileUC,
		logoutUC:  in.LogoutUC,
	}

	in.Router.Endpoint(http.MethodGet, "/users/profile", he.Profile)
	in.Router.Endpoint(http.MethodPost, "/users/logout", he.Logout)
}
