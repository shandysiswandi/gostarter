package gql

import (
	"reflect"
	"testing"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
)

func TestRegisterGQLEndpoint(t *testing.T) {
	type args struct {
		router *httpserver.Router
		h      *Endpoint
		cfg    config.Config
		jwte   jwt.JWT
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterGQLEndpoint(tt.args.router, tt.args.h, tt.args.cfg, tt.args.jwte)
		})
	}
}

func TestEndpoint_Mutation(t *testing.T) {
	tests := []struct {
		name string
		e    *Endpoint
		want ql.MutationResolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Mutation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Mutation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Query(t *testing.T) {
	tests := []struct {
		name string
		e    *Endpoint
		want ql.QueryResolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
