package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEnum int

const (
	testEnumUnknown testEnum = iota
	testEnumFailed
	testEnumSuccess
)

func (testEnum) Values() map[Enumerate]string {
	return map[Enumerate]string{
		testEnumUnknown: "unknown",
		testEnumFailed:  "failed",
		testEnumSuccess: "success",
	}
}

type testEnumDefault int

const (
	testEnumDefaultUnknown testEnumDefault = iota
	testEnumDefaultOk
)

func (testEnumDefault) Values() map[Enumerate]string {
	return map[Enumerate]string{
		testEnumDefaultUnknown: "unknown",
		testEnumDefaultOk:      "ok",
	}
}

func (testEnumDefault) Default() string { return "n/a" }

func TestNew(t *testing.T) {
	type args struct {
		e testEnum
	}
	tests := []struct {
		name string
		args args
		want Enum[testEnum]
	}{
		{
			name: "unknown",
			args: args{e: testEnumUnknown},
			want: Enum[testEnum]{enum: testEnumUnknown},
		},
		{
			name: "failed",
			args: args{e: testEnumFailed},
			want: Enum[testEnum]{enum: testEnumFailed},
		},
		{
			name: "success",
			args: args{e: testEnumSuccess},
			want: Enum[testEnum]{enum: testEnumSuccess},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New(tt.args.e)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnum_Enum(t *testing.T) {
	tests := []struct {
		name string
		enum Enum[testEnum]
		want testEnum
	}{
		{
			name: "unknown",
			enum: New(testEnumUnknown),
			want: testEnumUnknown,
		},
		{
			name: "failed",
			enum: New(testEnumFailed),
			want: testEnumFailed,
		},
		{
			name: "success",
			enum: New(testEnumSuccess),
			want: testEnumSuccess,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.enum.Enum()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnum_String(t *testing.T) {
	testEnumTests := []struct {
		name string
		enum Enum[testEnum]
		want string
	}{
		{
			name: "invalidTestEnum",
			enum: New(testEnumUnknown - 1),
			want: "UNKNOWN",
		},
		{
			name: "unknown",
			enum: New(testEnumUnknown),
			want: "unknown",
		},
		{
			name: "failed",
			enum: New(testEnumFailed),
			want: "failed",
		},
		{
			name: "success",
			enum: New(testEnumSuccess),
			want: "success",
		},
	}

	testEnumDefaultTests := []struct {
		name string
		enum Enum[testEnumDefault]
		want string
	}{
		{
			name: "invalidTestEnumDefault",
			enum: New(testEnumDefaultUnknown - 1),
			want: "n/a",
		},
		{
			name: "testEnumDefaultUnknown",
			enum: New(testEnumDefaultUnknown),
			want: "unknown",
		},
		{
			name: "testEnumDefaultOk",
			enum: New(testEnumDefaultOk),
			want: "ok",
		},
	}

	for _, tt := range testEnumTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.enum.String()
			assert.Equal(t, tt.want, got)
		})
	}

	for _, tt := range testEnumDefaultTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.enum.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnum_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		args    args
		enum    Enum[testEnum]
		wantErr error
		want    string
	}{
		{
			name:    "ArrayUint8",
			args:    args{value: []uint8{115, 117, 99, 99, 101, 115, 115}},
			enum:    New(testEnumSuccess), // for mock scan
			wantErr: nil,
			want:    "success",
		},
		{
			name:    "ArrayByte",
			args:    args{value: []uint8{115, 117, 99, 99, 101, 115, 115}},
			enum:    New(testEnumSuccess), // for mock scan
			wantErr: nil,
			want:    "success",
		},
		{
			name:    "String",
			args:    args{value: "failed"},
			enum:    New(testEnumSuccess), // for mock scan
			wantErr: nil,
			want:    "failed",
		},
		{
			name:    "Error",
			args:    args{value: 10},
			enum:    New(testEnumUnknown - 1), // for mock scan
			wantErr: ErrInvalidEnum,
			want:    "UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.enum.Scan(tt.args.value)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, tt.enum.String())
		})
	}
}

func TestEnum_Value(t *testing.T) {
	tests := []struct {
		name string
		enum Enum[testEnum]
		want string
	}{
		{
			name: "invalidTestEnum",
			enum: New(testEnumUnknown - 1),
			want: "UNKNOWN",
		},
		{
			name: "success",
			enum: New(testEnumSuccess),
			want: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.enum.Value()
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want testEnum
	}{
		{
			name: "invalidTestEnum",
			args: args{s: "not_found"},
			want: testEnumUnknown,
		},
		{
			name: "success",
			args: args{s: "success"},
			want: testEnumSuccess,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Parse[testEnum](tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
