package driver

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

var drivercommand = &cli.Command{
	Cmd: &cobra.Command{
		Use:   "driver",
		Short: "Manage drivers",
		Long:  "Manage drivers.",
	},
	Subcommands: []*cli.Command{
		{
			Cmd: &cobra.Command{
				Use:     "ls",
				Aliases: []string{"list"},
				Args:    cobra.NoArgs,
				Short:   "List available drivers",
				Long:    "List available drivers.",
				Run:     driverLsCommand,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:                   "show DRIVERNAME",
				Aliases:               []string{"get", "inspect", "describe"},
				Args:                  cobra.ExactValidArgs(1),
				ValidArgsFunction:     DrivernameValidArgs,
				Short:                 "Show details of a driver",
				RunE:                  driverShowCommand,
				DisableFlagsInUseLine: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:               "update DRIVERNAME",
				Aliases:           []string{"updateimages"},
				Args:              cobra.ExactValidArgs(1),
				ValidArgsFunction: DrivernameValidArgs,
				Short:             "Update image list for this driver",
				Long:              "Update image list for this driver.",
				RunE:              driverUpdateCommand,

				SilenceErrors: true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:               "select DRIVERNAME",
				Aliases:           []string{"setdefault", "default"},
				Args:              cobra.ExactValidArgs(1),
				ValidArgsFunction: DrivernameValidArgs,
				Short:             "Select default driver",
				RunE:              driverSelectCommand,
				SilenceErrors:     true,
			},
		},
		{
			Cmd: &cobra.Command{
				Use:           "unselect",
				Aliases:       []string{"cleardefault"},
				Args:          cobra.NoArgs,
				Short:         "Clear default driver",
				RunE:          driverUnselectCommand,
				SilenceErrors: true,
			},
		},
	},
}
