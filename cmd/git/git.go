package git

import (
	"github.com/spf13/cobra"
)

var GitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git related tooling",
	Long: "", 
	Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
	},
}

func init() {
    GitCmd.AddCommand(auditDirCmd)
    GitCmd.AddCommand(conventionalCommitCmd)
}
