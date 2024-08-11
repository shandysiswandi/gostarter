package config

import (
	"path"

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
	// Viper does not directly support binary; this treats it as a string and converts to bytes.
	return []byte(vc.v.GetString(key))
}

// GetArray retrieves the configuration value associated with the given key as a slice of strings.
// If the key does not exist, it returns an empty slice.
func (vc *ViperConfig) GetArray(key string) []string {
	return vc.v.GetStringSlice(key)
}

// GetMap retrieves the configuration value associated with the given key as a map of strings to strings.
// If the key does not exist, it returns an empty map.
func (vc *ViperConfig) GetMap(key string) map[string]string {
	return vc.v.GetStringMapString(key)
}

// Close performs any necessary cleanup. For ViperConfig, no resources need to be explicitly closed,
// so this method does nothing and returns nil.
func (vc *ViperConfig) Close() error {
	// No resources to close for ViperConfig; this is just for interface completeness.
	return nil
}
