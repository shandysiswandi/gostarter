package inbound

import (
	"context"
	"testing"

	"github.com/shandysiswandi/goreng/enum"
	"github.com/shandysiswandi/goreng/telemetry"
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func Test_gqlEndpoint_Mutation(t *testing.T) {
	tests := []struct {
		name string
		e    *gqlEndpoint
		want ql.MutationResolver
	}{
		{
			name: "Success",
			e:    &gqlEndpoint{},
			want: &gqlEndpoint{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.e.Mutation())
		})
	}
}

func Test_gqlEndpoint_Query(t *testing.T) {
	tests := []struct {
		name string
		e    *gqlEndpoint
		want ql.QueryResolver
	}{
		{
			name: "Success",
			e:    &gqlEndpoint{},
			want: &gqlEndpoint{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.e.Query())
		})
	}
}

func Test_gqlEndpoint_Fetch(t *testing.T) {
	cursor := "Mg"
	limit := "1"

	type args struct {
		ctx context.Context
		in  *ql.FetchInput
	}
	tests := []struct {
		name    string
		args    args
		want    *ql.FetchOutput
		wantErr error
		mockFn  func(a args) *gqlEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				in: &ql.FetchInput{
					Cursor: &cursor,
					Limit:  &limit,
					Status: &ql.AllStatus[0],
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *gqlEndpoint {
				fetchMock := mockz.NewMockFetch(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Fetch")
				defer span.End()

				in := domain.FetchInput{
					Cursor: "Mg",
					Limit:  "1",
					Status: "UNKNOWN",
				}
				fetchMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
					tel:     tel,
					fetchUC: fetchMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in: &ql.FetchInput{
					Cursor: &cursor,
					Limit:  &limit,
					Status: &ql.AllStatus[1],
				},
			},
			want: &ql.FetchOutput{
				Todos: []ql.Todo{{
					ID:          "10",
					UserID:      "11",
					Title:       "title",
					Description: "description",
					Status:      ql.StatusDone,
				}},
				Pagination: &ql.Pagination{
					NextCursor: "NTY",
					HasNext:    true,
				},
			},
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				fetchMock := mockz.NewMockFetch(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Fetch")
				defer span.End()

				in := domain.FetchInput{
					Cursor: "Mg",
					Limit:  "1",
					Status: a.in.Status.String(),
				}
				out := &domain.FetchOutput{
					Todos: []domain.Todo{{
						ID:          10,
						UserID:      11,
						Title:       "title",
						Description: "description",
						Status:      enum.New(domain.TodoStatusDone),
					}},
					NextCursor: "NTY",
					HasMore:    true,
				}
				fetchMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
					tel:     tel,
					fetchUC: fetchMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Fetch(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_gqlEndpoint_Find(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *ql.Todo
		wantErr error
		mockFn  func(a args) *gqlEndpoint
	}{
		{
			name: "ErrorParseToUint",
			args: args{
				ctx: context.Background(),
				id:  "a",
			},
			want:    nil,
			wantErr: errFailedParseToUint,
			mockFn: func(a args) *gqlEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Find")
				defer span.End()

				return &gqlEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				id:  "10",
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *gqlEndpoint {
				findMock := mockz.NewMockFind(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Find")
				defer span.End()

				in := domain.FindInput{ID: 10}
				findMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
					tel:    tel,
					findUC: findMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  "10",
			},
			want: &ql.Todo{
				ID:          "10",
				UserID:      "11",
				Title:       "title",
				Description: "description",
				Status:      ql.StatusDrop,
			},
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				findMock := mockz.NewMockFind(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Find")
				defer span.End()

				in := domain.FindInput{ID: 10}
				out := &domain.Todo{
					ID:          10,
					UserID:      11,
					Title:       "title",
					Description: "description",
					Status:      enum.New(domain.TodoStatusDrop),
				}
				findMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
					tel:    tel,
					findUC: findMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Find(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_gqlEndpoint_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		in  ql.CreateInput
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
		mockFn  func(a args) *gqlEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				in:  ql.CreateInput{Title: "title", Description: "description"},
			},
			want:    "",
			wantErr: assert.AnError,
			mockFn: func(a args) *gqlEndpoint {
				createMock := mockz.NewMockCreate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Create")
				defer span.End()

				in := domain.CreateInput{
					Title:       "title",
					Description: "description",
				}
				createMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
					tel:      tel,
					createUC: createMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in:  ql.CreateInput{Title: "title", Description: "description"},
			},
			want:    "10",
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				createMock := mockz.NewMockCreate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Create")
				defer span.End()

				in := domain.CreateInput{
					Title:       "title",
					Description: "description",
				}
				out := &domain.CreateOutput{ID: 10}
				createMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
					tel:      tel,
					createUC: createMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Create(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_gqlEndpoint_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
		mockFn  func(a args) *gqlEndpoint
	}{
		{
			name: "ErrorParseToUint",
			args: args{
				ctx: context.Background(),
				id:  "a",
			},
			want:    "",
			wantErr: errFailedParseToUint,
			mockFn: func(a args) *gqlEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Delete")
				defer span.End()

				return &gqlEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				id:  "10",
			},
			want:    "",
			wantErr: assert.AnError,
			mockFn: func(a args) *gqlEndpoint {
				deleteMock := mockz.NewMockDelete(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Delete")
				defer span.End()

				in := domain.DeleteInput{ID: 10}
				deleteMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
					tel:      tel,
					deleteUC: deleteMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  "10",
			},
			want:    "10",
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				deleteMock := mockz.NewMockDelete(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Delete")
				defer span.End()

				in := domain.DeleteInput{ID: 10}
				out := &domain.DeleteOutput{ID: 10}
				deleteMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
					tel:      tel,
					deleteUC: deleteMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Delete(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_gqlEndpoint_UpdateStatus(t *testing.T) {
	type args struct {
		ctx context.Context
		in  ql.UpdateStatusInput
	}
	tests := []struct {
		name    string
		args    args
		want    *ql.UpdateStatusOutput
		wantErr error
		mockFn  func(a args) *gqlEndpoint
	}{
		{
			name: "ErrorParseToUint",
			args: args{
				ctx: context.Background(),
				in:  ql.UpdateStatusInput{ID: "a", Status: ql.StatusDone},
			},
			want:    nil,
			wantErr: errFailedParseToUint,
			mockFn: func(a args) *gqlEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.UpdateStatus")
				defer span.End()

				return &gqlEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				in:  ql.UpdateStatusInput{ID: "10", Status: ql.StatusDone},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *gqlEndpoint {
				updateStateMock := mockz.NewMockUpdateStatus(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.UpdateStatus")
				defer span.End()

				in := domain.UpdateStatusInput{
					ID:     10,
					Status: ql.StatusDone.String(),
				}
				updateStateMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
					tel:            tel,
					updateStatusUC: updateStateMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in:  ql.UpdateStatusInput{ID: "10", Status: ql.StatusDone},
			},
			want:    &ql.UpdateStatusOutput{ID: "10", Status: ql.StatusDone},
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				updateStateMock := mockz.NewMockUpdateStatus(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.UpdateStatus")
				defer span.End()

				in := domain.UpdateStatusInput{
					ID:     10,
					Status: ql.StatusDone.String(),
				}
				out := &domain.UpdateStatusOutput{
					ID:     10,
					Status: enum.New(domain.TodoStatusDone),
				}
				updateStateMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
					tel:            tel,
					updateStatusUC: updateStateMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).UpdateStatus(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_gqlEndpoint_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		in  ql.UpdateInput
	}
	tests := []struct {
		name    string
		args    args
		want    *ql.Todo
		wantErr error
		mockFn  func(a args) *gqlEndpoint
	}{
		{
			name: "ErrorParseToUint",
			args: args{
				ctx: context.Background(),
				in: ql.UpdateInput{
					ID:          "a",
					Title:       "title",
					Description: "description",
					Status:      ql.StatusInProgress,
				},
			},
			want:    nil,
			wantErr: errFailedParseToUint,
			mockFn: func(a args) *gqlEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Update")
				defer span.End()

				return &gqlEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				in: ql.UpdateInput{
					ID:          "10",
					Title:       "title",
					Description: "description",
					Status:      ql.StatusInProgress,
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *gqlEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      ql.StatusInProgress.String(),
				}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
					tel:      tel,
					updateUC: updateMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in: ql.UpdateInput{
					ID:          "10",
					Title:       "title",
					Description: "description",
					Status:      ql.StatusInProgress,
				},
			},
			want: &ql.Todo{
				ID:          "10",
				UserID:      "11",
				Title:       "title",
				Description: "description",
				Status:      ql.StatusInProgress,
			},
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.gqlEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      ql.StatusInProgress.String(),
				}
				out := &domain.Todo{
					ID:          10,
					UserID:      11,
					Title:       "title",
					Description: "description",
					Status:      enum.New(domain.TodoStatusInProgress),
				}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
					tel:      tel,
					updateUC: updateMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Update(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
