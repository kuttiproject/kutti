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
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
				Short:             "Show details of node",
				RunE:              nodeShowCommand,
			},
			SetFlagsFunc: SetClusterFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:               "rm NODENAME",
				Aliases:           []string{"remove", "delete", "del"},
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
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
				ValidArgsFunction: NameValidArgs,
				RunE:              nodeStartCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().BoolP("force", "f", false, "forcibly start node (emergency use only)")
			},
		},
		{
			Cmd: &cobra.Command{
				Use:               "stop NODENAME...",
				Short:             "Stop one or more nodes",
				Args:              cobra.OnlyValidArgs,
				ValidArgsFunction: NameValidArgs,
				RunE:              nodeStopCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().BoolP("force", "f", false, "forcibly stop node (emergency use only)")
			},
		},
		// {
		// 	Cmd: &cobra.Command{
		// 		Use:               "recover NODENAME",
		// 		Short:             "Try to recover an unresponsive node",
		// 		Args:              cobra.ExactValidArgs(1),
		// 		ValidArgsFunction: NameValidArgs,
		// 		RunE:              nodeRecoverCommand,
		// 		SilenceErrors:     true,
		// 	},
		// 	SetFlagsFunc: SetClusterFlag,
		// },
		{
			Cmd: &cobra.Command{
				Use:               "publish NODENAME",
				Short:             "Publish a node port to a host port",
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
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
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
				RunE:              nodeUnpublishCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().IntP("nodeport", "n", 0, "port on the node to unmap")
				c.MarkFlagRequired("nodeport")
			},
		},
		{
			Cmd: &cobra.Command{
				Use:               "ssh NODENAME",
				Aliases:           []string{"attach", "shell"},
				Short:             "Open an SSH connection to the node",
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
				RunE:              nodeSSHCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().StringP("username", "u", "user1", "username for SSH connection")
				c.Flags().StringP("password", "p", "Pass@word1", "password for SSH connection")
			},
		},
		{
			Cmd: &cobra.Command{
				Use:     "scp SOURCE TARGET",
				Aliases: []string{"cp", "copy"},
				Short:   "Copy a file to or from the node",
				Long: `
Copy a file to or from the node.
			
Either the source or the target must begin with a nodename followed by a colon.

Examples:
	kutti node scp /some/file/on/host node1:/some/file
    kutti node scp node1:/some/file /some/file/on/host
	
	kutti node scp -r /some/directory/on/host node1:/some/directory
	kutti node scp -r node1:/some/directory /some/directory/on/host
 
`,
				Args:          cobra.ExactArgs(2),
				RunE:          nodeSCPCommand,
				SilenceErrors: true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetClusterFlag(c)

				c.Flags().BoolP("recurse", "r", false, "copy directories, recursively")
				c.Flags().StringP("username", "u", "user1", "username for SSH connection")
				c.Flags().StringP("password", "p", "Pass@word1", "password for SSH connection")
			},
		},
	},
}
