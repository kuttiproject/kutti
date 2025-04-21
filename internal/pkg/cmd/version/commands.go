package version

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

var versionCmd = &cli.Command{
	Cmd: &cobra.Command{
		Use:   "version",
		Short: "Manage Kubernetes versions",
	},
	Subcommands: []*cli.Command{
		{
			Cmd: &cobra.Command{
				Use:           "ls",
				Aliases:       []string{"list"},
				Args:          cobra.NoArgs,
				Short:         "List versions",
				RunE:          versionlsCommand,
				SilenceErrors: true,
			},
			SetFlagsFunc: SetDriverFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:               "show K8SVERSION",
				Aliases:           []string{"get", "inspect", "describe"},
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
				Short:             "Show details of a version",
				RunE:              versionShowCommand,
			},
			SetFlagsFunc: SetDriverFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:               "select K8SVERSION",
				Aliases:           []string{"setdefault", "default"},
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
				Short:             "Select default version",
				RunE:              versionSelectCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: SetDriverFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:                   "unselect",
				Aliases:               []string{"cleardefault"},
				Args:                  cobra.NoArgs,
				Short:                 "Clear default version",
				RunE:                  versionUnselectCommand,
				SilenceErrors:         true,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "pull [flags] K8SVERSION",
				Aliases:               []string{"fetch", "get"},
				Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction:     NameValidArgs,
				Short:                 "Download version image",
				RunE:                  versionPullCommand,
				SilenceErrors:         true,
				DisableFlagsInUseLine: true,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				SetDriverFlag(c)

				c.Flags().StringP("fromfile", "f", "", "local file path to import version image from")
				c.MarkFlagFilename("fromfile")
			},
		},
		{
			Cmd: &cobra.Command{
				Use:               "rm K8SVERSION",
				Aliases:           []string{"remove", "delete", "del", "purge", "purgelocal"},
				Args:              cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				ValidArgsFunction: NameValidArgs,
				Short:             "Remove version image",
				RunE:              versionRmCommand,
				SilenceErrors:     true,
			},
			SetFlagsFunc: SetDriverFlag,
		},
		{
			Cmd: &cobra.Command{
				Use:           "update",
				Args:          cobra.NoArgs,
				Short:         "Update image list",
				RunE:          versionUpdateCommand,
				SilenceErrors: true,
				SilenceUsage:  true,
			},
			SetFlagsFunc: SetDriverFlag,
		},
	},
}
