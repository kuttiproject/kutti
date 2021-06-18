package driver

import (
	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

func CommandTree() *cli.Command {
	return drivercommand
}

// DriverNameValidArgs returns driver names as per Cobra argument
// validation rules.
func DrivernameValidArgs(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	possibilities := kuttilib.DriverNames()
	return cli.StringCompletions(possibilities, toComplete)
}
