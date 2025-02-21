//nolint:dupl // this is not duplicate
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

func (cp *CreatePermission) Call(ctx context.Context, in domain.CreatePermissionInput) (
	*domain.CreatePermissionOutput, error,
) {
	ctx, span := cp.tele.Tracer().Start(ctx, "rbac.usecase.CreatePermission")
	defer span.End()

	if err := cp.validator.Validate(in); err != nil {
		cp.tele.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	perm, err := cp.store.FindPermissionByName(ctx, in.Name)
	if err != nil {
		cp.tele.Logger().Error(ctx, "failed to find permission by name", err)

		return nil, goerror.NewServerInternal(err)
	}
	if perm != nil {
		return &domain.CreatePermissionOutput{ID: perm.ID}, nil
	}

	permData := domain.Permission{
		ID:          cp.uidnumber.Generate(),
		Name:        in.Name,
		Description: in.Description,
	}
	if err := cp.store.SavePermission(ctx, permData); err != nil {
		cp.tele.Logger().Error(ctx, "failed to save permission", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.CreatePermissionOutput{ID: permData.ID}, nil
}
