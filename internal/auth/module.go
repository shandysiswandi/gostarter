package auth

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
)

type Expose struct{}

type Dependency struct {
	Database     *sql.DB
	Transaction  dbops.Tx
	QueryBuilder goqu.DialectWrapper
	Telemetry    *telemetry.Telemetry
	Router       *framework.Router
	GRPCServer   *grpc.Server
	Validator    validation.Validator
	UIDNumber    uid.NumberID
	Hash         hash.Hash
	SecHash      hash.Hash
	JWT          jwt.JWT
	Clock        clock.Clocker
}

func New(dep Dependency) (*Expose, error) {
	// This block initializes outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlAuth := outbound.NewSQLAuth(dep.Database, dep.QueryBuilder, dep.Telemetry)

	// This block initializes core business logic or use cases to handle user interaction
	ucDep := usecase.Dependency{
		Telemetry:   dep.Telemetry,
		Validator:   dep.Validator,
		UIDNumber:   dep.UIDNumber,
		Hash:        dep.Hash,
		SecHash:     dep.SecHash,
		JWT:         dep.JWT,
		Clock:       dep.Clock,
		Transaction: dep.Transaction,
	}

	loginUC := usecase.NewLogin(ucDep, sqlAuth)
	registerUC := usecase.NewRegister(ucDep, sqlAuth)
	refreshTokenUC := usecase.NewRefreshToken(ucDep, sqlAuth)
	forgotPasswordUC := usecase.NewForgotPassword(ucDep, sqlAuth)
	resetPasswordUC := usecase.NewResetPassword(ucDep, sqlAuth)

	// This block initializes REST & gRPC API endpoints to handle core user workflows:
	inbound := inbound.Inbound{
		Router:     dep.Router,
		GRPCServer: dep.GRPCServer,
		//
		LoginUC:          loginUC,
		RegisterUC:       registerUC,
		RefreshTokenUC:   refreshTokenUC,
		ForgotPasswordUC: forgotPasswordUC,
		ResetPasswordUC:  resetPasswordUC,
	}
	inbound.RegisterAuthServiceServer()

	return &Expose{}, nil
}
