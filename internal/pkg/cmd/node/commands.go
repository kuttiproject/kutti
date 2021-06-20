package node

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

var nodeCmd = &cli.Command{
	Cmd: &cobra.Command{
		Use:   "node",
		Short: "Manage nodes",
	},
	Subcommands: []*cli.Command{
		{
			Cmd: &cobra.Command{
				Use:           "ls",
				Aliases:       []string{"list"},
				Args:          cobra.NoArgs,
				Short:         "List nodes",
				RunE:          nodeLsCommand,
				SilenceErrors: true,
			},
			SetFlagsFunc: SetClusterFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:               "show NODENAME",
				Aliases:           []string{"get", "inspect", "describe"},
				Args:              cobra.ExactValidArgs(1),
				ValidArgsFunction: NodeNameValidArgs,
				Short:             "Show details of node",
				RunE:              nodeShowCommand,
			},
			SetFlagsFunc: SetClusterFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:               "rm NODENAME",
				Aliases:           []string{"remove", "delete", "del"},
				Args:              cobra.ExactValidArgs(1),
				ValidArgsFunction: NodeNameValidArgs,
				Short:             "Remove node",
				RunE:              nodeRmCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().BoolP("force", "f", false, "forcibly delete running nodes.")
			},
		},
		{
			Cmd: &cobra.Command{
				Use:           "create NODENAME",
				Aliases:       []string{"add"},
				Short:         "Create a new node",
				Args:          cobra.ExactArgs(1),
				RunE:          nodeCreateCommand,
				SilenceErrors: true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().IntP("sshport", "p", 0, "host port to forward node SSH port")
			},
		},
		{
			Cmd: &cobra.Command{
				Use:               "start NODENAME...",
				Short:             "Start one or more nodes",
				Args:              cobra.OnlyValidArgs,
				ValidArgsFunction: NodeNameValidArgs,
				RunE:              nodeStartCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: SetClusterFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:               "stop NODENAME...",
				Short:             "Stop one or more nodes",
				Args:              cobra.OnlyValidArgs,
				ValidArgsFunction: NodeNameValidArgs,
				RunE:              nodeStopCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: SetClusterFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:               "publish NODENAME",
				Short:             "Publish a node port to a host port",
				Args:              cobra.ExactValidArgs(1),
				ValidArgsFunction: NodeNameValidArgs,
				RunE:              nodePublishCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().IntP("hostport", "p", 0, "port on the host")
				c.Flags().IntP("nodeport", "n", 0, "port on the node")

				c.MarkFlagRequired("hostport")
				c.MarkFlagRequired("nodeport")
			},
		},
		{
			Cmd: &cobra.Command{
				Use:               "unpublish NODENAME",
				Short:             "Un-publish a node port",
				Args:              cobra.ExactValidArgs(1),
				ValidArgsFunction: NodeNameValidArgs,
				RunE:              nodeUnpublishCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().IntP("nodeport", "n", 0, "port on the node to unmap")
				c.MarkFlagRequired("nodeport")
			},
		},
	},
}
