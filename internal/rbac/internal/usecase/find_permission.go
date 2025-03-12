//nolint:dupl // this is not duplicate
package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
)

type FindPermissionStore interface {
	FindPermission(ctx context.Context, id uint64) (*domain.Permission, error)
}

type FindPermission struct {
	tele  *telemetry.Telemetry
	store FindPermissionStore
}

func NewFindPermission(dep Dependency, s FindPermissionStore) *FindPermission {
	return &FindPermission{
		tele:  dep.Telemetry,
		store: s,
	}
}

func (fp *FindPermission) Call(ctx context.Context, in domain.FindPermissionInput) (
	*domain.Permission, error,
) {
	ctx, span := fp.tele.Tracer().Start(ctx, "rbac.usecase.FindPermission")
	defer span.End()

	perm, err := fp.store.FindPermission(ctx, in.ID)
	if err != nil {
		fp.tele.Logger().Error(ctx, "failed to find permission by id", err)

		return nil, goerror.NewServerInternal(err)
	}

	if perm == nil {
		fp.tele.Logger().Warn(ctx, "permission is not found")

		return nil, goerror.NewBusiness("permission not found", goerror.CodeNotFound)
	}

	return perm, nil
}
