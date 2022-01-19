package node

import (
	"fmt"
	"os"
	"regexp"

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

	forceflag, _ := c.Flags().GetBool("force")
	if len(args) == 1 {
		return clustercmd.StartNode(cluster, args[0], forceflag)
	}

	for _, nodename := range args {
		err = clustercmd.StartNode(cluster, nodename, forceflag)
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

	forceflag, _ := c.Flags().GetBool("force")
	if len(args) == 1 {
		return clustercmd.StopNode(cluster, args[0], forceflag)
	}

	for _, nodename := range args {
		err = clustercmd.StopNode(cluster, nodename, forceflag)
		if err != nil {
			kuttilog.Printf(kuttilog.Info, "Warning: %v.", err)
		}
	}

	return nil
}

// Triaged to a future release
// func nodeRecoverCommand(c *cobra.Command, args []string) error {
// 	c.SilenceUsage = true

// 	cluster, err := getCluster(c)
// 	if err != nil {
// 		return err
// 	}

// 	node, ok := cluster.GetNode(args[0])
// 	if !ok {
// 		return cli.WrapErrorMessagef(
// 			2,
// 			"node %v not found",
// 			args[0],
// 		)
// 	}

// 	// First, capture the status
// 	status := node.Status()

// 	// Then try to force start the node
// 	err = node.ForceStart()

// 	// If that works, stop
// 	if err == nil {
// 		// TODO: Message here
// 		return nil
// 	}
// 	// Then try to force stop the node
// 	err = node.ForceStop()

// 	// If that does not work, stop with error
// 	if err != nil {
// 		// TODO: Better error message here
// 		return cli.WrapError(1, err)
// 	}

// 	// If captured status was stopped, stop
// 	if status == kuttilib.NodeStatusStopped {
// 		// TODO: Message here
// 		return nil
// 	}
// 	// Try to normal start the node
// 	err = node.Start()

// 	// If that works, stop
// 	if err == nil {
// 		// TODO: Message here
// 		return nil
// 	}
// 	// Try to force start the node
// 	err = node.ForceStart()

// 	// If that works, stop
// 	if err == nil {
// 		// TODO: Message here
// 		return nil
// 	}

// 	// Stop with error
// 	return cli.WrapErrorMessagef(
// 		1,
// 		"could not recover node. You may have to delete and recreate it. The last error returned was: %v",
// 		err,
// 	)
// }

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

func getNodeSSHPort(cluster *kuttilib.Cluster, nodename string) (int, error) {
	node, ok := cluster.GetNode(nodename)
	if !ok {
		return 0, cli.WrapErrorMessagef(
			2,
			"node '%v' not found",
			nodename,
		)
	}

	if node.Status() != "Running" {
		return 0, cli.WrapErrorMessagef(
			1,
			"node '%v' is not running",
			nodename,
		)
	}

	sshport, ok := node.Ports()[22]
	if !ok {
		return 0, cli.WrapErrorMessagef(
			1,
			"the SSH port of node '%s' has not been forwarded",
			nodename,
		)
	}

	return sshport, nil
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
			"the SSH command currently only works on clusters that use NAT networking",
		)
	}

	nodename := args[0]
	sshport, err := getNodeSSHPort(cluster, nodename)
	if err != nil {
		return err
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

type cparg struct {
	nodename          string
	filepath          string
	hasnodename       bool
	iswindowsfilepath bool
	localfileexists   bool
	localisdirectory  bool
}

func parseCPArg(arg string) (*cparg, error) {
	// Look for a nodename followed by a colon followed by a path,
	// or a drive letter followed by a path.
	// If it is a nodename, the nodename plus the colon will
	// be submatch[1], submatch[2] will be just the nodename, and
	// submatch[5] will be the path.
	// If it is a drive letter, the letter plus the colon will
	// be submatch[1], submatch[2] will be just the letter, and
	// submatch[5] will be the path.
	// The nodename followed by a colon may not appear at all, in
	// which case submatch[1] and [2] will be empty, and
	// submatch[5] will be just the path.
	cpargregex, _ := regexp.Compile("^((([A-Z])|([a-z][a-z0-9]{0,9})):){0,1}([^:]*)$")
	results := cpargregex.FindStringSubmatch(arg)

	// If no match, argument is invalid
	if len(results) < 6 {
		return nil, cli.WrapErrorMessagef(
			1,
			"could not understand '%v'",
			arg,
		)
	}

	result := &cparg{}

	// If match, and second submatch is empty,
	// the argment is just a file path
	if results[2] == "" {
		result.filepath = results[5]
	} else {
		if len(results[2]) == 1 {
			// First submatch is a drive letter plus colon,
			// fifth submatch has a path
			result.filepath = results[1] + results[5]
			result.iswindowsfilepath = true
		} else {

			// Second submatch is node name, fifth
			// submatch is file path.
			result.nodename = results[2]
			result.filepath = results[5]
			result.hasnodename = true
		}
	}

	// Do a standard OS check on the path
	// It may exist on the host
	fi, err := os.Stat(result.filepath)
	if err == nil {
		result.localfileexists = true
		result.localisdirectory = fi.IsDir()
		return result, nil
	}

	// If the error is IsNotExist, and not
	// anything else, file path seems to be
	// valid.
	if os.IsNotExist(err) {
		return result, nil
	}

	// Otherwise, file path is invalid
	return nil, cli.WrapError(
		1,
		err,
	)
}

func nodeSCPCommand(c *cobra.Command, args []string) error {
	c.SilenceUsage = true

	cluster, err := getCluster(c)
	if err != nil {
		return err
	}

	if !cluster.Driver().UsesNATNetworking() {
		return cli.WrapErrorMessage(
			1,
			"the scp command currently only works on clusters that use NAT networking",
		)
	}

	// Parse the arguments
	arg1, err := parseCPArg(args[0])
	if err != nil {
		return err
	}

	arg2, err := parseCPArg(args[1])
	if err != nil {
		return err
	}

	// If neither the first argument, nor the second
	// have a nodename, the user should be using the
	// cp or copy command instead.
	if !(arg1.hasnodename || arg2.hasnodename) {
		return cli.WrapErrorMessage(
			1,
			"must specify at least one node",
		)

	}

	// If both arguments have a nodename, it is an
	// error.
	if arg1.hasnodename && arg2.hasnodename {
		return cli.WrapErrorMessage(
			1,
			"copying between nodes is not supported",
		)
	}

	// If the first (source) argument does not have
	// a nodename, then the file or directory
	// specified must exist on the host.
	if (!arg1.hasnodename) && (!arg1.localfileexists) {
		return cli.WrapErrorMessagef(
			2,
			"'%v': no such file or directory",
			arg1.filepath,
		)
	}

	// If the second (destination) argument has a
	// nodename but not a path, we should assume
	// the current directoy on the node.
	if arg2.hasnodename && arg2.filepath == "" {
		arg2.filepath = "."
	}

	recurseFlag, _ := c.Flags().GetBool("recurse")

	// If the first (source) argument does not have
	// a nodename, and is a directory, the --recurse
	// flag should be specified.
	if (!arg1.hasnodename) &&
		arg1.localisdirectory &&
		(!recurseFlag) {

		return cli.WrapErrorMessagef(
			1,
			"'%v' is a directory. Use the --recurse option.",
			arg1.filepath,
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

	// If the first (source) argument has a nodename, this is a
	// copy from node operation.
	if arg1.hasnodename {
		sshport, err := getNodeSSHPort(cluster, arg1.nodename)
		if err != nil {
			return err
		}

		kuttilog.Printf(kuttilog.Info, "Copying from node %s...", arg1.nodename)
		address := fmt.Sprintf("localhost:%v", sshport)
		client := sshclient.NewWithPassword(username, password)

		err = client.CopyFrom(address, arg1.filepath, arg2.filepath, recurseFlag)

		if err != nil {
			return err
		}
	} else {
		sshport, err := getNodeSSHPort(cluster, arg2.nodename)
		if err != nil {
			return err
		}

		kuttilog.Printf(kuttilog.Info, "Copying to node %s...", arg2.nodename)
		address := fmt.Sprintf("localhost:%v", sshport)
		client := sshclient.NewWithPassword(username, password)

		err = client.CopyTo(address, arg1.filepath, arg2.filepath, recurseFlag)
		if err != nil {
			return err
		}
	}

	return nil
}
