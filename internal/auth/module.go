package auth

import (
	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/jwt"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
	"google.golang.org/grpc"
)

type Expose struct{}

type Dependency struct {
	SQLKitDB   *sqlkit.DB
	Telemetry  *telemetry.Telemetry
	Router     *framework.Router
	GRPCServer *grpc.Server
	Validator  validation.Validator
	UIDNumber  uid.NumberID
	Hash       hash.Hash
	SecHash    hash.Hash
	JWT        jwt.JWT
	Clock      clock.Clocker
}

func New(dep Dependency) (*Expose, error) {
	// This block initializes outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlAuth := outbound.NewSQL(dep.SQLKitDB, dep.Telemetry)

	// This block initializes core business logic or use cases to handle user interaction
	ucDep := usecase.Dependency{
		Telemetry:   dep.Telemetry,
		Validator:   dep.Validator,
		UIDNumber:   dep.UIDNumber,
		Hash:        dep.Hash,
		SecHash:     dep.SecHash,
		JWT:         dep.JWT,
		Clock:       dep.Clock,
		Transaction: dep.SQLKitDB.Tx(),
	}

	loginUC := usecase.NewLogin(ucDep, sqlAuth)
	registerUC := usecase.NewRegister(ucDep, sqlAuth)
	verifyUC := usecase.NewVerify(ucDep, sqlAuth)
	refreshTokenUC := usecase.NewRefreshToken(ucDep, sqlAuth)
	forgotPasswordUC := usecase.NewForgotPassword(ucDep, sqlAuth)
	resetPasswordUC := usecase.NewResetPassword(ucDep, sqlAuth)

	// This block initializes REST & gRPC API endpoints to handle core user workflows:
	inbound := inbound.Inbound{
		Router:    dep.Router,
		Telemetry: dep.Telemetry,
		//
		LoginUC:          loginUC,
		RegisterUC:       registerUC,
		VerifyUC:         verifyUC,
		RefreshTokenUC:   refreshTokenUC,
		ForgotPasswordUC: forgotPasswordUC,
		ResetPasswordUC:  resetPasswordUC,
	}
	inbound.RegisterAuthServiceServer()

	return &Expose{}, nil
}
