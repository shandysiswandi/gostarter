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
}

func (in Inbound) RegisterUserServiceServer() {
	he := &httpEndpoint{
		profileUC: in.ProfileUC,
	}

	in.Router.Endpoint(http.MethodGet, "/users/profile", he.Profile)
}
