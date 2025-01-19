//nolint:dupl // this is not duplicate
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

type FetchPermissionStore interface {
	FetchPermission(ctx context.Context, filter map[string]any) ([]domain.Permission, error)
}

type FetchPermission struct {
	tele  *telemetry.Telemetry
	store FetchPermissionStore
}

func NewFetchPermission(dep Dependency, s FetchPermissionStore) *FetchPermission {
	return &FetchPermission{
		tele:  dep.Telemetry,
		store: s,
	}
}

func (fr *FetchPermission) Call(ctx context.Context, in domain.FetchPermissionInput) (
	*domain.FetchPermissionOutput, error,
) {
	ctx, span := fr.tele.Tracer().Start(ctx, "rbac.usecase.FetchPermission")
	defer span.End()

	cursor, limit := pagination.ParseCursorBased(in.Cursor, in.Limit)

	filter := map[string]any{
		"cursor": cursor,
		"limit":  limit,
	}

	if in.Name != "" {
		filter["name"] = in.Name
	}

	pems, err := fr.store.FetchPermission(ctx, filter)
	if err != nil {
		fr.tele.Logger().Error(ctx, "failed to fetch permissions filter by name", err)

		return nil, goerror.NewServerInternal(err)
	}

	nc := ""
	hm := len(pems) > limit

	if hm {
		nc = base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatUint(pems[limit].ID, 10)))
		pems = pems[:limit]
	}

	return &domain.FetchPermissionOutput{
		Permissions: pems,
		NextCursor:  nc,
		HasMore:     hm,
	}, nil
}
