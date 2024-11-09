package auth

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Expose struct{}

type Dependency struct {
	Database     *sql.DB
	QueryBuilder goqu.DialectWrapper
	Telemetry    *telemetry.Telemetry
	Router       *httprouter.Router
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
	endpoint := inbound.NewEndpoint(
		loginUC,
		registerUC,
		refreshTokenUC,
		forgotPasswordUC,
		resetPasswordUC,
	)
	inbound.RegisterRESTEndpoint(
		dep.Router,
		dep.Telemetry.Logger(),
		endpoint,
	)

	return &Expose{}, nil
}
