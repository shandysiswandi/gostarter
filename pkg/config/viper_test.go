package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewViperConfig(t *testing.T) {
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
			got, err := NewViperConfig(tt.pathFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestViperConfig_GetInt(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   int64
		mockFn func(a args) *ViperConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: 101,
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", int64(101))

				return &ViperConfig{v: v}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := tt.mockFn(tt.args)
			got := vc.GetInt(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestViperConfig_GetBool(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   bool
		mockFn func(a args) *ViperConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: true,
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", true)

				return &ViperConfig{v: v}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := tt.mockFn(tt.args)
			got := vc.GetBool(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestViperConfig_GetFloat(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   float64
		mockFn func(a args) *ViperConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: 3.14,
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", float64(3.14))

				return &ViperConfig{v: v}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := tt.mockFn(tt.args)
			got := vc.GetFloat(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestViperConfig_GetString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(a args) *ViperConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: "value",
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", "value")

				return &ViperConfig{v: v}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := tt.mockFn(tt.args)
			got := vc.GetString(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestViperConfig_GetBinary(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   []byte
		mockFn func(a args) *ViperConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: []byte("foo"),
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", []byte("Zm9v"))

				return &ViperConfig{v: v}
			},
		},
		{
			name: "Error",
			args: args{key: "key"},
			want: nil,
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", "==")

				return &ViperConfig{v: v}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := tt.mockFn(tt.args)
			got := vc.GetBinary(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestViperConfig_GetArray(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   []string
		mockFn func(a args) *ViperConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: []string{"one", "two"},
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", "one,two")

				return &ViperConfig{v: v}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := tt.mockFn(tt.args)
			got := vc.GetArray(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestViperConfig_GetMap(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   map[string]string
		mockFn func(a args) *ViperConfig
	}{
		{
			name: "Success",
			args: args{key: "key"},
			want: map[string]string{"key": "value", "key1": "value1"},
			mockFn: func(a args) *ViperConfig {
				v := viper.New()
				v.Set("key", "key:value,key1:value1")

				return &ViperConfig{v: v}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vc := tt.mockFn(tt.args)
			got := vc.GetMap(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestViperConfig_Close(t *testing.T) {
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
			vc := &ViperConfig{}
			err := vc.Close()
			assert.Equal(t, tt.want, err)
		})
	}
}
