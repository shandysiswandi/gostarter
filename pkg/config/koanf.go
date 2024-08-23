// Package config provides an interface for accessing configuration values from
// various sources such as environment variables, configuration files, or remote
// configuration services.
package config

import (
	"encoding/base64"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// KoanfConfig is an implementation of the Config interface using the Koanf library.
type KoanfConfig struct {
	k *koanf.Koanf
}

// NewKoanfConfig creates a new KoanfConfig instance, initializing Koanf with file and environment variable providers.
func NewKoanfConfig(pathFile string) (*KoanfConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(pathFile), yaml.Parser()); err != nil {
		return nil, err
	}

	return &KoanfConfig{k: k}, nil
}

// GetInt retrieves the configuration value associated with the given key as an int64.
func (kc *KoanfConfig) GetInt(key string) int64 {
	return kc.k.Int64(key)
}

// GetBool retrieves the configuration value associated with the given key as a bool.
func (kc *KoanfConfig) GetBool(key string) bool {
	return kc.k.Bool(key)
}

// GetFloat retrieves the configuration value associated with the given key as a float64.
func (kc *KoanfConfig) GetFloat(key string) float64 {
	return kc.k.Float64(key)
}

// GetString retrieves the configuration value associated with the given key as a string.
func (kc *KoanfConfig) GetString(key string) string {
	return kc.k.String(key)
}

// GetBinary retrieves the configuration value associated with the given key as a byte slice.
func (kc *KoanfConfig) GetBinary(key string) []byte {
	data, err := base64.StdEncoding.DecodeString(kc.k.String(key))
	if err != nil {
		return nil
	}

	return data
}

// GetArray retrieves the configuration value associated with the given key as a slice of strings.
func (kc *KoanfConfig) GetArray(key string) []string {
	return strings.Split(kc.k.String(key), ",")
}

// GetMap retrieves the configuration value associated with the given key as a map of strings to strings.
func (kc *KoanfConfig) GetMap(key string) map[string]string {
	pairs := strings.Split(kc.k.String(key), ",")
	m := make(map[string]string)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}

	return m
}

// Close is implemented to satisfy the io.Closer interface but doesn't do anything in this case.
func (kc *KoanfConfig) Close() error {
	// No resources to close in this simplified example.
	return nil
}
