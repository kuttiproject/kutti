package node

import (
	"fmt"
	"os"

	"github.com/kuttiproject/kuttilog"

	"github.com/kuttiproject/kuttilib"

	"github.com/kuttiproject/kutti/internal/pkg/cli"
	clustercmd "github.com/kuttiproject/kutti/internal/pkg/cmd/cluster"
	"github.com/kuttiproject/kutti/internal/pkg/sshclient"

	"github.com/spf13/cobra"
)

func getCluster(c *cobra.Command) (*kuttilib.Cluster, error) {
	clustername, _ := c.Flags().GetString("cluster")

	if clustername == "" {
		clustername, _ = cli.Default("cluster")
	}

	if clustername == "" {
		return nil, cli.WrapErrorMessage(
			1,
			"no cluster specified and default cluster not set. Use --cluster, or select a default cluster using 'kutti cluster select'",
		)
	}

	cluster, ok := kuttilib.GetCluster(clustername)
	if !ok {
		return nil, cli.WrapErrorMessagef(
			2,
			"cluster '%v' not found",
			clustername,
		)
	}

	return cluster, nil
}

func nodeLsCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	quiet, _ := c.Root().PersistentFlags().GetBool("quiet")
	if quiet {
		nodenames := cluster.NodeNames()
		for _, nodename := range nodenames {
			fmt.Println(nodename)
		}
		return nil
	}

	var nodelsFormatter = cli.NewTableRenderer(
		"nodels",
		[]*cli.TableColumn{
			{Name: "Name", Width: 15, DefaultCheck: true},
			{Name: "Status", Width: 15},
			{Name: "CreatedAt", Title: "Created", Width: 15, FormatPrefix: "prettytime"},
		},
		"",
	)

	nodelsFormatter.Render(os.Stdout, cluster.Nodes())

	return nil
}

func nodeShowCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	nodename := args[0]

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	node, ok := cluster.GetNode(nodename)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"node '%v' not found",
			nodename,
		)
	}

	renderer := cli.NewJSONRenderer(2)
	renderer.Render(os.Stdout, node)

	return nil
}

func nodeRmCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	nodename := args[0]
	forceflag, _ := c.Flags().GetBool("force")

	// kuttilog.Printf(kuttilog.Info, "Deleting node %s...\n", nodename)
	err = cluster.DeleteNode(nodename, forceflag)
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not delete node '%s': %v",
			nodename,
			err,
		)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(kuttilog.Info, "Node '%s' deleted.", nodename)
	} else {
		kuttilog.Println(kuttilog.Minimal, nodename)
	}

	return nil
}

func nodeCreateCommand(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	// Get cluster to create node in
	cluster, err := getCluster(cmd)
	if err != nil {
		return err
	}

	// Check validity of node name
	nodename := args[0]
	err = cluster.ValidateNodeName(nodename)
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not create node '%v': %v",
			nodename,
			err,
		)
	}

	// Check for sshport for drivers that require it
	driver := cluster.Driver()
	sshport, _ := cmd.Flags().GetInt("sshport")
	if driver.UsesNATNetworking() && sshport == 0 {
		return cli.WrapErrorMessagef(
			1,
			"SSH port forwarding required for nodes in the '%v' cluster",
			cluster.Name(),
		)
	}

	// Check if sshport is occupied
	if sshport != 0 {
		err = cluster.CheckHostport(sshport)
		if err != nil {
			return cli.WrapErrorMessagef(
				1,
				"cannot use host port %v: %v",
				sshport,
				err,
			)
		}
	}

	kuttilog.Printf(kuttilog.Info, "Creating node '%v' on cluster %v...", nodename, cluster.Name())
	newnode, err := cluster.NewUninitializedNode(nodename)
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not create node '%v': %v",
			nodename,
			err,
		)
	}

	// Forward SSH port
	// Belt and suspenders if condition
	if driver.UsesNATNetworking() && sshport != 0 {
		err = newnode.ForwardSSHPort(sshport)
		if err != nil {
			kuttilog.Printf(kuttilog.Quiet, "Warning: Could not forward SSH port: %v.", err)
			// Don't fail node creation
			kuttilog.Printf(kuttilog.Quiet, "Warning: Try manually mapping the SSH port, or delete and re-create this node.")
		}
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(kuttilog.Info, "Node '%s' created.", nodename)
	} else {
		kuttilog.Println(kuttilog.Quiet, nodename)
	}

	return nil
}

func nodeStartCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return cli.WrapErrorMessage(
			1,
			"at least one node name expected",
		)
	}

	if len(args) == 1 {
		return clustercmd.StartNode(cluster, args[0])
	}

	for _, nodename := range args {
		err = clustercmd.StartNode(cluster, nodename)
		if err != nil {
			kuttilog.Printf(kuttilog.Info, "Warning: %v.", err)
		}
	}

	return nil
}

func nodeStopCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return cli.WrapErrorMessage(
			1,
			"at least one node name expected",
		)
	}

	if len(args) == 1 {
		return clustercmd.StopNode(cluster, args[0])
	}

	for _, nodename := range args {
		err = clustercmd.StopNode(cluster, nodename)
		if err != nil {
			kuttilog.Printf(kuttilog.Info, "Warning: %v.", err)
		}
	}

	return nil
}

func nodePublishCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	nodename := args[0]
	node, ok := cluster.GetNode(nodename)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"node '%v' not found",
			nodename,
		)
	}

	nodeport, _ := c.Flags().GetInt("nodeport")
	if !kuttilib.ValidPort(nodeport) {
		return cli.WrapErrorMessage(
			1,
			"please provide a valid nodeport. Valid ports are between 1 and 65535",
		)
	}

	hostport, _ := c.Flags().GetInt("hostport")
	if !kuttilib.ValidPort(hostport) {
		return cli.WrapErrorMessage(
			1,
			"please provide a valid hostport. Valid ports are between 1 and 65535",
		)
	}

	err = cluster.CheckHostport(hostport)
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"cannot forward to host port %v: %v",
			hostport,
			err,
		)
	}

	err = node.ForwardPort(hostport, nodeport)
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not forward node port %v to host port %v: %v",
			nodeport,
			hostport,
			err,
		)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(
			kuttilog.Info,
			"Forwarded node port %v to host port %v.\n",
			nodeport,
			hostport,
		)
	} else {
		kuttilog.Println(kuttilog.Minimal, hostport)
	}

	return nil
}

func nodeUnpublishCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	nodename := args[0]
	node, ok := cluster.GetNode(nodename)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"node '%v' not found",
			nodename,
		)
	}

	nodeport, _ := c.Flags().GetInt("nodeport")
	if !kuttilib.ValidPort(nodeport) {
		return cli.WrapErrorMessage(
			1,
			"please provide a valid nodeport. Valid ports are between 1 and 65535",
		)
	}

	err = node.UnforwardPort(nodeport)
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"could not unforward node port %v: %v",
			nodeport,
			err,
		)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(
			kuttilog.Info,
			"Node port %v unforwarded.\n",
			nodeport,
		)
	} else {
		kuttilog.Println(kuttilog.Minimal, nodeport)
	}

	return nil
}

func nodeSSHCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	if !cluster.Driver().UsesNATNetworking() {
		return cli.WrapErrorMessage(
			1,
			"the SSH command currently on works on clusters that use NAT networking",
		)
	}

	nodename := args[0]
	node, ok := cluster.GetNode(nodename)
	if !ok {
		return cli.WrapErrorMessagef(
			2,
			"node '%v' not found",
			nodename,
		)
	}

	if node.Status() != "Running" {
		return cli.WrapErrorMessagef(
			1,
			"node '%v' is not running",
			nodename,
		)
	}

	sshport, ok := node.Ports()[22]
	if !ok {
		return cli.WrapErrorMessagef(
			1,
			"the SSH port of node '%s' has not been forwarded",
			nodename,
		)
	}

	username, _ := c.Flags().GetString("username")
	if username == "" {
		username = "user1"
	}

	password, _ := c.Flags().GetString("password")
	if username == "" {
		password = "Pass@word1"
	}

	kuttilog.Printf(kuttilog.Info, "Connecting to node %s...", nodename)
	address := fmt.Sprintf("localhost:%v", sshport)
	client := sshclient.NewWithPassword(username, password)

	client.RunInterativeShell(address)

	return nil
}
