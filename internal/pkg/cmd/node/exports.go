package node

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/cluster"

	"github.com/spf13/cobra"
)

// CommandTree returns the top level node command
func CommandTree() *cli.Command {
	return nodeCmd
}

// SetClusterFlag adds a "--cluster" flag to a Cobra command,
// and sets it up for autocompletion with cluster names.
func SetClusterFlag(c *cobra.Command) {
	c.Flags().StringP("cluster", "c", "", "cluster name")

	c.RegisterFlagCompletionFunc(
		"cluster",
		cluster.NameValidArgs,
	)
}

// NameValidArgs returns node names as per Cobra argument
// validation rules.
func NameValidArgs(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cluster, err := getCluster(c)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
	}

	possibilities := cluster.NodeNames()
	return cli.StringCompletions(possibilities, toComplete)
}
