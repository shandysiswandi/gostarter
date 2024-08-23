// Package config provides an interface for accessing configuration values from
// various sources such as environment variables, configuration files, or remote
// configuration services.
package config

import (
	"encoding/base64"
	"path"
	"strings"

	"github.com/spf13/viper"
)

// ViperConfig is an implementation of the Config interface using the Viper library.
// It provides methods to retrieve configuration values from a file or other sources
// supported by Viper.
type ViperConfig struct {
	v *viper.Viper
}

// NewViperConfig creates a new instance of ViperConfig and loads the configuration
// from the specified file path. The function logs a fatal error if the configuration
// file cannot be read.
func NewViperConfig(pathFile string) (*ViperConfig, error) {
	v := viper.New()

	filename := path.Base(pathFile)
	filePath := path.Dir(pathFile)

	configName := path.Base(filename[:len(filename)-len(path.Ext(filename))])

	v.AddConfigPath(filePath)
	v.SetConfigName(configName)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &ViperConfig{v: v}, nil
}

// GetInt retrieves the configuration value associated with the given key as an int64.
// If the key does not exist or the value cannot be converted to an integer, it returns 0.
func (vc *ViperConfig) GetInt(key string) int64 {
	return vc.v.GetInt64(key)
}

// GetBool retrieves the configuration value associated with the given key as a bool.
// If the key does not exist or the value cannot be converted to a boolean, it returns false.
func (vc *ViperConfig) GetBool(key string) bool {
	return vc.v.GetBool(key)
}

// GetFloat retrieves the configuration value associated with the given key as a float64.
// If the key does not exist or the value cannot be converted to a float, it returns 0.0.
func (vc *ViperConfig) GetFloat(key string) float64 {
	return vc.v.GetFloat64(key)
}

// GetString retrieves the configuration value associated with the given key as a string.
// If the key does not exist, it returns an empty string.
func (vc *ViperConfig) GetString(key string) string {
	return vc.v.GetString(key)
}

// GetBinary retrieves the configuration value associated with the given key as a byte slice.
// If the key does not exist or the value cannot be converted to a byte slice, it returns nil.
func (vc *ViperConfig) GetBinary(key string) []byte {
	data, err := base64.StdEncoding.DecodeString(vc.v.GetString(key))
	if err != nil {
		return nil
	}

	return data
}

// GetArray retrieves the configuration value associated with the given key as a slice of strings.
// If the key does not exist, it returns an empty slice.
func (vc *ViperConfig) GetArray(key string) []string {
	return strings.Split(vc.v.GetString(key), ",")
}

// GetMap retrieves the configuration value associated with the given key as a map of strings to strings.
// If the key does not exist, it returns an empty map.
func (vc *ViperConfig) GetMap(key string) map[string]string {
	pairs := strings.Split(vc.v.GetString(key), ",")
	m := make(map[string]string)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}

	return m
}

// Close performs any necessary cleanup. For ViperConfig, no resources need to be explicitly closed,
// so this method does nothing and returns nil.
func (vc *ViperConfig) Close() error {
	// No resources to close for ViperConfig; this is just for interface completeness.
	return nil
}
