package cmd

import (
	"fmt"

	githubapi "github.com/jacobsee/applier-gen/pkg/github_api"
	"github.com/spf13/cobra"
)

// getLatestVersionCmd represents the get-latest-version command
var getLatestVersionCmd = &cobra.Command{
	Use:   "get-latest-version",
	Short: "Display the latest version of OpenShift-Applier",
	Long: `Uses the GitHub.com API to determine the latest released version
of OpenShift-Applier. For information purposes only, and not a prerequisite
for any other command (init performs this action automatically).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(githubapi.GetLatestVersionInfo().TagName)
	},
}

func init() {
	rootCmd.AddCommand(getLatestVersionCmd)
}
