package cli

import "github.com/spf13/cobra"

// Command represents a kutti CLI command.
// It wraps a Cobra command object.
type Command struct {
	Cmd          *cobra.Command
	SetFlagsFunc func(*cobra.Command)
	Subcommands  []*Command
}

func (kc *Command) Process(parent *cobra.Command) {
	if kc.Cmd == nil {
		panic("No command to process.")
	}

	if kc.SetFlagsFunc != nil {
		kc.SetFlagsFunc(kc.Cmd)
	}

	if kc.Subcommands != nil {
		for _, sc := range kc.Subcommands {
			sc.Process(kc.Cmd)
		}
	}

	if parent != nil {
		parent.AddCommand(kc.Cmd)
	}
}
