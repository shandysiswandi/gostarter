package auth

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
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
				return Dependency{
					Database:     nil,
					Transaction:  nil,
					QueryBuilder: goqu.DialectWrapper{},
					Telemetry:    telemetry.NewTelemetry(),
					Router:       framework.NewRouter(),
					GRPCServer:   grpc.NewServer(),
					Validator:    nil,
					UIDNumber:    nil,
					Hash:         nil,
					SecHash:      nil,
					JWT:          nil,
					Clock:        nil,
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
