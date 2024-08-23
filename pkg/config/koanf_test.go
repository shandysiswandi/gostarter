package config

import (
	"os"
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewKoanfConfig(t *testing.T) {
	tests := []struct {
		name     string
		pathFile string
		wantErr  bool
		beforeFn func()
		afterFn  func()
	}{
		{
			name:     "Error",
			pathFile: "invalid_config.yaml",
			wantErr:  true,
			beforeFn: func() {},
			afterFn:  func() {},
		},
		{
			name:     "Success",
			pathFile: "valid_config.yaml",
			wantErr:  false,
			beforeFn: func() {
				os.WriteFile("valid_config.yaml", []byte(``), 0644)
			},
			afterFn: func() {
				os.Remove("valid_config.yaml")
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.beforeFn()
			defer tt.afterFn()
			got, err := NewKoanfConfig(tt.pathFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestKoanfConfig_GetInt(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   int64
		mockFn func(a args) *KoanfConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: 101,
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", 101)
				return &KoanfConfig{k: k}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kc := tt.mockFn(tt.args)
			got := kc.GetInt(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKoanfConfig_GetBool(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   bool
		mockFn func(a args) *KoanfConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: true,
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", true)
				return &KoanfConfig{k: k}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kc := tt.mockFn(tt.args)
			got := kc.GetBool(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKoanfConfig_GetFloat(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   float64
		mockFn func(a args) *KoanfConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: 3.14,
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", 3.14)
				return &KoanfConfig{k: k}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kc := tt.mockFn(tt.args)
			got := kc.GetFloat(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKoanfConfig_GetString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(a args) *KoanfConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: "value",
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", "value")
				return &KoanfConfig{k: k}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kc := tt.mockFn(tt.args)
			got := kc.GetString(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKoanfConfig_GetBinary(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   []byte
		mockFn func(a args) *KoanfConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: []byte("foo"),
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", "Zm9v")
				return &KoanfConfig{k: k}
			},
		},
		{
			name: "Error",
			args: args{key: "key"},
			want: nil,
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", "==")
				return &KoanfConfig{k: k}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kc := tt.mockFn(tt.args)
			got := kc.GetBinary(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKoanfConfig_GetArray(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   []string
		mockFn func(a args) *KoanfConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: []string{"one", "two"},
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", "one,two")
				return &KoanfConfig{k: k}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kc := tt.mockFn(tt.args)
			got := kc.GetArray(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKoanfConfig_GetMap(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   map[string]string
		mockFn func(a args) *KoanfConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: map[string]string{"key": "value", "key1": "value1"},
			mockFn: func(a args) *KoanfConfig {
				k := koanf.New(".")
				k.Set("key", "key:value,key1:value1")
				return &KoanfConfig{k: k}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kc := tt.mockFn(tt.args)
			got := kc.GetMap(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKoanfConfig_Close(t *testing.T) {
	tests := []struct {
		name string
		want error
	}{
		{
			name: "Success",
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := &KoanfConfig{}
			err := vc.Close()
			assert.Equal(t, tt.want, err)
		})
	}
}
