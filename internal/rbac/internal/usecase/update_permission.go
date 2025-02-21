//nolint:dupl // this is not duplicate
package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdatePermissionStore interface {
	FindPermission(ctx context.Context, id uint64) (*domain.Permission, error)
	EditPermission(ctx context.Context, m domain.Permission) error
}

type UpdatePermission struct {
	tele      *telemetry.Telemetry
	validator validation.Validator
	store     UpdatePermissionStore
}

func NewUpdatePermission(dep Dependency, s UpdatePermissionStore) *UpdatePermission {
	return &UpdatePermission{
		tele:      dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (up *UpdatePermission) Call(ctx context.Context, in domain.UpdatePermissionInput) (
	*domain.UpdatePermissionOutput, error,
) {
	ctx, span := up.tele.Tracer().Start(ctx, "rbac.usecase.UpdatePermission")
	defer span.End()

	if err := up.validator.Validate(in); err != nil {
		up.tele.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	perm, err := up.store.FindPermission(ctx, in.ID)
	if err != nil {
		up.tele.Logger().Error(ctx, "failed to find permission by id", err)

		return nil, goerror.NewServerInternal(err)
	}

	if perm == nil {
		up.tele.Logger().Warn(ctx, "permission is not found")

		return nil, goerror.NewBusiness("permission not found", goerror.CodeNotFound)
	}

	if err := up.store.EditPermission(ctx, domain.Permission(in)); err != nil {
		up.tele.Logger().Error(ctx, "failed to update permission", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.UpdatePermissionOutput{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
	}, nil
}
