package driver

import (
	"os"

	"github.com/kuttiproject/kuttilog"

	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

func driverLsCommand(c *cobra.Command, args []string) {
	defaultdriver, _ := cli.Default("driver")

	var driverlsFormatter = cli.NewTableRenderer(
		"driverls",
		[]*cli.TableColumn{
			{Name: "Name", Width: 10, DefaultCheck: true},
			{Name: "Description", Width: 35},
			{Name: "Status", Width: 10},
		},
		defaultdriver,
	)

	driverlsFormatter.Render(os.Stdout, kuttilib.Drivers())
}

func driverShowCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	drivername := args[0]
	driver, ok := kuttilib.GetDriver(drivername)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"driver '%s' not found",
			drivername,
		)
	}

	renderer := cli.NewJSONRenderer(2)
	renderer.Render(os.Stdout, driver)

	return nil
}

func driverUpdateCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	drivername := args[0]
	driver, ok := kuttilib.GetDriver(drivername)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"driver '%s' not found",
			drivername,
		)
	}

	kuttilog.Println(kuttilog.Minimal, "Updating driver versions...")
	err := driver.UpdateVersionList()
	if err != nil {
		return err
	}

	kuttilog.Println(kuttilog.Minimal, "Driver versions updated.")
	return nil
}

func driverSelectCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	drivername := args[0]
	_, ok := kuttilib.GetDriver(drivername)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"driver '%s' not found",
			drivername,
		)
	}

	return cli.SetDefault("driver", drivername)
}

func driverUnselectCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true
	return cli.RemoveDefault("driver")
}
