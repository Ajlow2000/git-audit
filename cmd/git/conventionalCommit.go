package git

import (
	"github.com/Ajlow2000/toolbox/internal/git"
	"github.com/spf13/cobra"
)


var conventionalCommitCmd = &cobra.Command{
	Use:   "conventional-commit",
	Short: "Generates a conventional commit",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
        git.ConventionalCommit()
	},
}

func init() {
}
