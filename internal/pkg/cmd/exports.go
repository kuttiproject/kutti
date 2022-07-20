package cmd

import (
	"fmt"
	"os"

	"github.com/kuttiproject/kutti/internal/pkg/cli"
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd.Process(nil)

	if err := rootCmd.Cmd.Execute(); err != nil {
		result, ok := cli.UnwrapError(err)

		if !ok {
			fmt.Fprintf(os.Stderr, "Error: %v.\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Error: %v.\n", result.Error())
		os.Exit(result.Exitcode)
	}
}

// SetVersion sets the semantic version string for the current version of kutti
func SetVersion(version string) {
	rootCmd.Cmd.Version = version
}

// ProcessCobraCommandTree builds the command tree, and then invokes the
// supplied callback function, passing it the root Cobra command. This
// is for internal tool use.
func ProcessCobraCommandTree(callback func(c *cobra.Command) error) error {
	cobra.EnableCommandSorting = false
	rootCmd.Process(nil)

	err := callback(rootCmd.Cmd)
	return err
}
