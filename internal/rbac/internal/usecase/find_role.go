//nolint:dupl // this is not duplicate
package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type FindRoleStore interface {
	FindRole(ctx context.Context, id uint64) (*domain.Role, error)
}

type FindRole struct {
	tele  *telemetry.Telemetry
	store FindRoleStore
}

func NewFindRole(dep Dependency, s FindRoleStore) *FindRole {
	return &FindRole{
		tele:  dep.Telemetry,
		store: s,
	}
}

func (fr *FindRole) Call(ctx context.Context, in domain.FindRoleInput) (*domain.Role, error) {
	ctx, span := fr.tele.Tracer().Start(ctx, "rbac.usecase.FindRole")
	defer span.End()

	role, err := fr.store.FindRole(ctx, in.ID)
	if err != nil {
		fr.tele.Logger().Error(ctx, "failed to find role by id", err)

		return nil, goerror.NewServerInternal(err)
	}

	if role == nil {
		fr.tele.Logger().Warn(ctx, "role is not found")

		return nil, goerror.NewBusiness("role not found", goerror.CodeNotFound)
	}

	return role, nil
}
