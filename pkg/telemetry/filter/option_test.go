package filter

import (
	"reflect"
	"testing"
)

func TestWithHeaders(t *testing.T) {
	type args struct {
		header []string
	}
	tests := []struct {
		name string
		args args
		want OptionFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHeaders(tt.args.header...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithQueries(t *testing.T) {
	type args struct {
		query []string
	}
	tests := []struct {
		name string
		args args
		want OptionFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithQueries(tt.args.query...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithFields(t *testing.T) {
	type args struct {
		field []string
	}
	tests := []struct {
		name string
		args args
		want OptionFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithFields(tt.args.field...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
