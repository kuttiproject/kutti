package version

import (
	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/driver"

	"github.com/spf13/cobra"
)

func CommandTree() *cli.Command {
	return versionCmd
}

// VersionNameValidArgs returns Kubernetes version strings as per Cobra argument
// validation rules.
func VersionNameValidArgs(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	driver, err := getDriver(c)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
	}

	possibilities := driver.VersionNames()
	return cli.StringCompletions(possibilities, toComplete)
}

// SetDriverFlag adds a "--driver" flag to a Cobra command,
// and sets it up for autocompletion with driver names.
func SetDriverFlag(c *cobra.Command) {
	c.Flags().StringP("driver", "d", "", "driver name")

	c.RegisterFlagCompletionFunc(
		"driver",
		driver.DrivernameValidArgs,
	)
}

// GetVersion gets the version and driver from the command line context,
// given the specified Kubernetes version string.
func GetVersion(c *cobra.Command, k8sversion string) (*kuttilib.Version, *kuttilib.Driver, error) {
	driver, err := getDriver(c)
	if err != nil {
		return nil, nil, err
	}

	version, err := driver.GetVersion(k8sversion)
	if err != nil {
		return nil, nil, cli.WrapError(
			2,
			err,
		)
	}

	return version, driver, nil
}
