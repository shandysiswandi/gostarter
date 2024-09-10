package domain

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTodoStatus(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want TodoStatus
	}{
		{name: "Unknown", s: "UNKNOWN", want: TodoStatusUnknown},
		{name: "Initiate", s: "INITIATE", want: TodoStatusInitiate},
		{name: "InProgres", s: "IN_PROGRESS", want: TodoStatusInProgress},
		{name: "Drop", s: "DROP", want: TodoStatusDrop},
		{name: "Done", s: "DONE", want: TodoStatusDone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseTodoStatus(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTodoStatus_String(t *testing.T) {
	tests := []struct {
		name string
		ts   TodoStatus
		want string
	}{
		{name: "Unknown", ts: TodoStatusUnknown, want: "UNKNOWN"},
		{name: "Initiate", ts: TodoStatusInitiate, want: "INITIATE"},
		{name: "InProgres", ts: TodoStatusInProgress, want: "IN_PROGRESS"},
		{name: "Drop", ts: TodoStatusDrop, want: "DROP"},
		{name: "Done", ts: TodoStatusDone, want: "DONE"},
		{name: "Invalid", ts: 100, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.ts.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTodoStatus_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		ts      *TodoStatus
		args    args
		wantErr error
	}{
		{name: "IsString", ts: new(TodoStatus), args: args{value: "MOB"}, wantErr: nil},
		{name: "IsBytes", ts: new(TodoStatus), args: args{value: []uint8{1, 2}}, wantErr: nil},
		{name: "Error", ts: new(TodoStatus), args: args{value: 1}, wantErr: ErrScanTodoStatus},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.ts.Scan(tt.args.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTodoStatus_Value(t *testing.T) {
	tests := []struct {
		name    string
		ts      TodoStatus
		want    driver.Value
		wantErr error
	}{
		{name: "Unknown", ts: TodoStatusUnknown, want: TodoStatusUnknown.String(), wantErr: nil},
		{name: "Initiate", ts: TodoStatusInitiate, want: TodoStatusInitiate.String(), wantErr: nil},
		{name: "InProgress", ts: TodoStatusInProgress, want: TodoStatusInProgress.String(), wantErr: nil},
		{name: "Drop", ts: TodoStatusDrop, want: TodoStatusDrop.String(), wantErr: nil},
		{name: "Done", ts: TodoStatusDone, want: TodoStatusDone.String(), wantErr: nil},
		{name: "InvalidData", ts: -1, want: TodoStatusUnknown.String(), wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.ts.Value()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
