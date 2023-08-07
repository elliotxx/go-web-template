package options

import (
	"time"

	"github.com/elliotxx/errors"
	"github.com/elliotxx/go-web-template/cmd/options/types"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
)

var _ types.Options = &NetworkOptions{}

// NetworkOptions is a Network options struct
type NetworkOptions struct {
	Port                  int           `json:"port,omitempty" yaml:"port,omitempty"`
	CorsAllowedOriginList []string      `json:"corsAllowedOriginList,omitempty" yaml:"corsAllowedOriginList,omitempty"`
	RequestTimeout        time.Duration `json:"requestTimeout,omitempty" yaml:"requestTimeout,omitempty"`
}

// NewNetworkOptions returns a NetworkOptions instance with the default values
func NewNetworkOptions() *NetworkOptions {
	return &NetworkOptions{
		Port:                  80,
		CorsAllowedOriginList: []string{},
		RequestTimeout:        30 * time.Second,
	}
}

// Validate checks NetworkOptions and return a slice of found error(s)
func (o *NetworkOptions) Validate() error {
	if o == nil {
		return errors.Errorf("options is nil")
	}

	var err *multierror.Error

	if o.RequestTimeout <= 0 {
		err = multierror.Append(err, errors.Errorf("--request-timeout must be greater than 0"))
	}

	if o.Port <= 0 {
		err = multierror.Append(err, errors.Errorf("--port must be greater than 0"))
	}

	return err.ErrorOrNil()
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *NetworkOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringSliceVar(&o.CorsAllowedOriginList, "cors-allowed-origins", o.CorsAllowedOriginList,
		"List of allowed origins for CORS, comma separated")

	fs.DurationVar(&o.RequestTimeout, "request-timeout", o.RequestTimeout,
		"An optional field indicating the duration a handler must keep a request open before timing it out")

	fs.IntVarP(&o.Port, "port", "p", o.Port, "Port")
}
