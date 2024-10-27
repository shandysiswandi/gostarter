package auth

import (
	"database/sql"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Expose struct{}

type Dependency struct {
	Config    config.Config
	Database  *sql.DB
	Telemetry *telemetry.Telemetry
	Router    *httprouter.Router
	Validator validation.Validator
	UIDNumber uid.NumberID
	Hash      hash.Hash
	SecHash   hash.Hash
	JWT       jwt.JWT
}

func New(dep Dependency) (*Expose, error) {
	// init outbound | database | http client | grpc client | redis | etc.
	sqlAuth := outbound.NewSQLAuth(dep.Database, dep.Config)

	// init services | useCases | business logic
	loginUC := service.NewLogin(dep.Telemetry, dep.Validator, dep.UIDNumber,
		dep.Hash, dep.SecHash, dep.JWT, sqlAuth)

	registerUC := service.NewRegister(dep.Telemetry, dep.Validator, dep.UIDNumber, dep.Hash, sqlAuth)

	refreshTokenUC := service.NewRefreshToken(dep.Telemetry, dep.Validator, dep.UIDNumber,
		dep.SecHash, dep.JWT, sqlAuth)

	forgotPasswordUC := service.NewForgotPassword(dep.Telemetry, dep.Validator, dep.UIDNumber, dep.SecHash, sqlAuth)

	resetPasswordUC := service.NewResetPassword(dep.Telemetry, dep.Validator, dep.Hash, sqlAuth)

	// register endpoint REST
	inbound.RegisterRESTEndpoint(dep.Router, dep.Telemetry.Logger(), inbound.NewEndpoint(
		loginUC,
		registerUC,
		refreshTokenUC,
		forgotPasswordUC,
		resetPasswordUC,
	))

	return &Expose{}, nil
}
