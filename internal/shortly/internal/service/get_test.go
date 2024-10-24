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

func TestNewGet(t *testing.T) {
	type args struct {
		store GetStore
		v     validation.Validator
		t     *telemetry.Telemetry
	}
	tests := []struct {
		name string
		args args
		want *Get
	}{
		{name: "Success", args: args{}, want: NewGet(nil, nil, nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewGet(tt.args.store, tt.args.v, tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGet_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.GetInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.GetOutput
		wantErr error
		mockFn  func(a args) *Get
	}{
		{
			name:    "ErrorValidation",
			args:    args{ctx: context.TODO(), in: domain.GetInput{}},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *Get {
				mv := new(vMock.MockValidator)
				mtel := telemetry.NewTelemetry()

				mv.EXPECT().Validate(a.in).Return(assert.AnError).Once()

				return &Get{
					store:     nil,
					validator: mv,
					telemetry: mtel, // TODO
				}
			},
		},
		{
			name:    "ErrorStoreGet",
			args:    args{ctx: context.TODO(), in: domain.GetInput{}},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) *Get {
				mv := new(vMock.MockValidator)
				mtel := telemetry.NewTelemetry()
				ms := new(mockz.MockGetStore)

				mv.EXPECT().Validate(a.in).Return(nil).Once()
				ms.EXPECT().Get(a.ctx, a.in.Key).Return(nil, assert.AnError).Once()

				return &Get{
					store:     ms,
					validator: mv,
					telemetry: mtel,
				}
			},
		},
		{
			name:    "NotFound",
			args:    args{ctx: context.TODO(), in: domain.GetInput{}},
			want:    &domain.GetOutput{},
			wantErr: nil,
			mockFn: func(a args) *Get {
				mv := new(vMock.MockValidator)
				mtel := telemetry.NewTelemetry()

				ms := new(mockz.MockGetStore)

				mv.EXPECT().Validate(a.in).Return(nil).Once()
				ms.EXPECT().Get(a.ctx, a.in.Key).Return(nil, nil).Once()

				return &Get{
					store:     ms,
					validator: mv,
					telemetry: mtel,
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), in: domain.GetInput{Key: "key"}},
			want:    &domain.GetOutput{URL: "www.google.com"},
			wantErr: nil,
			mockFn: func(a args) *Get {
				mv := new(vMock.MockValidator)
				mtel := telemetry.NewTelemetry()
				ms := new(mockz.MockGetStore)

				mv.EXPECT().Validate(a.in).Return(nil).Once()
				ms.EXPECT().Get(a.ctx, a.in.Key).
					Return(&domain.Short{
						Key:     a.in.Key,
						URL:     "www.google.com",
						Slug:    "",
						Expired: time.Time{},
					}, nil).Once()

				return &Get{
					store:     ms,
					validator: mv,
					telemetry: mtel,
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
