package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/region/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/region/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRESTEndpoint(t *testing.T) {
	type args struct {
		router *httprouter.Router
		h      *Endpoint
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Success", args: args{router: httprouter.New()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterRESTEndpoint(tt.args.router, tt.args.h)
		})
	}
}

func TestNewEndpoint(t *testing.T) {
	type args struct {
		search domain.Search
	}
	tests := []struct {
		name string
		args args
		want *Endpoint
	}{
		{name: "success", args: args{}, want: &Endpoint{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewEndpoint(tt.args.search)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEndpoint_Search(t *testing.T) {
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
			name: "ErrorSearchCall",
			args: func() args {
				r := httptest.NewRequest(http.MethodGet, "/regions", nil)
				q := r.URL.Query()
				q.Add("by", "cities")
				q.Add("pid", "11")
				q.Add("ids", "1,2,3")
				r.URL.RawQuery = q.Encode()

				return args{
					ctx: context.TODO(),
					r:   r,
				}
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *Endpoint {
				search := mockz.NewMockSearch(t)

				search.EXPECT().
					Call(a.ctx, domain.SearchInput{By: "cities", ParentID: "11", IDs: "1,2,3"}).
					Return(nil, assert.AnError)

				return &Endpoint{
					search: search,
				}
			},
		},
		{
			name: "Success",
			args: func() args {
				r := httptest.NewRequest(http.MethodGet, "/regions", nil)
				q := r.URL.Query()
				q.Add("by", "cities")
				q.Add("pid", "11")
				q.Add("ids", "1,2,3")
				r.URL.RawQuery = q.Encode()

				return args{
					ctx: context.TODO(),
					r:   r,
				}
			},
			want:    SearchResponse{Type: "cities", Regions: []Region{{ID: "1", Name: "test 1"}}},
			wantErr: nil,
			mockFn: func(a args) *Endpoint {
				search := mockz.NewMockSearch(t)

				search.EXPECT().
					Call(a.ctx, domain.SearchInput{By: "cities", ParentID: "11", IDs: "1,2,3"}).
					Return([]domain.Region{{ID: "1", Name: "test 1"}}, nil)

				return &Endpoint{
					search: search,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := tt.args()
			e := tt.mockFn(arg)
			got, err := e.Search(arg.ctx, arg.r)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
