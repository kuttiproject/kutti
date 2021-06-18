package cmd

import (
	"github.com/kuttiproject/kuttilib"

	"github.com/spf13/cobra"
)

func setverbosity(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		kuttilib.SetVerbosityLevel(kuttilib.VerbosityDebug)
	} else {
		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			kuttilib.SetVerbosityLevel(kuttilib.VerbosityQuiet)
		}
	}
}
