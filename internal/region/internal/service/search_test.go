package service

import (
	"context"
	"strings"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/region/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/region/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	vMock "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func testconvertArgs(ids string) []any {
	idss := strings.Split(ids, ",")

	var anyArgs []any

	for _, arg := range idss {
		anyArgs = append(anyArgs, arg)
	}

	return anyArgs
}

func TestNewSearch(t *testing.T) {
	type args struct {
		validate validation.Validator
		store    SearchStore
	}
	tests := []struct {
		name string
		args args
		want *Search
	}{
		{name: "Success", args: args{}, want: &Search{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSearch(tt.args.validate, tt.args.store)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSearch_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.SearchInput
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.Region
		wantErr error
		mockFn  func(a args) *Search
	}{
		{
			name:    "ErrorValidation",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{}},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)

				validate.EXPECT().Validate(a.in).Return(assert.AnError)

				return &Search{
					validate: validate,
				}
			},
		},
		{
			name:    "Empty",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{}},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				return &Search{
					validate: validate,
				}
			},
		},
		{
			name:    "ErrorFromProvinces",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "provinces"}},
			want:    nil,
			wantErr: goerror.NewServer("failed to search provinces", assert.AnError),
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Provinces(a.ctx).Return(nil, assert.AnError)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
		{
			name:    "ErrorFromCities",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "cities"}},
			want:    nil,
			wantErr: goerror.NewServer("failed to search cities", assert.AnError),
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Cities(a.ctx, a.in.ParentID).Return(nil, assert.AnError)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
		{
			name:    "ErrorFromDistricts",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "districts"}},
			want:    nil,
			wantErr: goerror.NewServer("failed to search districts", assert.AnError),
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Districts(a.ctx, a.in.ParentID).Return(nil, assert.AnError)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
		{
			name:    "ErrorFromVillages",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "villages"}},
			want:    nil,
			wantErr: goerror.NewServer("failed to search villages", assert.AnError),
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Villages(a.ctx, a.in.ParentID).Return(nil, assert.AnError)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
		{
			name:    "SuccessFromProvinces",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "provinces", IDs: "1,2,3", ParentID: ""}},
			want:    []domain.Region{{ID: "1", Name: "test 1"}},
			wantErr: nil,
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Provinces(a.ctx, testconvertArgs(a.in.IDs)...).
					Return([]domain.Province{{ID: "1", Name: "test 1"}}, nil)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
		{
			name:    "SuccessFromCities",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "cities", IDs: "1,2", ParentID: "2"}},
			want:    []domain.Region{{ID: "1", Name: "test 1"}},
			wantErr: nil,
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Cities(a.ctx, a.in.ParentID, testconvertArgs(a.in.IDs)...).
					Return([]domain.City{{ID: "1", Name: "test 1"}}, nil)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
		{
			name:    "SuccessFromDistricts",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "districts", IDs: "1,2", ParentID: "2"}},
			want:    []domain.Region{{ID: "1", Name: "test 1"}},
			wantErr: nil,
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Districts(a.ctx, a.in.ParentID, testconvertArgs(a.in.IDs)...).
					Return([]domain.District{{ID: "1", Name: "test 1"}}, nil)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
		{
			name:    "SuccessFromVillages",
			args:    args{ctx: context.TODO(), in: domain.SearchInput{By: "villages", IDs: "1,2", ParentID: "2"}},
			want:    []domain.Region{{ID: "1", Name: "test 1"}},
			wantErr: nil,
			mockFn: func(a args) *Search {
				validate := vMock.NewMockValidator(t)
				store := mockz.NewMockSearchStore(t)

				validate.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Villages(a.ctx, a.in.ParentID, testconvertArgs(a.in.IDs)...).
					Return([]domain.Village{{ID: "1", Name: "test 1"}}, nil)

				return &Search{
					store:    store,
					validate: validate,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.mockFn(tt.args)
			got, err := s.Call(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
