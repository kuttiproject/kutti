package cmd

import (
	"fmt"
	"os"

	"github.com/kuttiproject/kutti/internal/pkg/cli"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
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
