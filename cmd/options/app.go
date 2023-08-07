package options

import (
	"encoding/json"
	"fmt"

	"github.com/elliotxx/go-web-template/cmd/options/types"
	"github.com/elliotxx/go-web-template/pkg/server"
	"github.com/elliotxx/go-web-template/pkg/util/cmdutil"
	"github.com/elliotxx/go-web-template/pkg/util/configutil"
	"github.com/elliotxx/go-web-template/pkg/version"
	"github.com/hashicorp/go-multierror"
	"github.com/koding/multiconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"k8s.io/component-base/cli/flag"
)

// EnvironmentLoader satisifies the loader interface.
// It loads the configuration from the environment variables in the form of
// STRUCTNAME_FIELDNAME.
var envLoader = &multiconfig.EnvironmentLoader{
	// CamelCase adds a separator for field names in camelcase form. A
	// fieldname of "AccessKey" would generate a environment name of
	// "STRUCTNAME_ACCESSKEY". If CamelCase is enabled, the environment name
	// will be generated in the form of "STRUCTNAME_ACCESS_KEY"
	CamelCase: true,
}

// AppOptions runs a App server.
type AppOptions struct {
	Generic  *GenericOptions  `json:"generic,omitempty" yaml:"generic,omitempty"`
	Logging  *LoggingOptions  `json:"logging,omitempty" yaml:"logging,omitempty"`
	Network  *NetworkOptions  `json:"network,omitempty" yaml:"network,omitempty"`
	Database *DatabaseOptions `json:"database,omitempty" yaml:"database,omitempty"`
}

// NewAppOptions creates a new AppOptions object with default parameters
func NewAppOptions() *AppOptions {
	return &AppOptions{
		Generic:  NewGenericOptions(),
		Logging:  NewLoggingOptions(),
		Network:  NewNetworkOptions(),
		Database: NewDatabaseOptions(),
	}
}

// Flags returns flags for a specific APIServer by section name
func (o *AppOptions) Flags() (fss flag.NamedFlagSets) {
	// Add the generic flags
	o.Logging.AddFlags(fss.FlagSet("logging"))
	o.Network.AddFlags(fss.FlagSet("network"))
	o.Generic.AddFlags(fss.FlagSet("generic"))
	o.Database.AddFlags(fss.FlagSet("database"))
	return fss
}

func (o *AppOptions) Complete(args []string) error {
	// NOTE: Priority of configuration effectiveness:
	//   Default value < Command Flags < ConfigFile < Environment

	// Load configuration from file
	if o.Generic.ConfigFile != "" {
		if err := o.LoadConfigFromFile(o.Generic.ConfigFile); err != nil {
			return errors.Wrap(err, "failed to load config file")
		}
	}

	// Load configuration from environment
	err := envLoader.Load(o)
	if err != nil {
		return errors.Wrap(err, "failed to load config file")
	}

	return nil
}

// Validate checks AppOptions and return a slice of found error(s)
func (o *AppOptions) Validate() (err error) {
	if o == nil {
		return errors.Errorf("options is nil")
	}

	err = multierror.Append(err, multierror.Flatten(o.Generic.Validate()))
	if !o.Generic.DumpVersion && !o.Generic.DumpEnvs {
		err = multierror.Append(err, multierror.Flatten(o.Logging.Validate()))
		err = multierror.Append(err, multierror.Flatten(o.Network.Validate()))
		// err = multierror.Append(err, multierror.Flatten(o.Database.Validate()))
	}

	return
}

func (o *AppOptions) Config() *server.Config {
	cfg := server.NewConfig()
	o.Generic.ApplyTo(cfg)
	o.Database.ApplyTo(cfg)
	o.Logging.ApplyTo(cfg)
	return cfg
}

func (o *AppOptions) Run() (err error) {
	if o.Generic.DumpVersion {
		fmt.Println(version.YAML())
		return cmdutil.ErrExit
	}

	if o.Generic.DumpEnvs {
		envLoader.PrintEnvs(o)
		return cmdutil.ErrExit
	}

	// Init logrus configuration by options
	if err := o.Logging.InitLogging(types.ProjectName); err != nil {
		return err
	}

	if o.Logging.DumpCurrentConfig {
		logrus.Info("Dumping the currently used server configuration ...")
		output, err := json.MarshalIndent(o, "", "    ")
		if err != nil {
			logrus.Warn(err)
		}
		logrus.Info(string(output))
	}

	// Start app server
	logrus.Info("Start instantiating App Server ...")
	s, err := o.Config().New()
	if err != nil {
		return err
	}
	logrus.Info("Successfully complete instance!")

	logrus.Info("Start executing predecessor tasks ...")
	if err = s.PreRun(); err != nil {
		return err
	}
	logrus.Info("Successfully complete the predecessor task execution!")

	logrus.Info("Starting the server ...")
	if err = s.Run(fmt.Sprintf(":%d", o.Network.Port)); err != nil {
		return err
	}

	return nil
}

func (o *AppOptions) LoadConfigFromFile(configFile string) error {
	// Validate
	if configFile == "" {
		return errors.Errorf("config file is not specified")
	}

	if !configutil.IsValidConfigFilename(configFile) {
		return errors.Errorf("invalid config file: %s", configFile)
	}

	// Load configuration from file
	err := configutil.FromFile(afero.NewOsFs(), configFile, o)
	if err != nil {
		return err
	}

	return nil
}
