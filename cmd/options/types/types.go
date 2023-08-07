package types

import (
	"github.com/spf13/pflag"
)

type Options interface {
	// Validate checks Options and return a slice of found error(s)
	Validate() error
	// AddFlags adds flags for a specific Option to the specified FlagSet
	AddFlags(fs *pflag.FlagSet)
}

const (
	ProjectName = "app"
	MaskString  = "******"
)
