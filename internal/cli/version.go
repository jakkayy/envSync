package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is set dynamically during build time via ldflags
	Version = "v0.1.0-dev"
	Commit  = "none"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of envSync",
	Long:  `All software has versions. This is envSync's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("envSync CLI %s (commit: %s, built at: %s)\n", Version, Commit, Date)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
