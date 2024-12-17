package inbound

import (
	"context"
	"testing"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func Test_grpcEndpoint_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.CreateResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateRequest{
					Title:       "title",
					Description: "description",
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				createMock := mockz.NewMockCreate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Create")
				defer span.End()

				in := domain.CreateInput{
					Title:       a.req.GetTitle(),
					Description: a.req.GetDescription(),
				}
				createMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					tel:      tel,
					createUC: createMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateRequest{
					Title:       "title",
					Description: "description",
				},
			},
			want:    &pb.CreateResponse{Id: 10},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				createMock := mockz.NewMockCreate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Create")
				defer span.End()

				in := domain.CreateInput{
					Title:       a.req.GetTitle(),
					Description: a.req.GetDescription(),
				}
				out := &domain.CreateOutput{ID: 10}
				createMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					tel:      tel,
					createUC: createMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.Create(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_grpcEndpoint_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DeleteResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.DeleteRequest{Id: 10},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				deleteMock := mockz.NewMockDelete(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Delete")
				defer span.End()

				in := domain.DeleteInput{ID: a.req.GetId()}
				deleteMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					tel:      tel,
					deleteUC: deleteMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.DeleteRequest{Id: 10},
			},
			want:    &pb.DeleteResponse{Id: 10},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				deleteMock := mockz.NewMockDelete(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Delete")
				defer span.End()

				in := domain.DeleteInput{ID: a.req.GetId()}
				out := &domain.DeleteOutput{ID: 10}
				deleteMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					tel:      tel,
					deleteUC: deleteMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.Delete(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_grpcEndpoint_Find(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.FindRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.FindResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.FindRequest{Id: 10},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				findMock := mockz.NewMockFind(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Find")
				defer span.End()

				in := domain.FindInput{ID: a.req.GetId()}
				findMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					tel:    tel,
					findUC: findMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.FindRequest{Id: 10},
			},
			want: &pb.FindResponse{
				Id:          10,
				UserId:      11,
				Title:       "title",
				Description: "description",
				Status:      pb.Status_STATUS_DONE,
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				findMock := mockz.NewMockFind(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Find")
				defer span.End()

				in := domain.FindInput{ID: a.req.GetId()}
				out := &domain.Todo{
					ID:          10,
					UserID:      11,
					Title:       "title",
					Description: "description",
					Status:      enum.New(domain.TodoStatusDone),
				}
				findMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					tel:    tel,
					findUC: findMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.Find(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_grpcEndpoint_Fetch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.FetchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.FetchResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.FetchRequest{
					Cursor: "Mg",
					Limit:  "1",
					Status: pb.Status_STATUS_DONE,
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				fetchMock := mockz.NewMockFetch(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Fetch")
				defer span.End()

				in := domain.FetchInput{
					Cursor: "Mg",
					Limit:  "1",
					Status: a.req.Status.String(),
				}
				fetchMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					tel:     tel,
					fetchUC: fetchMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.FetchRequest{
					Cursor: "Mg",
					Limit:  "1",
					Status: pb.Status_STATUS_DONE,
				},
			},
			want: &pb.FetchResponse{
				Todos: []*pb.Todo{{
					Id:          10,
					UserId:      11,
					Title:       "title",
					Description: "description",
					Status:      pb.Status_STATUS_DONE,
				}},
				Pagination: &pb.Pagination{
					NextCursor: "NTY",
					HasMore:    true,
				},
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				fetchMock := mockz.NewMockFetch(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Fetch")
				defer span.End()

				in := domain.FetchInput{
					Cursor: "Mg",
					Limit:  "1",
					Status: a.req.Status.String(),
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

				return &grpcEndpoint{
					tel:     tel,
					fetchUC: fetchMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.Fetch(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_grpcEndpoint_UpdateStatus(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateStatusRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateStatusResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateStatusRequest{
					Id:     10,
					Status: pb.Status_STATUS_DONE,
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				updateStatusMock := mockz.NewMockUpdateStatus(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.UpdateStatus")
				defer span.End()

				in := domain.UpdateStatusInput{
					ID:     10,
					Status: "DONE",
				}
				updateStatusMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					tel:            tel,
					updateStatusUC: updateStatusMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateStatusRequest{
					Id:     10,
					Status: pb.Status_STATUS_DONE,
				},
			},
			want: &pb.UpdateStatusResponse{
				Id:     10,
				Status: pb.Status_STATUS_DONE,
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				updateStatusMock := mockz.NewMockUpdateStatus(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.UpdateStatus")
				defer span.End()

				in := domain.UpdateStatusInput{
					ID:     10,
					Status: "DONE",
				}
				out := &domain.UpdateStatusOutput{
					ID:     10,
					Status: enum.New(domain.TodoStatusDone),
				}
				updateStatusMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					tel:            tel,
					updateStatusUC: updateStatusMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.UpdateStatus(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_grpcEndpoint_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateResponse
		wantErr error
		mockFn  func(a args) *grpcEndpoint
	}{
		{
			name: "ErrorCallUC",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateRequest{
					Id:          10,
					Title:       "title",
					Description: "description",
					Status:      pb.Status_STATUS_DROP,
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      "DROP",
				}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					tel:      tel,
					updateUC: updateMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateRequest{
					Id:          10,
					Title:       "title",
					Description: "description",
					Status:      pb.Status_STATUS_DROP,
				},
			},
			want: &pb.UpdateResponse{
				Id:          10,
				UserId:      11,
				Title:       "title",
				Description: "description",
				Status:      pb.Status_STATUS_DROP,
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.inbound.grpcEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      "DROP",
				}
				out := &domain.Todo{
					ID:          10,
					UserID:      11,
					Title:       "title",
					Description: "description",
					Status:      enum.New(domain.TodoStatusDrop),
				}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
					tel:      tel,
					updateUC: updateMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := tt.mockFn(tt.args)
			got, err := g.Update(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
