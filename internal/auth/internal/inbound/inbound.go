package inbound

import (
	"net/http"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/auth"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"google.golang.org/grpc"
)

type Inbound struct {
	Router     *framework.Router
	GRPCServer *grpc.Server
	//
	LoginUC          domain.Login
	RegisterUC       domain.Register
	RefreshTokenUC   domain.RefreshToken
	ForgotPasswordUC domain.ForgotPassword
	ResetPasswordUC  domain.ResetPassword
}

func (in Inbound) RegisterAuthServiceServer() {
	he := &httpEndpoint{
		loginUC:          in.LoginUC,
		registerUC:       in.RegisterUC,
		refreshTokenUC:   in.RefreshTokenUC,
		forgotPasswordUC: in.ForgotPasswordUC,
		resetPasswordUC:  in.ResetPasswordUC,
	}

	ge := &GrpcEndpoint{
		loginUC:          in.LoginUC,
		registerUC:       in.RegisterUC,
		refreshTokenUC:   in.RefreshTokenUC,
		forgotPasswordUC: in.ForgotPasswordUC,
		resetPasswordUC:  in.ResetPasswordUC,
	}

	in.Router.Endpoint(http.MethodPost, "/auth/login", he.Login)
	in.Router.Endpoint(http.MethodPost, "/auth/register", he.Register)
	in.Router.Endpoint(http.MethodPost, "/auth/refresh-token", he.RefreshToken)
	in.Router.Endpoint(http.MethodPost, "/auth/forgot-password", he.ForgotPassword)
	in.Router.Endpoint(http.MethodPost, "/auth/reset-password", he.ResetPassword)

	pb.RegisterAuthServiceServer(in.GRPCServer, ge)
}
