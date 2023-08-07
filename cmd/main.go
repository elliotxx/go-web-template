package main

import (
	"expvar"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"

	"github.com/elliotxx/go-web-template/cmd/options"
	"github.com/elliotxx/go-web-template/pkg/util/cmdutil"
)

// NewAppCommand creates a *cobra.Command object with default parameters
func NewAppCommand() *cobra.Command {
	o := options.NewAppOptions()
	// Publish the app options to the public variables,
	// You can visit http://localhost/debug/vars to view all public variables
	expvar.Publish("appOptions", expvar.Func(func() any {
		return o
	}))

	// Initialize root command
	cmd := &cobra.Command{
		Use:  "app",
		Long: `App is a service that operates and maintains configuration code based on gitops technology`,
		// stop printing usage when the command errors
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	// Add flags of each option to root command
	fs := cmd.Flags()
	namedFlagSets := o.Flags()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	// Group options by flag set
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cliflag.SetUsageAndHelpFunc(cmd, namedFlagSets, cols)

	return cmd
}

func main() {
	cmd := NewAppCommand()

	if err := cmd.Execute(); err != nil {
		cmdutil.CheckErr(err)
	}
}
