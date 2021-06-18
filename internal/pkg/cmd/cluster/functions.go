package cluster

import (
	"os"

	"github.com/kuttiproject/kuttilog"

	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/version"

	"github.com/spf13/cobra"
)

func getimagename(c *cobra.Command) (string, error) {
	imagename, _ := c.Flags().GetString("image")
	if imagename != "" {
		return imagename, nil
	}

	imagename, ok := cli.Default("image")
	if !ok {
		return "", cli.WrapErrorMessage(
			1,
			"no image specified and default image not set. Use --image, or select a default image using 'kutti image select'",
		)
	}

	return imagename, nil
}

func clusterLsCommand(c *cobra.Command, args []string) {
	defaultcluster, _ := cli.Default("cluster")
	var clusterlsFormatter = cli.NewTableRenderer(
		"driverls",
		[]*cli.TableColumn{
			{Name: "Name", Title: "Name", Width: 15, DefaultCheck: true},
			{Name: "DriverName", Title: "Driver", Width: 15},
			{Name: "K8sVersion", Title: "K8s Version", Width: 15},
			{Name: "Type", Width: 15},
		},
		defaultcluster,
	)

	clusterlsFormatter.Render(os.Stdout, kuttilib.Clusters())
}

func clusterShowCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	clustername := args[0]
	cluster, ok := kuttilib.GetCluster(clustername)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"cluster '%v' not found",
			clustername,
		)
	}

	renderer := cli.NewJSONRenderer(2)
	renderer.Render(os.Stdout, cluster)

	return nil
}

func clusterSelectCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	clustername := args[0]
	_, ok := kuttilib.GetCluster(clustername)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"cluster '%v' not found",
			clustername,
		)
	}

	err := cli.SetDefault("cluster", clustername)
	if err != nil {
		return err
	}

	kuttilog.Printf(
		kuttilog.Info,
		"Cluster '%v' selected as default.",
		clustername,
	)
	return nil
}

func clusterUnselectCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	return cli.RemoveDefault("cluster")
}

func clusterRmCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	clustername := args[0]
	forceflag, _ := c.Flags().GetBool("force")

	kuttilog.Printf(kuttilog.Info, "Removing cluster '%v'...\n", clustername)
	err := kuttilib.DeleteCluster(clustername, forceflag)
	if err != nil {
		return cli.WrapError(1, err)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(kuttilog.Info, "Cluster '%v' removed.\n", clustername)
	} else {
		kuttilog.Println(kuttilog.Minimal, clustername)
	}

	defaultcluster, ok := cli.Default("cluster")
	if ok && (defaultcluster == clustername) {
		cli.RemoveDefault("cluster")
		kuttilog.Println(kuttilog.Info, "Default cluster reset.")
	}

	return nil
}

func clusterCreateCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	clustername := args[0]
	err := kuttilib.ValidateClusterName(clustername)
	if err != nil {
		return err
	}

	imagename, err := getimagename(c)
	if err != nil {
		return err
	}

	image, driver, err := version.GetVersion(c, imagename)
	if err != nil {
		return err
	}
	if image.Status() != kuttilib.VersionStatusDownloaded {
		return cli.WrapErrorMessagef(
			1,
			"local copy of image '%v' has not been downloaded. Cannot create cluster",
			imagename,
		)
	}

	unmanaged, _ := c.Flags().GetBool("unmanaged")
	if !unmanaged {
		return cli.WrapErrorMessage(
			1,
			"managed cluster creation not yet implemented",
		)
	}

	kuttilog.Printf(2, "Creating cluster '%s'...\n", clustername)

	err = kuttilib.NewEmptyCluster(clustername, imagename, driver.Name())
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not create cluster '%v': %v",
			clustername,
			err.Error(),
		)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(kuttilog.Info, "Cluster '%v' created.\n", clustername)
	} else {
		kuttilog.Println(kuttilog.Minimal, clustername)
	}

	setdefault, _ := c.Flags().GetBool("select")
	if setdefault {
		cli.SetDefault("cluster", clustername)
		kuttilog.Printf(kuttilog.Info, "Default cluster set to '%v'.\n", clustername)
	}

	return nil
}

func getclustername(args []string) (string, error) {
	if len(args) == 0 {
		clustername, ok := cli.Default("cluster")
		if !ok {
			return "", cli.WrapErrorMessage(
				1,
				"no cluster specified and default cluster not set. Use --cluster, or select a default cluster using 'kutti cluster select'",
			)
		}

		return clustername, nil
	}

	return args[0], nil
}

func clusterUpCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	clustername, err := getclustername(args)
	if err != nil {
		return err
	}

	cluster, ok := kuttilib.GetCluster(clustername)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"cluster '%v' not found",
			clustername,
		)
	}

	for _, nodename := range cluster.NodeNames() {
		StartNode(cluster, nodename)
	}

	return nil
}

func clusterDownCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	clustername, err := getclustername(args)
	if err != nil {
		return err
	}

	cluster, ok := kuttilib.GetCluster(clustername)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"cluster '%v' not found",
			clustername,
		)
	}

	for _, nodename := range cluster.NodeNames() {
		StopNode(cluster, nodename)
	}

	return nil
}
