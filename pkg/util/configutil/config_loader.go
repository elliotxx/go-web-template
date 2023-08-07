package configutil

import (
	"path/filepath"
	"strings"

	"github.com/elliotxx/go-web-template/third_party/metadecoders"
	"github.com/spf13/afero"
)

var (
	ValidConfigFileExtensions                    = []string{"toml", "yaml", "yml", "json"}
	validConfigFileExtensionsMap map[string]bool = make(map[string]bool)
)

func init() {
	for _, ext := range ValidConfigFileExtensions {
		validConfigFileExtensionsMap[ext] = true
	}
}

// IsValidConfigFilename returns whether filename is one of the supported
// config formats in Hugo.
func IsValidConfigFilename(filename string) bool {
	return validConfigFileExtensionsMap[GetFileExtension(filename)]
}

// FromFile loads the configuration from the given filename.
func FromFile(fs afero.Fs, filename string, data any) error {
	content, err := afero.ReadFile(fs, filename)
	if err != nil {
		return err
	}

	err = FromConfigString(string(content), GetFileExtension(filename), data)
	if err != nil {
		return err
	}

	return nil
}

// FromConfigString creates a config from the given YAML, JSON or TOML config. This is useful in tests.
func FromConfigString(config, configType string, data any) error {
	return metadecoders.Default.UnmarshalTo([]byte(config), metadecoders.FormatFromString(configType), data)
}

// GetFileExtension return the extension of specfied filename
func GetFileExtension(filename string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
}
