package payment

import (
	"database/sql"
	"hash"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Expose struct{}

type Dependency struct {
	Database     *sql.DB
	Transaction  dbops.Tx
	QueryBuilder goqu.DialectWrapper
	Telemetry    *telemetry.Telemetry
	Router       *framework.Router
	Validator    validation.Validator
	UIDNumber    uid.NumberID
	Hash         hash.Hash
	SecHash      hash.Hash
	Clock        clock.Clocker
}

func New(dep Dependency) (*Expose, error) {
	// This block initializes outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlPayment := outbound.NewSQLPayment(dep.Database, dep.QueryBuilder, dep.Telemetry)

	// This block initializes core business logic or use cases to handle user interaction
	ucDep := usecase.Dependency{
		UIDNumber:   dep.UIDNumber,
		Validator:   dep.Validator,
		Transaction: dep.Transaction,
		Telemetry:   dep.Telemetry,
	}
	paymentTopupUC := usecase.NewPaymentTopup(ucDep, sqlPayment)

	// This block initializes REST, SSE, gRPC, and graphQL API endpoints to handle core user workflows:
	inbound := inbound.Inbound{
		Router:         dep.Router,
		PaymentTopupUC: paymentTopupUC,
	}
	inbound.RegisterPaymentServiceServer()

	return &Expose{}, nil
}
