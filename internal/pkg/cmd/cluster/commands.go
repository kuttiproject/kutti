package cluster

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/version"

	"github.com/spf13/cobra"
)

var clusterCmd = &cli.Command{
	Cmd: &cobra.Command{
		Use:   "cluster",
		Short: "Manage clusters",
		Long:  "Manage Clusters.",
	},
	Subcommands: []*cli.Command{
		{
			Cmd: &cobra.Command{
				Use:                   "ls",
				Aliases:               []string{"list"},
				Short:                 "List available clusters",
				Run:                   clusterLsCommand,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "show CLUSTERNAME",
				Aliases:               []string{"get", "inspect", "describe"},
				Args:                  cobra.ExactValidArgs(1),
				ValidArgsFunction:     NameValidArgs,
				Short:                 "Show details of a cluster",
				RunE:                  clusterShowCommand,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "select CLUSTERNAME",
				Aliases:               []string{"setdefault", "default"},
				Args:                  cobra.ExactValidArgs(1),
				ValidArgsFunction:     NameValidArgs,
				Short:                 "Select default cluster",
				RunE:                  clusterSelectCommand,
				SilenceErrors:         true,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "unselect",
				Aliases:               []string{"cleardefault"},
				Args:                  cobra.NoArgs,
				Short:                 "Clear default cluster",
				RunE:                  clusterUnselectCommand,
				SilenceErrors:         true,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "rm CLUSTERNAME",
				Aliases:               []string{"remove", "delete", "del"},
				Args:                  cobra.ExactValidArgs(1),
				ValidArgsFunction:     NameValidArgs,
				Short:                 "Remove cluster",
				RunE:                  clusterRmCommand,
				SilenceErrors:         true,
				DisableFlagsInUseLine: true,
			},
			// SetFlagsFunc: func(c *cobra.Command) {
			// 	c.Flags().BoolP("force", "f", false, "forcibly remove cluster")
			// },
		},
		{
			Cmd: &cobra.Command{
				Use:           "create CLUSTERNAME",
				Aliases:       []string{"add"},
				Short:         "Create a new cluster",
				Args:          cobra.ExactArgs(1),
				RunE:          clusterCreateCommand,
				SilenceErrors: true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				version.SetDriverFlag(c)

				c.Flags().StringP("version", "v", "", "K8s version for the cluster")
				c.RegisterFlagCompletionFunc("version", version.NameValidArgs)

				c.Flags().BoolP(
					"unmanaged",
					"u",
					false,
					"create an unmanaged cluster with no nodes",
				)

				c.Flags().BoolP(
					"select",
					"s",
					false,
					"set the new cluster as default",
				)
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "up [CLUSTERNAME]",
				Aliases:               []string{"start"},
				Args:                  cobra.RangeArgs(0, 1),
				ValidArgsFunction:     NameValidArgs,
				Short:                 "Start cluster",
				RunE:                  clusterUpCommand,
				SilenceErrors:         true,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "down [CLUSTERNAME]",
				Aliases:               []string{"stop"},
				Args:                  cobra.RangeArgs(0, 1),
				ValidArgsFunction:     NameValidArgs,
				Short:                 "Stop cluster",
				RunE:                  clusterDownCommand,
				SilenceErrors:         true,
				DisableFlagsInUseLine: true,
			},
		},
	},
}
