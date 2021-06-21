package cmd

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/cluster"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/completions"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/driver"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/node"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/setting"
	"github.com/kuttiproject/kutti/internal/pkg/cmd/version"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cli.Command{
	Cmd: &cobra.Command{
		Use:              "kutti",
		Short:            "Manage multi-node kubernetes clusters in a local environment",
		Long:             `Manage multi-node kubernetes clusters in a local environment.`,
		PersistentPreRun: setverbosity,
	},
	SetFlagsFunc: func(c *cobra.Command) {
		c.PersistentFlags().BoolP("quiet", "q", false, "produce minimum output")
		c.PersistentFlags().Bool("debug", false, "produce maximum output")
	},
	Subcommands: []*cli.Command{
		completions.CommandTree(),
		setting.CommandTree(),
		driver.CommandTree(),
		version.CommandTree(),
		cluster.CommandTree(),
		node.CommandTree(),
		// Add more commands here
	},
}
