package setting

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

var configcommand = &cli.Command{
	Cmd: &cobra.Command{
		Use:   "setting",
		Short: "Manage configuration settings",
		Long:  `Manage configuration settings.`,
	},
	SetFlagsFunc: nil,
	Subcommands: []*cli.Command{
		{
			Cmd: &cobra.Command{
				Use:                   "ls",
				Aliases:               []string{"list"},
				Args:                  cobra.NoArgs,
				Short:                 "Shows all configuration settings",
				Long:                  `Shows all configuration settings.`,
				Run:                   configlsCommand,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:     "get SETTINGNAME",
				Aliases: []string{"show"},
				Args:    cobra.ExactArgs(1),
				Short:   "Gets a configuration setting",
				Long: `Gets a configuration setting. If the specified setting does not exist, 
outputs an empty string and exits with error code 2.`,
				RunE: configGetCommand,
			},
			SetFlagsFunc: func(c *cobra.Command) {
				c.Flags().BoolP(
					"no-error", "n", false,
					"do not return exit code 2 if setting does not exist.",
				)
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "set SETTINGNAME VALUE",
				Args:                  cobra.ExactArgs(2),
				Short:                 "Sets a configuration setting value",
				Long:                  `Sets a configuration setting value.`,
				RunE:                  configSetCommand,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:     "rm SETTINGNAME",
				Aliases: []string{"remove", "delete", "del"},
				Args:    cobra.ExactArgs(1),
				Short:   "Removes a configuration setting",
				Long:    `Removes a configuration setting.`,
				RunE:    comfigRmCommand,
			},
		},
	},
}
