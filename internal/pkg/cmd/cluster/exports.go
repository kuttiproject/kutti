package cluster

import (
	"github.com/kuttiproject/kuttilib"
	"github.com/kuttiproject/kuttilog"

	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

// CommandTree returns the top level cluster command
func CommandTree() *cli.Command {
	return clusterCmd
}

// NameValidArgs returns cluster names as per Cobra argument
// validation rules.
func NameValidArgs(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	possibilities := kuttilib.ClusterNames()
	return cli.StringCompletions(possibilities, toComplete)
}

// StartNode starts a node.
func StartNode(cluster *kuttilib.Cluster, nodename string, force bool) error {
	node, ok := cluster.GetNode(nodename)
	if !ok {
		return cli.WrapErrorMessagef(
			1,
			"node '%v' not found",
			nodename,
		)
	}

	nodestatus := node.Status()

	if nodestatus == kuttilib.NodeStatusRunning && (!force) {
		return cli.WrapErrorMessagef(
			1,
			"node '%v' already started",
			nodename,
		)
	}

	if nodestatus == kuttilib.NodeStatusError ||
		nodestatus == kuttilib.NodeStatusUnknown {

		return cli.WrapErrorMessagef(
			1,
			"cannot start node '%v': status unknown",
			nodename,
		)
	}

	kuttilog.Printf(kuttilog.Info, "Starting node %v...", nodename)
	var err error
	if force {
		err = node.ForceStart()
	} else {
		err = node.Start()
	}
	if err != nil {
		kuttilog.Printf(
			kuttilog.Info,
			"Node '%v' could not be started: %v",
			nodename,
			err,
		)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(kuttilog.Info, "Node '%s' started.", nodename)
	} else {
		kuttilog.Println(kuttilog.Quiet, nodename)
	}

	return nil
}

// StopNode stops a node.
func StopNode(cluster *kuttilib.Cluster, nodename string, force bool) error {
	node, ok := cluster.GetNode(nodename)
	if !ok {
		return cli.WrapErrorMessagef(
			1,
			"node '%v' not found",
			nodename,
		)
	}

	nodestatus := node.Status()

	if nodestatus == kuttilib.NodeStatusStopped && (!force) {
		return cli.WrapErrorMessagef(
			1,
			"node '%v' already stopped",
			nodename,
		)
	}

	if nodestatus == kuttilib.NodeStatusError ||
		nodestatus == kuttilib.NodeStatusUnknown {

		return cli.WrapErrorMessagef(
			1,
			"cannot stop node '%v': status unknown",
			nodename,
		)
	}

	var err error
	kuttilog.Printf(kuttilog.Info, "Stopping node %v...", nodename)
	if force {
		err = node.ForceStop()
	} else {
		err = node.Stop()
	}
	if err != nil {
		return cli.WrapErrorMessagef(
			1,
			"Node '%v' could not be stopped: %v",
			nodename,
			err,
		)
	}

	if kuttilog.V(kuttilog.Info) {
		kuttilog.Printf(kuttilog.Info, "Node '%s' stopped.", nodename)
	} else {
		kuttilog.Println(kuttilog.Quiet, nodename)
	}

	return nil
}
