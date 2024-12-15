package user

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
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
					Database:     nil,
					QueryBuilder: goqu.DialectWrapper{},
					Validator:    nil,
					Router:       framework.NewRouter(),
					Telemetry:    telemetry.NewTelemetry(),
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
