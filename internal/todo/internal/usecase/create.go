package usecase

import (
	"context"
	"errors"

	"github.com/shandysiswandi/goreng/enum"
	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type CreateStore interface {
	Create(ctx context.Context, in domain.Todo) error
}

type Create struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	store     CreateStore
}

func NewCreate(dep Dependency, s CreateStore) *Create {
	return &Create{
		telemetry: dep.Telemetry,
		uidnumber: dep.UIDNumber,
		validator: dep.Validator,
		store:     s,
	}
}

func (s *Create) Call(ctx context.Context, in domain.CreateInput) (*domain.CreateOutput, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "todo.usecase.Create")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	id := s.uidnumber.Generate()
	userID := uint64(0)
	if clm := lib.GetJWTClaim(ctx); clm != nil {
		userID = clm.AuthID
	}

	err := s.store.Create(ctx, domain.Todo{
		ID:          id,
		UserID:      userID,
		Title:       in.Title,
		Description: in.Description,
		Status:      enum.New(domain.TodoStatusInitiate),
	})
	if errors.Is(err, domain.ErrTodoNotCreated) {
		s.telemetry.Logger().Warn(ctx, "todo created but db not affected")

		return nil, goerror.NewBusiness("failed to create todo", goerror.CodeUnknown)
	}

	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to create", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.CreateOutput{ID: id}, nil
}
