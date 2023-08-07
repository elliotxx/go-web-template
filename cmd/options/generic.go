package options

import (
	"github.com/elliotxx/go-web-template/cmd/options/types"
	"github.com/elliotxx/go-web-template/pkg/server"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var _ types.Options = &GenericOptions{}

// GenericOptions is a generic options struct
type GenericOptions struct {
	ConfigFile  string `json:"configFile" yaml:"configFile"`
	DumpVersion bool   `json:"dumpVersion" yaml:"dumpVersion"`
	DumpEnvs    bool   `json:"dumpEnvs" yaml:"dumpEnvs"`
}

// NewGenericOptions returns a GenericOptions instance with the default values
func NewGenericOptions() *GenericOptions {
	return &GenericOptions{
		ConfigFile:  "",
		DumpVersion: false,
		DumpEnvs:    false,
	}
}

// Validate checks GenericOptions and return a slice of found error(s)
func (o *GenericOptions) Validate() error {
	if o == nil {
		return errors.Errorf("options is nil")
	}

	var err *multierror.Error
	if o.DumpVersion && o.ConfigFile != "" {
		err = multierror.Append(err, errors.Errorf("--dump-version and --config-file/--cors-allowed-origins are mutually exclusive"))
	}

	return err.ErrorOrNil()
}

// ApplyTo apply generic options to the server config
func (o *GenericOptions) ApplyTo(config *server.Config) {}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *GenericOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVarP(&o.ConfigFile, "config-file", "f", "",
		"The path to the configuration file. Valid extension: toml, yaml, yml, json")

	fs.BoolVarP(&o.DumpVersion, "version", "V", o.DumpVersion,
		"Print the version information and exit")

	fs.BoolVarP(&o.DumpEnvs, "envs", "E", o.DumpEnvs,
		"Output all the environment variable names that can be set up, and exit")
}
