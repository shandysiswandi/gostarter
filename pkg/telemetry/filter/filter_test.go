package filter

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/stretchr/testify/assert"
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
		{
			name: "Success",
			args: args{opts: []OptionFilter{
				WithHeaders("test"),
				WithQueries("test"),
				WithFields("test"),
			}},
			want: &Filter{
				headers: []string{"test"},
				queries: []string{"test"},
				fields:  []string{"test"},
				json:    codec.NewJSONCodec(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewFilter(tt.args.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilter_Query(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(a args) *Filter
	}{
		{
			name: "Success",
			args: args{rawURL: "https://example.com?test_key=some_value"},
			want: "https://example.com?test_key=***",
			mockFn: func(a args) *Filter {
				return NewFilter(WithQueries("test_key"))
			},
		},
		{
			name: "MultipleQueries",
			args: args{rawURL: "https://example.com?test_key=some_value&another_key=another_value"},
			want: "https://example.com?test_key=***&another_key=another_value",
			mockFn: func(a args) *Filter {
				return NewFilter(WithQueries("test_key"))
			},
		},
		{
			name: "NoMatchingQuery",
			args: args{rawURL: "https://example.com?non_matching_key=some_value"},
			want: "https://example.com?non_matching_key=some_value",
			mockFn: func(a args) *Filter {
				return NewFilter(WithQueries("test_key"))
			},
		},
		{
			name: "CaseInsensitiveQuery",
			args: args{rawURL: "https://example.com?TEST_KEY=some_value"},
			want: "https://example.com?TEST_KEY=***",
			mockFn: func(a args) *Filter {
				return NewFilter(WithQueries("test_key"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.mockFn(tt.args).Query(tt.args.rawURL)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilter_Body(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name   string
		f      *Filter
		args   args
		want   map[string]any
		mockFn func(a args) *Filter
	}{
		{
			name: "Success",
			args: args{body: []byte(`{"test_key":"some_value"}`)},
			want: map[string]any{"test_key": "***"},
			mockFn: func(a args) *Filter {
				return NewFilter(WithFields("test_key"))
			},
		},
		{
			name: "BodyNil",
			args: args{body: nil},
			want: map[string]any{},
			mockFn: func(a args) *Filter {
				return NewFilter(WithFields("test_key"))
			},
		},
		{
			name: "InvalidBodyJSON",
			args: args{body: []byte(`{invalid_json}`)},
			want: map[string]any{},
			mockFn: func(a args) *Filter {
				return NewFilter(WithFields("test_key"))
			},
		},
		{
			name: "MultipleFields",
			args: args{body: []byte(`{"test_key":"some_value","another_key":"another_value"}`)},
			want: map[string]any{"test_key": "***", "another_key": "another_value"},
			mockFn: func(a args) *Filter {
				return NewFilter(WithFields("test_key"))
			},
		},
		{
			name: "NoMatchingField",
			args: args{body: []byte(`{"non_matching_key":"some_value"}`)},
			want: map[string]any{"non_matching_key": "some_value"},
			mockFn: func(a args) *Filter {
				return NewFilter(WithFields("test_key"))
			},
		},
		{
			name: "NestedFields",
			args: args{body: []byte(`{"nested":{"test_key":"some_value"}}`)},
			want: map[string]any{"nested": map[string]any{"test_key": "***"}},
			mockFn: func(a args) *Filter {
				return NewFilter(WithFields("test_key"))
			},
		},
		{
			name: "ArrayFields",
			args: args{body: []byte(`{"array":[{"test_key":"some_value"},{"test_key":"another_value"}]}`)},
			want: map[string]any{"array": []any{map[string]any{"test_key": "***"}, map[string]any{"test_key": "***"}}},
			mockFn: func(a args) *Filter {
				return NewFilter(WithFields("test_key"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.mockFn(tt.args).Body(tt.args.body)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilter_Header(t *testing.T) {
	type args struct {
		hh map[string][]string
	}
	tests := []struct {
		name   string
		args   args
		want   map[string]string
		mockFn func(a args) *Filter
	}{
		{
			name: "Success",
			args: args{hh: map[string][]string{"Test-Key": {"some_value"}}},
			want: map[string]string{"test-key": "***"},
			mockFn: func(a args) *Filter {
				return NewFilter(WithHeaders("test-key"))
			},
		},
		{
			name: "MultipleHeaders",
			args: args{hh: map[string][]string{"Test-Key": {"some_value"}, "Another-Key": {"another_value"}}},
			want: map[string]string{"test-key": "***", "another-key": "another_value"},
			mockFn: func(a args) *Filter {
				return NewFilter(WithHeaders("test-key"))
			},
		},
		{
			name: "NoMatchingHeader",
			args: args{hh: map[string][]string{"Non-Matching-Key": {"some_value"}}},
			want: map[string]string{"non-matching-key": "some_value"},
			mockFn: func(a args) *Filter {
				return NewFilter(WithHeaders("test-key"))
			},
		},
		{
			name: "CaseInsensitiveHeader",
			args: args{hh: map[string][]string{"TEST-KEY": {"some_value"}}},
			want: map[string]string{"test-key": "***"},
			mockFn: func(a args) *Filter {
				return NewFilter(WithHeaders("test-key"))
			},
		},
		{
			name: "EmptyHeaders",
			args: args{hh: map[string][]string{}},
			want: map[string]string{},
			mockFn: func(a args) *Filter {
				return NewFilter(WithHeaders("test-key"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.mockFn(tt.args).Header(tt.args.hh)
			assert.Equal(t, tt.want, got)
		})
	}
}
