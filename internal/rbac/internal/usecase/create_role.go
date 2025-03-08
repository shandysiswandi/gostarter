//nolint:dupl // this is not duplicate
package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
)

type CreateRoleStore interface {
	FindRoleByName(ctx context.Context, name string) (*domain.Role, error)
	SaveRole(ctx context.Context, m domain.Role) error
}

type CreateRole struct {
	tele      *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	store     CreateRoleStore
}

func NewCreateRole(dep Dependency, s CreateRoleStore) *CreateRole {
	return &CreateRole{
		tele:      dep.Telemetry,
		validator: dep.Validator,
		uidnumber: dep.UIDNumber,
		store:     s,
	}
}

func (cr *CreateRole) Call(ctx context.Context, in domain.CreateRoleInput) (*domain.CreateRoleOutput, error) {
	ctx, span := cr.tele.Tracer().Start(ctx, "rbac.usecase.CreateRole")
	defer span.End()

	if err := cr.validator.Validate(in); err != nil {
		cr.tele.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	role, err := cr.store.FindRoleByName(ctx, in.Name)
	if err != nil {
		cr.tele.Logger().Error(ctx, "failed to find role by name", err)

		return nil, goerror.NewServerInternal(err)
	}
	if role != nil {
		return &domain.CreateRoleOutput{ID: role.ID}, nil
	}

	roleData := domain.Role{
		ID:          cr.uidnumber.Generate(),
		Name:        in.Name,
		Description: in.Description,
	}
	if err := cr.store.SaveRole(ctx, roleData); err != nil {
		cr.tele.Logger().Error(ctx, "failed to save role", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.CreateRoleOutput{ID: roleData.ID}, nil
}
