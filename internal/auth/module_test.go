package auth

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dep     func() Dependency
		want    *Expose
		wantErr error
	}{
		{
			name: "Success",
			dep: func() Dependency {
				return Dependency{
					Router:     httpserver.New(),
					GRPCServer: grpc.NewServer(),
					Telemetry:  telemetry.NewTelemetry(),
				}
			},
			want:    &Expose{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tt.dep())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
