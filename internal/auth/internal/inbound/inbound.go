package inbound

import (
	"net/http"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
)

type Inbound struct {
	Router    *framework.Router
	Telemetry *telemetry.Telemetry
	//
	LoginUC          domain.Login
	RegisterUC       domain.Register
	VerifyUC         domain.Verify
	RefreshTokenUC   domain.RefreshToken
	ForgotPasswordUC domain.ForgotPassword
	ResetPasswordUC  domain.ResetPassword
}

func (in Inbound) RegisterAuthServiceServer() {
	he := &httpEndpoint{
		telemetry: in.Telemetry,
		//
		loginUC:          in.LoginUC,
		registerUC:       in.RegisterUC,
		verifyUC:         in.VerifyUC,
		refreshTokenUC:   in.RefreshTokenUC,
		forgotPasswordUC: in.ForgotPasswordUC,
		resetPasswordUC:  in.ResetPasswordUC,
	}

	in.Router.Endpoint(http.MethodPost, "/auth/login", he.Login)
	in.Router.Endpoint(http.MethodPost, "/auth/register", he.Register)
	in.Router.Endpoint(http.MethodPost, "/auth/verify", he.Verify)
	in.Router.Endpoint(http.MethodPost, "/auth/refresh-token", he.RefreshToken)
	in.Router.Endpoint(http.MethodPost, "/auth/forgot-password", he.ForgotPassword)
	in.Router.Endpoint(http.MethodPost, "/auth/reset-password", he.ResetPassword)
}
