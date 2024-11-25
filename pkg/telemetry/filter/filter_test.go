package filter

import (
	"reflect"
	"testing"
)

func TestNewFilter(t *testing.T) {
	type args struct {
		opts []OptionFilter
	}
	tests := []struct {
		name string
		args args
		want *Filter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFilter(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Query(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name string
		f    *Filter
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Query(tt.args.rawURL); got != tt.want {
				t.Errorf("Filter.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Body(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name string
		f    *Filter
		args args
		want map[string]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Body(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter.Body() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Header(t *testing.T) {
	type args struct {
		hh map[string][]string
	}
	tests := []struct {
		name string
		f    *Filter
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Header(tt.args.hh); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}
