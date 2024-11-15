package auth

import (
	"database/sql"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	pb "github.com/shandysiswandi/gostarter/api/gen-proto/auth"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
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
	QueryBuilder goqu.DialectWrapper
	Telemetry    *telemetry.Telemetry
	Router       *httpserver.Router
	GRPCServer   *grpc.Server
	Validator    validation.Validator
	UIDNumber    uid.NumberID
	Hash         hash.Hash
	SecHash      hash.Hash
	JWT          jwt.JWT
}

//nolint:funlen // it's long line because it format param dependency
func New(dep Dependency) (*Expose, error) {
	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes outbound dependencies for core services.
	// This includes setups for outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlAuth := outbound.NewSQLAuth(dep.Database, dep.QueryBuilder, dep.Telemetry)

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes core business logic or use cases to handle user interaction
	loginUC := service.NewLogin(
		dep.Telemetry,
		dep.Validator,
		dep.UIDNumber,
		dep.Hash,
		dep.SecHash,
		dep.JWT,
		sqlAuth,
	)

	registerUC := service.NewRegister(
		dep.Telemetry,
		dep.Validator,
		dep.UIDNumber,
		dep.Hash,
		sqlAuth,
	)

	refreshTokenUC := service.NewRefreshToken(
		dep.Telemetry,
		dep.Validator,
		dep.UIDNumber,
		dep.SecHash,
		dep.JWT,
		sqlAuth,
	)

	forgotPasswordUC := service.NewForgotPassword(
		dep.Telemetry,
		dep.Validator,
		dep.UIDNumber,
		dep.SecHash,
		sqlAuth,
	)

	resetPasswordUC := service.NewResetPassword(
		dep.Telemetry,
		dep.Validator,
		dep.Hash,
		sqlAuth,
	)

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes REST API endpoints to handle core user workflows:
	hEndpoint := inbound.NewHTTPEndpoint(
		loginUC,
		registerUC,
		refreshTokenUC,
		forgotPasswordUC,
		resetPasswordUC,
	)
	dep.Router.Endpoint(http.MethodPost, "/auth/login", hEndpoint.Login)
	dep.Router.Endpoint(http.MethodPost, "/auth/register", hEndpoint.Register)
	dep.Router.Endpoint(http.MethodPost, "/auth/refresh-token", hEndpoint.RefreshToken)
	dep.Router.Endpoint(http.MethodPost, "/auth/forgot-password", hEndpoint.ForgotPassword)
	dep.Router.Endpoint(http.MethodPost, "/auth/reset-password", hEndpoint.ResetPassword)

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes gRPC API endpoints to handle core user workflows:
	gEndpoint := inbound.NewGrpcEndpoint(
		dep.Telemetry,
		loginUC,
		registerUC,
		refreshTokenUC,
		forgotPasswordUC,
		resetPasswordUC,
	)
	pb.RegisterAuthServiceServer(dep.GRPCServer, gEndpoint)

	return &Expose{}, nil
}
