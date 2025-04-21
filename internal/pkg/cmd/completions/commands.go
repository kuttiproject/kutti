package completions

import (
	"github.com/kuttiproject/kutti/internal/pkg/cli"

	"github.com/spf13/cobra"
)

var completionCmd = &cli.Command{
	Cmd: &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

  $ source <(kutti completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ kutti completion bash > /etc/bash_completion.d/kutti
  # macOS:
  $ kutti completion bash > /usr/local/etc/bash_completion.d/kutti

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ kutti completion zsh > "${fpath[1]}/_yourprogram"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ kutti completion fish | source

  # To load completions for each session, execute once:
  $ kutti completion fish > ~/.config/fish/completions/kutti.fish

PowerShell:

  PS> kutti completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> kutti completion powershell > kutti.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:                   completionsCommand,
	},
}
