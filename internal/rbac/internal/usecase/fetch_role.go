package usecase

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/pagination"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type FetchRoleStore interface {
	FetchRole(ctx context.Context, filter map[string]any) ([]domain.Role, error)
}

type FetchRole struct {
	tele  *telemetry.Telemetry
	store FetchRoleStore
}

func NewFetchRole(dep Dependency, s FetchRoleStore) *FetchRole {
	return &FetchRole{
		tele:  dep.Telemetry,
		store: s,
	}
}

func (fr *FetchRole) Call(ctx context.Context, in domain.FetchRoleInput) (*domain.FetchRoleOutput, error) {
	ctx, span := fr.tele.Tracer().Start(ctx, "rbac.usecase.FetchRole")
	defer span.End()

	cursor, limit := pagination.ParseCursorBased(in.Cursor, in.Limit)

	filter := map[string]any{
		"cursor": cursor,
		"limit":  limit,
	}

	if in.Name != "" {
		filter["name"] = in.Name
	}

	roles, err := fr.store.FetchRole(ctx, filter)
	if err != nil {
		fr.tele.Logger().Error(ctx, "failed to fetch roles filter by name", err)

		return nil, goerror.NewServerInternal(err)
	}

	nc := ""
	hm := len(roles) > limit

	if hm {
		nc = base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatUint(roles[limit].ID, 10)))
		roles = roles[:limit]
	}

	return &domain.FetchRoleOutput{
		Roles:      roles,
		NextCursor: nc,
		HasMore:    hm,
	}, nil
}
