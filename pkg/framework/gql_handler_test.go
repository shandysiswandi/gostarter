package framework

import (
	"reflect"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

func TestWithIntrospection(t *testing.T) {
	tests := []struct {
		name string
		want GQlOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithIntrospection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIntrospection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlerGQL(t *testing.T) {
	type args struct {
		es   graphql.ExecutableSchema
		opts []GQlOption
	}
	tests := []struct {
		name string
		args args
		want *handler.Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandlerGQL(tt.args.es, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandlerGQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
