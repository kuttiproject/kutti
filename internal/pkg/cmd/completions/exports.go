package completions

import "github.com/kuttiproject/kutti/internal/pkg/cli"

// CommandTree returns the top level completion command
func CommandTree() *cli.Command {
	return completionCmd
}
