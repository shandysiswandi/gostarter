package service

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	vMock "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	type args struct {
		store SetStore
		v     validation.Validator
		t     *telemetry.Telemetry
	}
	tests := []struct {
		name string
		args args
		want *Set
	}{
		{name: "Success", args: args{}, want: &Set{now: time.Now}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSet(tt.args.store, tt.args.v, tt.args.t)
			assert.NotNil(t, tt.want)
			assert.NotNil(t, got)
		})
	}
}

func TestSet_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.SetInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.SetOutput
		wantErr error
		mockFn  func(a args) *Set
	}{
		{
			name: "ErrorValidation",
			args: args{
				ctx: context.TODO(),
				in:  domain.SetInput{URL: "www.google.com", Slug: ""},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *Set {
				mv := new(vMock.MockValidator)
				mtel := telemetry.NewTelemetry()
				ms := new(mockz.MockSetStore)

				mv.EXPECT().Validate(a.in).Return(assert.AnError).Once()

				return &Set{
					store:     ms,
					validator: mv,
					telemetry: mtel,
					now:       func() time.Time { return time.Now() },
				}
			},
		},
		{
			name: "ErrorStoreSet",
			args: args{
				ctx: context.TODO(),
				in:  domain.SetInput{URL: "www.google.com", Slug: ""},
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *Set {
				mv := new(vMock.MockValidator)
				mtel := telemetry.NewTelemetry()
				ms := new(mockz.MockSetStore)

				mv.EXPECT().Validate(a.in).Return(nil).Once()

				input := domain.Short{
					Key:     "MTE2NTEzNzk0OTQ4MzgyMDY0NjQ=",
					URL:     a.in.URL,
					Slug:    a.in.Slug,
					Expired: time.Time{},
				}
				ms.EXPECT().Set(a.ctx, input).Return(assert.AnError).Once()

				return &Set{
					store:     ms,
					validator: mv,
					telemetry: mtel,
					now:       func() time.Time { return time.Time{} },
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.TODO(),
				in:  domain.SetInput{URL: "www.google.com", Slug: ""},
			},
			want:    &domain.SetOutput{Key: "MTE2NTEzNzk0OTQ4MzgyMDY0NjQ="},
			wantErr: nil,
			mockFn: func(a args) *Set {
				mv := new(vMock.MockValidator)
				mtel := telemetry.NewTelemetry()
				ms := new(mockz.MockSetStore)

				mv.EXPECT().Validate(a.in).Return(nil).Once()

				input := domain.Short{
					Key:     "MTE2NTEzNzk0OTQ4MzgyMDY0NjQ=",
					URL:     a.in.URL,
					Slug:    a.in.Slug,
					Expired: time.Time{},
				}
				ms.EXPECT().Set(a.ctx, input).Return(nil).Once()

				return &Set{
					store:     ms,
					validator: mv,
					telemetry: mtel,
					now:       func() time.Time { return time.Time{} },
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Call(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
