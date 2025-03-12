package inbound

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/shandysiswandi/goreng/enum"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/stretchr/testify/assert"
)

func Test_httpEndpoint_Create(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPost, "/todos", body)

				return c.Build()
			},
			want:    nil,
			wantErr: errInvalidBody,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Create")
				defer span.End()

				return &httpEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"title":"title","description":"description"}`)
				c := framework.NewTestContext(http.MethodPost, "/todos", body)

				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				createMock := mockz.NewMockCreate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Create")
				defer span.End()

				in := domain.CreateInput{
					Title:       "title",
					Description: "description",
				}
				createMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:      tel,
					createUC: createMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"title":"title","description":"description"}`)
				c := framework.NewTestContext(http.MethodPost, "/todos", body)

				return c.Build()
			},
			want:    CreateResponse{ID: 12},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				createMock := mockz.NewMockCreate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Create")
				defer span.End()

				in := domain.CreateInput{
					Title:       "title",
					Description: "description",
				}
				out := &domain.CreateOutput{ID: 12}
				createMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:      tel,
					createUC: createMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Create(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_Delete(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorParseToUint",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodDelete, "/todos/1", nil)
				c.SetParam("id", "n/a")

				return c.Build()
			},
			want:    nil,
			wantErr: errFailedParseToUint,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Delete")
				defer span.End()

				return &httpEndpoint{
					tel:    tel,
					findUC: nil,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodPost, "/todos/1", nil)
				c.SetParam("id", "1")

				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				deleteMock := mockz.NewMockDelete(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Delete")
				defer span.End()

				in := domain.DeleteInput{ID: 1}
				deleteMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:      tel,
					deleteUC: deleteMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodPost, "/todos/1", nil)
				c.SetParam("id", "1")

				return c.Build()
			},
			want:    DeleteResponse{ID: 1},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				deleteMock := mockz.NewMockDelete(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Delete")
				defer span.End()

				in := domain.DeleteInput{ID: 1}
				out := &domain.DeleteOutput{ID: 1}
				deleteMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:      tel,
					deleteUC: deleteMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Delete(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_Find(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorParseToUint",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/todos/1", nil)
				c.SetParam("id", "n/a")

				return c.Build()
			},
			want:    nil,
			wantErr: errFailedParseToUint,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Find")
				defer span.End()

				return &httpEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/todos/1", nil)
				c.SetParam("id", "11")

				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				findMock := mockz.NewMockFind(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Find")
				defer span.End()

				in := domain.FindInput{ID: 11}
				findMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:    tel,
					findUC: findMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/todos/1", nil)
				c.SetParam("id", "11")

				return c.Build()
			},
			want: FindResponse{
				ID:          11,
				UserID:      89,
				Title:       "title",
				Description: "description",
				Status:      enum.New(domain.TodoStatusDone).String(),
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				findMock := mockz.NewMockFind(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Find")
				defer span.End()

				in := domain.FindInput{ID: 11}
				out := &domain.Todo{
					ID:          11,
					UserID:      89,
					Title:       "title",
					Description: "description",
					Status:      enum.New(domain.TodoStatusDone),
				}
				findMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:    tel,
					findUC: findMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Find(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_Fetch(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/todos", nil)
				c.SetQuery("limit", "1")
				c.SetQuery("cursor", "Mg")
				c.SetQuery("status", "done")

				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				fetchMock := mockz.NewMockFetch(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Fetch")
				defer span.End()

				in := domain.FetchInput{
					Cursor: "Mg",
					Limit:  "1",
					Status: "done",
				}
				fetchMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:     tel,
					fetchUC: fetchMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/todos", nil)
				c.SetQuery("limit", "1")
				c.SetQuery("cursor", "Mg")
				c.SetQuery("status", "done")

				return c.Build()
			},
			want: FetchResponse{
				Todos: []Todo{{
					ID:          11,
					Title:       "title",
					Description: "description",
					Status:      enum.New(domain.TodoStatusDone).String(),
				}},
				Pagination: Pagination{
					NextCursor: "NTY",
					HasMore:    true,
				},
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				fetchMock := mockz.NewMockFetch(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Fetch")
				defer span.End()

				in := domain.FetchInput{
					Cursor: "Mg",
					Limit:  "1",
					Status: "done",
				}
				out := &domain.FetchOutput{
					Todos: []domain.Todo{{
						ID:          11,
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

				return &httpEndpoint{
					tel:     tel,
					fetchUC: fetchMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Fetch(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_UpdateStatus(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorParseToUint",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodPatch, "/todos/1", nil)
				c.SetParam("id", "n/a")

				return c.Build()
			},
			want:    nil,
			wantErr: errFailedParseToUint,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.UpdateStatus")
				defer span.End()

				return &httpEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPatch, "/todos/1", body)
				c.SetParam("id", "2")

				return c.Build()
			},
			want:    nil,
			wantErr: errInvalidBody,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.UpdateStatus")
				defer span.End()

				return &httpEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"status":"done"}`)
				c := framework.NewTestContext(http.MethodPatch, "/todos/2", body)
				c.SetParam("id", "2")

				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				updateStatusMock := mockz.NewMockUpdateStatus(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.UpdateStatus")
				defer span.End()

				in := domain.UpdateStatusInput{ID: 2, Status: "done"}
				updateStatusMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:            tel,
					updateStatusUC: updateStatusMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"status":"done"}`)
				c := framework.NewTestContext(http.MethodPatch, "/todos/2", body)
				c.SetParam("id", "2")

				return c.Build()
			},
			want:    UpdateStatusResponse{ID: 2, Status: enum.New(domain.TodoStatusDone).String()},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				updateStatusMock := mockz.NewMockUpdateStatus(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.UpdateStatus")
				defer span.End()

				in := domain.UpdateStatusInput{ID: 2, Status: "done"}
				out := &domain.UpdateStatusOutput{ID: 2, Status: enum.New(domain.TodoStatusDone)}
				updateStatusMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:            tel,
					updateStatusUC: updateStatusMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.UpdateStatus(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_Update(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorParseToUint",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodPut, "/todos/1", nil)
				c.SetParam("id", "n/a")

				return c.Build()
			},
			want:    nil,
			wantErr: errFailedParseToUint,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.UpdateStatus")
				defer span.End()

				return &httpEndpoint{
					tel:    tel,
					findUC: nil,
				}
			},
		},
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPut, "/todos/2", body)
				c.SetParam("id", "2")

				return c.Build()
			},
			want:    nil,
			wantErr: errInvalidBody,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Update")
				defer span.End()

				return &httpEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"title":"title","description":"description","status":"done"}`)
				c := framework.NewTestContext(http.MethodPut, "/todos/2", body)
				c.SetParam("id", "2")

				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{
					ID:          2,
					Title:       "title",
					Description: "description",
					Status:      "done",
				}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:      tel,
					updateUC: updateMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"title":"title","description":"description","status":"done"}`)
				c := framework.NewTestContext(http.MethodPut, "/todos/2", body)
				c.SetParam("id", "2")

				return c.Build()
			},
			want: UpdateResponse{
				ID:          2,
				UserID:      12,
				Title:       "title",
				Description: "description",
				Status:      enum.New(domain.TodoStatusDone).String(),
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "todo.inbound.httpEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{
					ID:          2,
					Title:       "title",
					Description: "description",
					Status:      "done",
				}
				out := &domain.Todo{
					ID:          2,
					UserID:      12,
					Title:       "title",
					Description: "description",
					Status:      enum.New(domain.TodoStatusDone),
				}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:      tel,
					updateUC: updateMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Update(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
