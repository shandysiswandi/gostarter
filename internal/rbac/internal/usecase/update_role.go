//nolint:dupl // this is not duplicate
package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
)

type UpdateRoleStore interface {
	FindRole(ctx context.Context, id uint64) (*domain.Role, error)
	EditRole(ctx context.Context, m domain.Role) error
}

type UpdateRole struct {
	tele      *telemetry.Telemetry
	validator validation.Validator
	store     UpdateRoleStore
}

func NewUpdateRole(dep Dependency, s UpdateRoleStore) *UpdateRole {
	return &UpdateRole{
		tele:      dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (ur *UpdateRole) Call(ctx context.Context, in domain.UpdateRoleInput) (*domain.UpdateRoleOutput, error) {
	ctx, span := ur.tele.Tracer().Start(ctx, "rbac.usecase.UpdateRole")
	defer span.End()

	if err := ur.validator.Validate(in); err != nil {
		ur.tele.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	role, err := ur.store.FindRole(ctx, in.ID)
	if err != nil {
		ur.tele.Logger().Error(ctx, "failed to find role by id", err)

		return nil, goerror.NewServerInternal(err)
	}

	if role == nil {
		ur.tele.Logger().Warn(ctx, "role is not found")

		return nil, goerror.NewBusiness("role not found", goerror.CodeNotFound)
	}

	if err := ur.store.EditRole(ctx, domain.Role(in)); err != nil {
		ur.tele.Logger().Error(ctx, "failed to update role", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.UpdateRoleOutput{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
	}, nil
}
