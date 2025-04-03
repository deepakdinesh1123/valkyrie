package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(valkyrie completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ valkyrie completion bash > /etc/bash_completion.d/valkyrie
  # macOS:
  $ valkyrie completion bash > /usr/local/etc/bash_completion.d/valkyrie

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ valkyrie completion zsh > "${fpath[1]}/_valkyrie"

Fish:
  $ valkyrie completion fish | source
  # To load completions for each session, execute once:
  $ valkyrie completion fish > ~/.config/fish/completions/valkyrie.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		}
	},
}
