package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/kuttiproject/kutti/internal/pkg/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	outDir := flag.String("o", "", "output directory")
	outType := flag.String("t", "manpages", "output type (manpages, markdown)")
	flag.Parse()

	if *outDir == "" {
		fmt.Fprintln(
			os.Stderr,
			"Error: must specify output directory with -o option.",
		)
		flag.Usage()
		os.Exit(1)
	}

	if *outType == "" {
		fmt.Fprintln(os.Stderr, "Error: must specify output type.")
		flag.Usage()
		os.Exit(1)
	}

	outTypes := map[string]func(*cobra.Command, string) error{
		"manpages": generateManFile,
		"markdown": generateMarkDown,
	}

	genFunc, ok := outTypes[*outType]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: invalid output type.")
		flag.Usage()
		os.Exit(1)
	}

	err := cmd.ProcessCobraCommandTree(func(cmd *cobra.Command) error {
		walkTree(cmd)
		return genFunc(cmd, *outDir)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v.\n", err)
		os.Exit(1)
	}
}

func walkTree(cmd *cobra.Command) {
	for _, subCmd := range cmd.Commands() {
		walkTree(subCmd)
	}

	cmd.DisableAutoGenTag = true
}

func validateParams(cmd *cobra.Command, directory string) error {
	if cmd == nil {
		return errors.New("no command to process")
	}

	fInfo, err := os.Stat(directory)
	if err != nil {
		return errors.New("error creating man pages: " + err.Error())
	}

	if !fInfo.IsDir() {
		return errors.New(directory + " is not a directory")
	}

	return nil
}

func generateManFile(cmd *cobra.Command, directory string) error {
	err := validateParams(cmd, directory)
	if err != nil {
		return err
	}

	header := &doc.GenManHeader{
		Title:   "kutti",
		Section: "1",
		Source:  "kuttiproject",
		Manual:  "Kutti CLI Manual",
	}

	err = doc.GenManTree(cmd, header, directory)
	return err
}

func generateMarkDown(cmd *cobra.Command, directory string) error {
	err := validateParams(cmd, directory)
	if err != nil {
		return err
	}

	err = doc.GenMarkdownTree(cmd, directory)
	return err
}
