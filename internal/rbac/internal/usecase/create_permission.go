package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type CreatePermissionStore interface {
	FindPermissionByName(ctx context.Context, name string) (*domain.Permission, error)
	SavePermission(ctx context.Context, m domain.Permission) error
}

type CreatePermission struct {
	tele      *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	store     CreatePermissionStore
}

func NewCreatePermission(dep Dependency, s CreatePermissionStore) *CreatePermission {
	return &CreatePermission{
		tele:      dep.Telemetry,
		validator: dep.Validator,
		uidnumber: dep.UIDNumber,
		store:     s,
	}
}

func (cr *CreatePermission) Call(ctx context.Context, in domain.CreatePermissionInput) (
	*domain.CreatePermissionOutput, error) {
	ctx, span := cr.tele.Tracer().Start(ctx, "rbac.usecase.CreatePermission")
	defer span.End()

	if err := cr.validator.Validate(in); err != nil {
		cr.tele.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	perm, err := cr.store.FindPermissionByName(ctx, in.Name)
	if err != nil {
		cr.tele.Logger().Error(ctx, "failed to find permission by name", err)

		return nil, goerror.NewServerInternal(err)
	}
	if perm != nil {
		return &domain.CreatePermissionOutput{ID: perm.ID}, nil
	}

	permData := domain.Permission{
		ID:          cr.uidnumber.Generate(),
		Name:        in.Name,
		Description: in.Description,
	}
	if err := cr.store.SavePermission(ctx, permData); err != nil {
		cr.tele.Logger().Error(ctx, "failed to save permission", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.CreatePermissionOutput{ID: permData.ID}, nil
}
