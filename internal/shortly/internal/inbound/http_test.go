package inbound

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHTTP(t *testing.T) {
	type args struct {
		router *httprouter.Router
		h      *Endpoint
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{router: &httprouter.Router{}, h: &Endpoint{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			RegisterHTTP(tt.args.router, tt.args.h)
		})
	}
}

func TestEndpoint_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		args    func() args
		want    any
		wantErr error
		mockFn  func(a args) *Endpoint
	}{
		{
			name:    "ErrorCallUC",
			args:    func() args { return args{ctx: context.TODO()} },
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *Endpoint {
				uc := new(mockz.MockGet)

				uc.EXPECT().Call(a.ctx, domain.GetInput{}).Return(nil, assert.AnError)

				return &Endpoint{
					GetUC: uc,
				}
			},
		},
		{
			name:    "ErrorNotFound",
			args:    func() args { return args{ctx: context.TODO()} },
			want:    map[string]any{"url": "NOT_FOUND"},
			wantErr: nil,
			mockFn: func(a args) *Endpoint {
				uc := new(mockz.MockGet)

				uc.EXPECT().Call(a.ctx, domain.GetInput{}).Return(&domain.GetOutput{URL: ""}, nil)

				return &Endpoint{
					GetUC: uc,
				}
			},
		},
		{
			name: "Success",
			args: func() args {
				req := httptest.NewRequest(http.MethodPost, "/shortly/value_of_key_param", nil)
				ctx := context.WithValue(context.TODO(), httprouter.ParamsKey, httprouter.Params{{Key: "key", Value: "value_of_key_param"}})

				return args{ctx: ctx, r: req}
			},
			want:    map[string]any{"url": "https://www.google.com"},
			wantErr: nil,
			mockFn: func(a args) *Endpoint {
				uc := new(mockz.MockGet)

				uc.EXPECT().Call(a.ctx, domain.GetInput{Key: "value_of_key_param"}).
					Return(&domain.GetOutput{URL: "https://www.google.com"}, nil)

				return &Endpoint{
					GetUC: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			arg := tt.args()
			e := tt.mockFn(arg)
			got, err := e.Get(arg.ctx, arg.r)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEndpoint_Set(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		args    func() args
		want    any
		wantErr error
		mockFn  func(a args) *Endpoint
	}{
		{
			name: "ErrorDecodeBody",
			args: func() args {
				body := []byte(`N/A`)
				req := httptest.NewRequest(http.MethodPost, "/shortly", bytes.NewBuffer(body))

				return args{ctx: context.TODO(), r: req}
			},
			want:    nil,
			wantErr: ErrDecodeBody,
			mockFn:  func(a args) *Endpoint { return &Endpoint{} },
		},
		{
			name: "ErrorCallUC",
			args: func() args {
				body := []byte(`{"url":"https://www.golang.com"}`)
				req := httptest.NewRequest(http.MethodPost, "/shortly", bytes.NewBuffer(body))

				return args{ctx: context.TODO(), r: req}
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *Endpoint {
				uc := new(mockz.MockSet)

				uc.EXPECT().Call(a.ctx, domain.SetInput{URL: "https://www.golang.com"}).
					Return(nil, assert.AnError)

				return &Endpoint{
					SetUC: uc,
				}
			},
		},
		{
			name: "Success",
			args: func() args {
				body := []byte(`{"url":"https://www.golang.com"}`)
				req := httptest.NewRequest(http.MethodPost, "/shortly", bytes.NewBuffer(body))

				return args{ctx: context.TODO(), r: req}
			},
			want:    map[string]any{"key": "ok"},
			wantErr: nil,
			mockFn: func(a args) *Endpoint {
				uc := new(mockz.MockSet)

				uc.EXPECT().Call(a.ctx, domain.SetInput{URL: "https://www.golang.com"}).
					Return(&domain.SetOutput{Key: "ok"}, nil)

				return &Endpoint{
					SetUC: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			arg := tt.args()
			e := tt.mockFn(arg)
			got, err := e.Set(arg.ctx, arg.r)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
