package outbound

import (
	"context"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type SQL struct {
	db        *sqlkit.DB
	telemetry *telemetry.Telemetry
}

func NewSQL(db *sqlkit.DB, tel *telemetry.Telemetry) *SQL {
	return &SQL{
		db:        db,
		telemetry: tel,
	}
}

func (s *SQL) User(ctx context.Context, id uint64) (*domain.User, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "user.outbound.SQL.User")
	defer span.End()

	return sqlkit.One[domain.User](ctx, s.db, sqlkit.Ex{"id": id})
}

func (s *SQL) UserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "user.outbound.SQL.UserByEmail")
	defer span.End()

	return sqlkit.One[domain.User](ctx, s.db, sqlkit.Ex{"email": email})
}

func (s *SQL) UserUpdate(ctx context.Context, user map[string]any) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "user.outbound.SQL.UserUpdate")
	defer span.End()

	id := user["id"]
	delete(user, "id")

	_, err := sqlkit.Update[domain.User](ctx, s.db, user, sqlkit.Ex{"id": id})

	return err
}

func (s *SQL) TokenDeleteByAccess(ctx context.Context, token string) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "user.outbound.SQL.TokenDeleteByAccess")
	defer span.End()

	_, err := sqlkit.Delete[domain.Token](ctx, s.db, sqlkit.Ex{"access_token": token})

	return err
}
