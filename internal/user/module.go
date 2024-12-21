package user

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/user/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/user/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/user/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Expose struct{}

type Dependency struct {
	Database     *sql.DB
	QueryBuilder goqu.DialectWrapper
	Validator    validation.Validator
	Router       *framework.Router
	Telemetry    *telemetry.Telemetry
}

func New(dep Dependency) (*Expose, error) {
	// This block initializes outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlUser := outbound.NewSQLUser(dep.Database, dep.QueryBuilder, dep.Telemetry)

	// This block initializes core business logic or use cases to handle user interaction
	ucDep := usecase.Dependency{
		Telemetry: dep.Telemetry,
		Validator: dep.Validator,
	}
	profile := usecase.NewProfile(ucDep, sqlUser)
	update := usecase.NewUpdate(ucDep, sqlUser)
	logout := usecase.NewLogout(ucDep, sqlUser)

	// This block initializes REST, SSE, gRPC, and graphQL API endpoints to handle core user workflows:
	inbound := inbound.Inbound{
		Router:    dep.Router,
		Telemetry: dep.Telemetry,
		//
		ProfileUC: profile,
		UpdateUC:  update,
		LogoutUC:  logout,
	}
	inbound.RegisterUserServiceServer()

	return &Expose{}, nil
}
