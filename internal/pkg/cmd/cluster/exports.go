package cluster

import (
	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

func CommandTree() *cli.Command {
	return clusterCmd
}

// ClusterNameValidArgs returns cluster names as per Cobra argument
// validation rules.
func ClusterNameValidArgs(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	possibilities := kuttilib.ClusterNames()
	return cli.StringCompletions(possibilities, toComplete)
}
