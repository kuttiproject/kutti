package driver

import (
	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

// CommandTree returns the top level driver command
func CommandTree() *cli.Command {
	return drivercommand
}

// DrivernameValidArgs returns driver names as per Cobra argument
// validation rules.
func DrivernameValidArgs(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	possibilities := kuttilib.DriverNames()
	return cli.StringCompletions(possibilities, toComplete)
}
