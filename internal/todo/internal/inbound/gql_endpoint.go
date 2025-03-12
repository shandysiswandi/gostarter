package inbound

import (
	"context"
	"strconv"

	"github.com/shandysiswandi/goreng/telemetry"
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

func getString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}

	return ""
}

func getStatusString(status *ql.Status) string {
	if status != nil && status.IsValid() {
		return status.String()
	}

	return ""
}

type gqlEndpoint struct {
	ql.Resolver

	tel *telemetry.Telemetry

	findUC         domain.Find
	fetchUC        domain.Fetch
	createUC       domain.Create
	deleteUC       domain.Delete
	updateUC       domain.Update
	updateStatusUC domain.UpdateStatus
}

func (l *gqlEndpoint) Mutation() ql.MutationResolver { return l }

func (l *gqlEndpoint) Query() ql.QueryResolver { return l }

func (l *gqlEndpoint) Fetch(ctx context.Context, in *ql.FetchInput) (*ql.FetchOutput, error) {
	ctx, span := l.tel.Tracer().Start(ctx, "todo.inbound.gqlEndpoint.Fetch")
	defer span.End()

	input := domain.FetchInput{}
	if in != nil {
		input = domain.FetchInput{
			Cursor: getString(in.Cursor),
			Limit:  getString(in.Limit),
			Status: getStatusString(in.Status),
		}
	}

	resp, err := l.fetchUC.Call(ctx, input)
	if err != nil {
		return nil, err
	}

	todos := make([]ql.Todo, 0)
	for _, todo := range resp.Todos {
		todos = append(todos, ql.Todo{
			ID:          strconv.FormatUint(todo.ID, 10),
			UserID:      strconv.FormatUint(todo.UserID, 10),
			Title:       todo.Title,
			Description: todo.Description,
			Status:      ql.Status(todo.Status.String()),
		})
	}

	return &ql.FetchOutput{
		Todos: todos,
		Pagination: &ql.Pagination{
			NextCursor: resp.NextCursor,
			HasNext:    resp.HasMore,
		},
	}, nil
}

func (l *gqlEndpoint) Find(ctx context.Context, id string) (*ql.Todo, error) {
	ctx, span := l.tel.Tracer().Start(ctx, "todo.inbound.gqlEndpoint.Find")
	defer span.End()

	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := l.findUC.Call(ctx, domain.FindInput{ID: idu64})
	if err != nil {
		return nil, err
	}

	return &ql.Todo{
		ID:          strconv.FormatUint(resp.ID, 10),
		UserID:      strconv.FormatUint(resp.UserID, 10),
		Title:       resp.Title,
		Description: resp.Description,
		Status:      ql.Status(resp.Status.String()),
	}, nil
}

func (l *gqlEndpoint) Create(ctx context.Context, in ql.CreateInput) (string, error) {
	ctx, span := l.tel.Tracer().Start(ctx, "todo.inbound.gqlEndpoint.Create")
	defer span.End()

	resp, err := l.createUC.Call(ctx, domain.CreateInput{Title: in.Title, Description: in.Description})
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(resp.ID, 10), nil
}

func (l *gqlEndpoint) Delete(ctx context.Context, id string) (string, error) {
	ctx, span := l.tel.Tracer().Start(ctx, "todo.inbound.gqlEndpoint.Delete")
	defer span.End()

	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", errFailedParseToUint
	}

	resp, err := l.deleteUC.Call(ctx, domain.DeleteInput{ID: idu64})
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(resp.ID, 10), nil
}

func (l *gqlEndpoint) UpdateStatus(ctx context.Context, in ql.UpdateStatusInput) (
	*ql.UpdateStatusOutput, error,
) {
	ctx, span := l.tel.Tracer().Start(ctx, "todo.inbound.gqlEndpoint.UpdateStatus")
	defer span.End()

	idu64, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := l.updateStatusUC.Call(ctx, domain.UpdateStatusInput{ID: idu64, Status: in.Status.String()})
	if err != nil {
		return nil, err
	}

	return &ql.UpdateStatusOutput{
		ID:     strconv.FormatUint(resp.ID, 10),
		Status: ql.Status(resp.Status.String()),
	}, nil
}

func (l *gqlEndpoint) Update(ctx context.Context, in ql.UpdateInput) (*ql.Todo, error) {
	ctx, span := l.tel.Tracer().Start(ctx, "todo.inbound.gqlEndpoint.Update")
	defer span.End()

	idu64, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := l.updateUC.Call(ctx, domain.UpdateInput{
		ID:          idu64,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status.String(),
	})
	if err != nil {
		return nil, err
	}

	return &ql.Todo{
		ID:          strconv.FormatUint(resp.ID, 10),
		UserID:      strconv.FormatUint(resp.UserID, 10),
		Title:       resp.Title,
		Description: resp.Description,
		Status:      ql.Status(resp.Status.String()),
	}, nil
}
