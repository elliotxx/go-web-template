package entity

import (
	"fmt"
	"time"
)

// SystemConfig represents the configuration of a system.
type SystemConfig struct {
	// Unique ID of the system
	ID uint `yaml:"id" json:"id"`
	// Tenant or organization that the system belongs to
	Tenant string `yaml:"tenant" json:"tenant"`
	// Environment where the system is deployed (e.g. prod, gray)
	Env Env `yaml:"env" json:"env"`
	// Type or category of the system (e.g. cache, message queue)
	Type string `yaml:"type" json:"type"`
	// Configuration data in JSON or YAML format
	Config string `yaml:"config,omitempty" json:"config,omitempty"`
	// Description or purpose of the system
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Username or ID of the user who created the system
	Creator string `yaml:"creator,omitempty" json:"creator,omitempty"`
	// Username or ID of the user who last modified the system
	Modifier string `yaml:"modifier,omitempty" json:"modifier,omitempty"`
	// Timestamp when the system was created
	CreatedAt time.Time `yaml:"createdAt,omitempty" json:"createdAt,omitempty"`
	// Timestamp when the system was last updated
	UpdatedAt time.Time `yaml:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// Validate checks if the system config is valid.
// It returns an error if the system config is not valid.
func (s *SystemConfig) Validate() error {
	if _, err := ParseEnv(string(s.Env)); err != nil {
		return err
	}

	return nil
}

// Env represents the environment.
type Env string

// These constants represent the possible environment.
const (
	// EnvPre represents the pre environment.
	EnvPre Env = "pre"

	// EnvGray represents the gray environment.
	EnvGray Env = "gray"

	// EnvProd represents the production environment.
	EnvProd Env = "prod"

	// EnvDev represents the DEV environment.
	EnvDev Env = "dev"

	// EnvTest represents the TEST environment.
	EnvTest Env = "test"

	// EnvStable represents the STABLE environment.
	EnvStable Env = "stable"
)

// ParseEnv parses a string into a Env.
// If the string is not a valid Env, it returns an error.
func ParseEnv(str string) (Env, error) {
	switch str {
	case "pre":
		return EnvPre, nil
	case "gray":
		return EnvGray, nil
	case "prod":
		return EnvProd, nil
	case "dev":
		return EnvDev, nil
	case "test":
		return EnvTest, nil
	case "stable":
		return EnvStable, nil
	default:
		return Env(""), fmt.Errorf("invalid environment: %q", str)
	}
}

// MustParseEnv parses a string into a Env.
// If the string is not a valid Env, it panics.
func MustParseEnv(str string) Env {
	env, err := ParseEnv(str)
	if err != nil {
		panic(err)
	}

	return env
}
