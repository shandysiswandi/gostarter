package user

import (
	"testing"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/stretchr/testify/assert"
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
					SQLKitDB:  nil,
					Hash:      nil,
					Validator: nil,
					Router:    framework.NewRouter(),
					Telemetry: telemetry.NewTelemetry(),
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
