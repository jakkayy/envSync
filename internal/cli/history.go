package cli

import (
	"fmt"
	"os"

	"github.com/jakkayy/envSync/internal/config"
	"github.com/jakkayy/envSync/pkg/client"
	"github.com/spf13/cobra"
)

var historyEnv string

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View environment configuration change history and timeline",
	Long:  `Displays a chronological history timeline of environment variable changes for the project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadProjectConfig(cfgFile)
		if err != nil {
			return err
		}

		apiClient := client.NewAPIClient(cfg.ServerURL, os.Getenv("ENVSYNC_TOKEN"))
		historyItems, err := apiClient.GetHistory(cfg.ProjectID, historyEnv)
		if err != nil {
			return fmt.Errorf("failed to fetch history: %w", err)
		}

		fmt.Printf("\n📜 Revision History for '%s' (%s):\n\n", cfg.ProjectName, cfg.ServerURL)
		fmt.Printf("%-6s %-20s %-16s %-30s\n", "REV", "DATE", "USER", "MESSAGE")
		fmt.Println("-------------------------------------------------------------------------")

		if len(historyItems) == 0 {
			fmt.Println("No history versions found.")
			return nil
		}

		for _, item := range historyItems {
			version := fmt.Sprintf("v%.0f", item["version"])
			date := fmt.Sprintf("%v", item["created_at"])
			user := fmt.Sprintf("@%v", item["created_by"])
			msg := fmt.Sprintf("%v", item["message"])
			if msg == "" {
				msg = "(no message)"
			}

			fmt.Printf("%-6s %-20s %-16s %-30s\n", version, date, user, msg)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	historyCmd.Flags().StringVarP(&historyEnv, "env", "e", "", "Filter by environment name (dev, staging, prod)")
	RootCmd.AddCommand(historyCmd)
}
