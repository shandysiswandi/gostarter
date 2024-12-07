package framework

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterCtx_Context(t *testing.T) {
	tests := []struct {
		name string
		rc   *RouterCtx
		want context.Context
	}{
		{
			name: "Success",
			rc: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.rc.Context()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRouterCtx_Request(t *testing.T) {
	tests := []struct {
		name string
		rc   *RouterCtx
		want *http.Request
	}{
		{
			name: "Nil",
			rc:   &RouterCtx{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.rc.Request()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRouterCtx_Header(t *testing.T) {
	tests := []struct {
		name string
		rc   *RouterCtx
		want http.Header
	}{
		{
			name: "Success",
			rc: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: map[string][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.rc.Header()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRouterCtx_Body(t *testing.T) {
	tests := []struct {
		name string
		rc   *RouterCtx
		want io.ReadCloser
	}{
		{
			name: "Success",
			rc: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: http.NoBody,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.rc.Body()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRouterCtx_Query(t *testing.T) {
	tests := []struct {
		name string
		key  string
		rc   *RouterCtx
		want string
	}{
		{
			name: "Success",
			key:  "key",
			rc: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.rc.Query(tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRouterCtx_Queries(t *testing.T) {
	tests := []struct {
		name string
		rc   *RouterCtx
		key  string
		want []string
	}{
		{
			name: "Success",
			rc: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			key:  "key",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.rc.Queries(tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRouterCtx_Param(t *testing.T) {
	tests := []struct {
		name string
		rc   *RouterCtx
		key  string
		want string
	}{
		{
			name: "Success",
			rc: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			key:  "key",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.rc.Param(tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewTestContext(t *testing.T) {
	type args struct {
		method string
		target string
		body   io.Reader
	}
	tests := []struct {
		name string
		args args
		want *TestCtx
	}{
		{
			name: "Success",
			args: args{
				method: "GET",
				target: "/",
				body:   nil,
			},
			want: &TestCtx{
				r:       httptest.NewRequest("GET", "/", nil),
				param:   make(map[string]string),
				queries: make(map[string][]string),
				headers: make(map[string][]string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewTestContext(tt.args.method, tt.args.target, tt.args.body)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTestCtx_SetParam(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		tc   *TestCtx
		args args
	}{
		{
			name: "Success",
			tc: &TestCtx{
				param: make(map[string]string),
			},
			args: args{
				key:   "key",
				value: "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.tc.SetParam(tt.args.key, tt.args.value)
		})
	}
}

func TestTestCtx_SetQuery(t *testing.T) {
	type args struct {
		key    string
		values []string
	}
	tests := []struct {
		name string
		tc   *TestCtx
		args args
	}{
		{
			name: "Success",
			tc: &TestCtx{
				queries: make(map[string][]string),
			},
			args: args{
				key:    "key",
				values: []string{"value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.tc.SetQuery(tt.args.key, tt.args.values...)
		})
	}
}

func TestTestCtx_SetHeader(t *testing.T) {
	type args struct {
		key    string
		values []string
	}
	tests := []struct {
		name string
		tc   *TestCtx
		args args
	}{
		{
			name: "Success",
			tc: &TestCtx{
				headers: make(map[string][]string),
			},
			args: args{
				key:    "key",
				values: []string{"value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.tc.SetHeader(tt.args.key, tt.args.values...)
		})
	}
}

func TestTestCtx_SetContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		tc   *TestCtx
		args args
	}{
		{
			name: "Success",
			tc: &TestCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.tc.SetContext(tt.args.ctx)
		})
	}
}

func TestTestCtx_Build(t *testing.T) {
	tests := []struct {
		name string
		tc   func() *TestCtx
		want *RouterCtx
	}{
		{
			name: "SuccessQuery",
			tc: func() *TestCtx {
				tc := NewTestContext(http.MethodGet, "/", nil)
				tc.SetQuery("key", "value")
				return tc
			},
			want: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
		},
		{
			name: "SuccessHeader",
			tc: func() *TestCtx {
				tc := NewTestContext(http.MethodGet, "/", nil)
				tc.SetHeader("key", "value")
				return tc
			},
			want: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
		},
		{
			name: "SuccessParam",
			tc: func() *TestCtx {
				tc := NewTestContext(http.MethodGet, "/", nil)
				tc.SetParam("key", "value")
				return tc
			},
			want: &RouterCtx{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.tc().Build()
			assert.Equal(t, tt.want.r.Method, got.r.Method)
		})
	}
}
