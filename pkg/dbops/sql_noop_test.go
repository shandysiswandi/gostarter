package dbops

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetContextNoopTx(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want context.Context
	}{
		{
			name: "Success",
			ctx:  context.Background(),
			want: SetContextNoopTx(context.Background()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SetContextNoopTx(tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewNoopDB(t *testing.T) {
	tests := []struct {
		name string
		want *NoopDB
	}{
		{
			name: "Success",
			want: &NoopDB{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewNoopDB()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoopDB_BeginTx(t *testing.T) {
	type args struct {
		ctx context.Context
		opt *sql.TxOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.Tx
		wantErr error
	}{
		{
			name:    "Success",
			args:    args{},
			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := &NoopDB{}
			got, err := n.BeginTx(tt.args.ctx, tt.args.opt)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
