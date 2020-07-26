package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/hsnodgrass/pufctl/internal/uitext"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:                   uitext.CompletionUse,
	Short:                 uitext.CompletionShort,
	Long:                  uitext.CompletionLong,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
