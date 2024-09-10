package job

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleJob_Start(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		e       *ExampleJob
	}{
		{name: "Success", wantErr: nil, e: &ExampleJob{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.e.Start()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestExampleJob_Stop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		e       *ExampleJob
	}{
		{name: "Success", wantErr: nil, e: &ExampleJob{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.e.Stop(tt.args.ctx)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
