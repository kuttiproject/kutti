package setting

import "github.com/kuttiproject/kutti/internal/pkg/cli"

// CommandTree returns the top level setting command
func CommandTree() *cli.Command {
	return configcommand
}
