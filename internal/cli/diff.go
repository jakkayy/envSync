package cli

import (
	"fmt"
	"os"

	"github.com/jakkayy/envSync/internal/config"
	"github.com/jakkayy/envSync/internal/ui"
	"github.com/jakkayy/envSync/pkg/env"
	"github.com/spf13/cobra"
)

var targetEnv string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare local .env file with remote environment configuration",
	Long:  `Displays a side-by-side colorized diff comparing your local .env file against remote environment variables.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.ConfigExists(cfgFile) {
			return fmt.Errorf("project not initialized. Run 'envsync init' first")
		}

		localEnvFile, err := env.ParseFile(".env")
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("local .env file not found")
			}
			return fmt.Errorf("failed to parse local .env file: %w", err)
		}

		remoteDummy := map[string]string{
			"DB_HOST":         "dev-db.internal.company.com",
			"DB_PORT":         "5432",
			"REDIS_TIMEOUT":   "3000",
			"MAX_CONNECTIONS": "20",
		}

		diff := env.CompareEnvs(localEnvFile, remoteDummy)
		ui.PrintDiff(diff, targetEnv)

		return nil
	},
}

func init() {
	diffCmd.Flags().StringVarP(&targetEnv, "env", "e", "dev", "Target environment to compare against")
	RootCmd.AddCommand(diffCmd)
}
