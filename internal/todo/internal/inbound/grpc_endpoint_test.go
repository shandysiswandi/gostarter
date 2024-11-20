package inbound

import (
	"context"
	"testing"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
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
				req: &pb.CreateRequest{Title: "title", Description: "description"},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				createMock := mockz.NewMockCreate(t)

				in := domain.CreateInput{
					Title:       a.req.GetTitle(),
					Description: a.req.GetDescription(),
				}
				createMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					createUC: createMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateRequest{Title: "title", Description: "description"},
			},
			want:    &pb.CreateResponse{Id: 10},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				createMock := mockz.NewMockCreate(t)

				in := domain.CreateInput{
					Title:       a.req.GetTitle(),
					Description: a.req.GetDescription(),
				}
				out := &domain.CreateOutput{ID: 10}
				createMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
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

				in := domain.DeleteInput{ID: a.req.GetId()}
				deleteMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
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

				in := domain.DeleteInput{ID: a.req.GetId()}
				out := &domain.DeleteOutput{ID: 10}
				deleteMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
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

				in := domain.FindInput{ID: a.req.GetId()}
				findMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
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
				Title:       "title",
				Description: "description",
				Status:      pb.Status_STATUS_DONE,
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				findMock := mockz.NewMockFind(t)

				in := domain.FindInput{ID: a.req.GetId()}
				out := &domain.Todo{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      domain.TodoStatusDone,
				}
				findMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
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
					Id:          10,
					Title:       "title",
					Description: "description",
					Status:      pb.Status_STATUS_DONE,
				},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				fetchMock := mockz.NewMockFetch(t)

				in := domain.FetchInput{
					ID:          "10",
					Title:       "title",
					Description: "description",
					Status:      "STATUS_DONE",
				}
				fetchMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					fetchUC: fetchMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.FetchRequest{
					Id:          10,
					Title:       "title",
					Description: "description",
					Status:      pb.Status_STATUS_DONE,
				},
			},
			want: &pb.FetchResponse{Todos: []*pb.Todo{
				{
					Id:          10,
					Title:       "title",
					Description: "description",
					Status:      pb.Status_STATUS_DONE,
				},
			}},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				fetchMock := mockz.NewMockFetch(t)

				in := domain.FetchInput{
					ID:          "10",
					Title:       "title",
					Description: "description",
					Status:      "STATUS_DONE",
				}
				out := []domain.Todo{{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      domain.TodoStatusDone,
				}}
				fetchMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
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
				req: &pb.UpdateStatusRequest{Id: 10, Status: pb.Status_STATUS_DONE},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *grpcEndpoint {
				updateStatusMock := mockz.NewMockUpdateStatus(t)

				in := domain.UpdateStatusInput{ID: 10, Status: "STATUS_DONE"}
				updateStatusMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
					updateStatusUC: updateStatusMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateStatusRequest{Id: 10, Status: pb.Status_STATUS_DONE},
			},
			want:    &pb.UpdateStatusResponse{Id: 10, Status: pb.Status_STATUS_DONE},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				updateStatusMock := mockz.NewMockUpdateStatus(t)

				in := domain.UpdateStatusInput{ID: 10, Status: "STATUS_DONE"}
				out := &domain.UpdateStatusOutput{ID: 10, Status: domain.TodoStatusDone}
				updateStatusMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
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

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      "STATUS_DROP",
				}
				updateMock.EXPECT().
					Call(a.ctx, in).
					Return(nil, assert.AnError)

				return &grpcEndpoint{
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
				Title:       "title",
				Description: "description",
				Status:      pb.Status_STATUS_DROP,
			},
			wantErr: nil,
			mockFn: func(a args) *grpcEndpoint {
				updateMock := mockz.NewMockUpdate(t)

				in := domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      "STATUS_DROP",
				}
				out := &domain.Todo{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      domain.TodoStatusDrop,
				}
				updateMock.EXPECT().
					Call(a.ctx, in).
					Return(out, nil)

				return &grpcEndpoint{
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
