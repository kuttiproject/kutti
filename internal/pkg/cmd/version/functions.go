package version

import (
	"os"

	"github.com/kuttiproject/kuttilog"

	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

func getDriver(c *cobra.Command) (*kuttilib.Driver, error) {
	drivername, _ := c.Flags().GetString("driver")

	if drivername == "" {
		drivername, _ = cli.Default("driver")
	}

	if drivername == "" {
		return nil, cli.WrapErrorMessage(
			1,
			"no driver specified and default driver not set. Use --driver, or select a default driver using 'kutti driver select'",
		)
	}

	driver, ok := kuttilib.GetDriver(drivername)
	if !ok {
		return nil, cli.WrapErrorMessagef(
			2,
			"driver '%v' not found",
			drivername,
		)
	}

	return driver, nil
}

func versionlsCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	driver, err := getDriver(c)
	if err != nil {
		return err
	}

	defaultversion, _ := cli.Default("version")
	var versionlsFormatter = cli.NewTableRenderer(
		"driverls",
		[]*cli.TableColumn{
			{Name: "K8sVersion", Title: "K8s Version", Width: 15, DefaultCheck: true},
			{Name: "Status", Width: 15},
		},
		defaultversion,
	)

	versionlsFormatter.Render(os.Stdout, driver.Versions())

	return nil
}

func versionShowCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	versionname := args[0]
	version, _, err := GetVersion(c, versionname)
	if err != nil {
		return err
	}

	renderer := cli.NewJSONRenderer(2)
	renderer.Render(os.Stdout, version)

	return nil
}

func versionSelectCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	versionname := args[0]
	version, driver, err := GetVersion(c, versionname)
	if err != nil {
		return err
	}

	err = cli.SetDefault("driver", driver.Name())
	if err != nil {
		return err
	}
	kuttilog.Printf(kuttilog.Info, "Default driver set to %v.\n", driver.Name())

	err = cli.SetDefault("version", version.K8sVersion())
	if err != nil {
		return err
	}
	kuttilog.Printf(kuttilog.Info, "Default version set to %v.\n", version.K8sVersion())

	return nil
}

func versionUnselectCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	return cli.RemoveDefault("version")
}

func versionPullCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	versionname := args[0]

	version, _, err := GetVersion(c, versionname)
	if err != nil {
		return err
	}

	filename, err := c.Flags().GetString("fromfile")
	if err != nil || filename == "" {
		kuttilog.Printf(kuttilog.Minimal, "Downloading image for Kubernetes version %s...", versionname)
		err = version.Fetch()
		if err != nil {
			return cli.WrapErrorMessagef(
				1,
				"could not download image for Kubernetes version %s: %v",
				versionname,
				err,
			)
		}

		if kuttilog.V(kuttilog.Minimal) {
			kuttilog.Printf(kuttilog.Minimal, "Downloaded image for Kubernetes version %s.", versionname)
		} else {
			kuttilog.Println(kuttilog.Quiet, versionname)
		}

		return nil
	}

	kuttilog.Printf(kuttilog.Info, "Importing image for version %v...", versionname)
	err = version.FromFile(filename)
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not import image: %v",
			err,
		)
	}

	if kuttilog.V(kuttilog.Minimal) {
		kuttilog.Printf(kuttilog.Minimal, "Image for version %v imported.", version.K8sVersion())
	} else {
		kuttilog.Println(kuttilog.Quiet, version.K8sVersion())
	}
	return nil
}

func versionRmCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	versionname := args[0]
	version, _, err := GetVersion(c, versionname)
	if err != nil {
		return err
	}

	kuttilog.Printf(kuttilog.Info, "Removing image for Kubernetes version '%v'...\n")
	err = version.PurgeLocal()
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not remove image for Kubernetes version '%v': %v",
			versionname,
			err,
		)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(kuttilog.Info, "Removed image for Kubernetes version '%v'.\n", versionname)
	} else {
		kuttilog.Println(kuttilog.Minimal, versionname)
	}

	return nil
}
