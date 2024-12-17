package todo

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	configMock "github.com/shandysiswandi/gostarter/pkg/config/mocker"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dep     func() Dependency
		wantErr error
	}{
		{
			name: "Success",
			dep: func() Dependency {
				mc := configMock.NewMockConfig(t)
				mc.EXPECT().GetBool("feature.flag.todo.job").Return(true).Once()

				return Dependency{
					Database:     nil,
					QueryBuilder: goqu.DialectWrapper{},
					Messaging:    nil,
					Config:       mc,
					UIDNumber:    nil,
					CodecJSON:    nil,
					Validator:    nil,
					Router:       framework.NewRouter(),
					GQLRouter:    framework.NewRouter(),
					GRPCServer:   grpc.NewServer(),
					Telemetry:    nil,
				}
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tt.dep())
			assert.Equal(t, tt.wantErr, err)
			assert.NotNil(t, got)
		})
	}
}
