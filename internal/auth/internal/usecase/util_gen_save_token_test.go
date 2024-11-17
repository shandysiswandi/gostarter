package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
)

func Test_generateAndSaveToken(t *testing.T) {
	type args struct {
		ctx     context.Context
		log     logger.Logger
		jwte    jwt.JWT
		secHash hash.Hash
		store   func(context.Context, domain.Token) error
		tid     uint64
		uid     uint64
		email   string
	}
	tests := []struct {
		name    string
		args    args
		want    *generateAndSaveTokenOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateAndSaveToken(tt.args.ctx, tt.args.log, tt.args.jwte, tt.args.secHash, tt.args.store, tt.args.tid, tt.args.uid, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateAndSaveToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateAndSaveToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
