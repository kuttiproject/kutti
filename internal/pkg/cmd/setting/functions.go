package setting

import (
	"fmt"
	"os"

	"github.com/kuttiproject/kuttilog"

	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

func configlsCommand(cmd *cobra.Command, args []string) {
	var configlsFormatter = cli.NewMapTableRenderer(
		"configls",
		[]*cli.TableColumn{
			{Name: "Setting", Width: 10},
			{Name: "Value", Width: 25},
		},
	)

	configlsFormatter.Render(os.Stdout, cli.Settings())
}

func configGetCommand(c *cobra.Command, args []string) error {
	setting := args[0]
	noerror, _ := c.Flags().GetBool("no-error")
	result, ok := cli.Setting(setting)
	if !ok && !noerror {
		return cli.WrapErrorMessagef(
			2,
			"setting '%v' does not exist",
			setting,
		)
	}

	kuttilog.Print(kuttilog.Quiet, result)
	return nil
}

func configSetCommand(c *cobra.Command, args []string) error {
	setting := args[0]
	value := args[1]

	err := cli.SetSetting(setting, value)
	if err != nil {
		return cli.WrapError(
			1,
			err,
		)
	}

	if kuttilog.V(kuttilog.Verbose) {
		kuttilog.Printf(kuttilog.Verbose, "Setting %v set to %v.\n", setting, value)
		return nil
	}

	fmt.Println(value)
	return nil
}

func comfigRmCommand(cmd *cobra.Command, args []string) error {
	setting := args[0]
	return cli.RemoveSetting(setting)
}
