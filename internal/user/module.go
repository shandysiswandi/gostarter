package user

import (
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/user/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/user/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/user/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type Expose struct{}

type Dependency struct {
	SQLKitDB  *sqlkit.DB
	Validator validation.Validator
	Hash      hash.Hash
	Router    *framework.Router
	Telemetry *telemetry.Telemetry
}

func New(dep Dependency) (*Expose, error) {
	// This block initializes outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlUser := outbound.NewSQL(dep.SQLKitDB, dep.Telemetry)

	// This block initializes core business logic or use cases to handle user interaction
	ucDep := usecase.Dependency{
		Telemetry: dep.Telemetry,
		Validator: dep.Validator,
		Hash:      dep.Hash,
	}
	profile := usecase.NewProfile(ucDep, sqlUser)
	update := usecase.NewUpdate(ucDep, sqlUser)
	updatePassword := usecase.NewUpdatePassword(ucDep, sqlUser)
	logout := usecase.NewLogout(ucDep, sqlUser)

	// This block initializes REST, SSE, gRPC, and graphQL API endpoints to handle core user workflows:
	inbound := inbound.Inbound{
		Router:    dep.Router,
		Telemetry: dep.Telemetry,
		//
		ProfileUC:        profile,
		UpdateUC:         update,
		UpdatePasswordUC: updatePassword,
		LogoutUC:         logout,
	}
	inbound.RegisterUserServiceServer()

	return &Expose{}, nil
}
