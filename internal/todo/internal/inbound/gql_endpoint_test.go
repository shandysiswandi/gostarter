package inbound

import (
	"context"
	"testing"

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
			got := tt.e.Mutation()
			assert.Equal(t, tt.want, got)
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
			got := tt.e.Query()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_gqlEndpoint_Fetch(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *ql.FetchInput
	}
	tests := []struct {
		name    string
		args    args
		want    []ql.Todo
		wantErr error
		mockFn  func(a args) *gqlEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				in:  &ql.FetchInput{},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *gqlEndpoint {
				fetchMock := mockz.NewMockFetch(t)

				in := domain.FetchInput{}
				fetchMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
					fetchUC: fetchMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in: &ql.FetchInput{
					ID:          new(string),
					Title:       new(string),
					Description: new(string),
					Status:      &ql.AllStatus[1],
				},
			},
			want: []ql.Todo{{
				ID:          "10",
				Title:       "title",
				Description: "description",
				Status:      ql.StatusDone,
			}},
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				fetchMock := mockz.NewMockFetch(t)

				in := domain.FetchInput{Status: a.in.Status.String()}
				out := []domain.Todo{{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      domain.TodoStatusDone,
				}}
				fetchMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
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
				return &gqlEndpoint{}
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

				in := domain.FindInput{ID: 10}
				findMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
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
				Title:       "title",
				Description: "description",
				Status:      ql.StatusDrop,
			},
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				findMock := mockz.NewMockFind(t)

				in := domain.FindInput{ID: 10}
				out := &domain.Todo{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      domain.TodoStatusDrop,
				}
				findMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
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

				in := domain.CreateInput{Title: "title", Description: "description"}
				createMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
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

				in := domain.CreateInput{Title: "title", Description: "description"}
				out := &domain.CreateOutput{ID: 10}
				createMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
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
				return &gqlEndpoint{}
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

				in := domain.DeleteInput{ID: 10}
				deleteMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
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

				in := domain.DeleteInput{ID: 10}
				out := &domain.DeleteOutput{ID: 10}
				deleteMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
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
				return &gqlEndpoint{}
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

				in := domain.UpdateStatusInput{ID: 10, Status: ql.StatusDone.String()}
				updateStateMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
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

				in := domain.UpdateStatusInput{ID: 10, Status: ql.StatusDone.String()}
				out := &domain.UpdateStatusOutput{ID: 10, Status: domain.TodoStatusDone}
				updateStateMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
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
				return &gqlEndpoint{}
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

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      ql.StatusInProgress.String(),
				}
				updateMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &gqlEndpoint{
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
				Title:       "title",
				Description: "description",
				Status:      ql.StatusInProgress,
			},
			wantErr: nil,
			mockFn: func(a args) *gqlEndpoint {
				updateMock := mockz.NewMockUpdate(t)

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      ql.StatusInProgress.String(),
				}
				out := &domain.Todo{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      domain.TodoStatusInProgress,
				}
				updateMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &gqlEndpoint{
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
