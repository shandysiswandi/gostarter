package local

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	type args struct {
		basePath string
	}
	tests := []struct {
		name string
		args args
		want *Storage
	}{
		{
			name: "Success",
			args: args{basePath: "testdata"},
			want: &Storage{basePath: "testdata"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewStorage(tt.args.basePath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStorage_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() *Storage
	}{
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() *Storage {
				return &Storage{basePath: "testdata"}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.mockFn()
			err := s.Close()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestStorage_Upload(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		data []byte
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantErr   error
		mockFn    func() *Storage
		mockClean func()
	}{
		{
			name: "Success",
			args: args{
				ctx:  context.Background(),
				key:  "test_file.txt",
				data: []byte("Hello, world!"),
			},
			wantErr: nil,
			want:    "test_file.txt",
			mockFn: func() *Storage {
				return NewStorage(".")
			},
			mockClean: func() {
				os.Remove("./test_file.txt")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.mockFn()
			got, err := s.Upload(tt.args.ctx, tt.args.key, tt.args.data)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)

			if tt.mockClean != nil {
				tt.mockClean()
			}
		})
	}
}

func TestStorage_Download(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		wantErr   error
		mockFn    func() *Storage
		mockClean func()
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				key: "existing_file.txt",
			},
			want:    []byte("file content"),
			wantErr: nil,
			mockFn: func() *Storage {
				_ = os.WriteFile(filepath.Join(".", "existing_file.txt"), []byte("file content"), 0o644)
				return NewStorage(".")
			},
			mockClean: func() {
				os.Remove("./existing_file.txt")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.mockFn()
			got, err := s.Download(tt.args.ctx, tt.args.key)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)

			if tt.mockClean != nil {
				tt.mockClean()
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func() *Storage
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				key: "remove.txt",
			},
			wantErr: nil,
			mockFn: func() *Storage {
				_ = os.WriteFile(filepath.Join(".", "remove.txt"), []byte("file"), 0o644)
				return NewStorage(".")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.mockFn()
			err := s.Delete(tt.args.ctx, tt.args.key)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestStorage_List(t *testing.T) {
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr error
		mockFn  func() *Storage
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				prefix: "",
			},
			want:    []string{"storage.go", "storage_test.go"},
			wantErr: nil,
			mockFn: func() *Storage {

				return NewStorage(".")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.mockFn()
			got, err := s.List(tt.args.ctx, tt.args.prefix)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
