package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jakkayy/envSync/internal/config"
	"github.com/spf13/cobra"
)

var (
	projectName string
	serverURL   string
	defaultEnv  string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize envSync configuration in current project",
	Long:  `Initialize envSync in current project directory by creating a .envsync.json configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.ConfigExists(cfgFile) {
			return fmt.Errorf("project already initialized (.envsync.json already exists)")
		}

		if projectName == "" {
			dir, err := os.Getwd()
			if err == nil {
				projectName = filepath.Base(dir)
			}
			if projectName == "" || projectName == "." {
				projectName = "my-project"
			}
		}

		projectID := "proj_" + uuid.New().String()[:8]

		cfg := &config.ProjectConfig{
			ProjectID:   projectID,
			ProjectName: projectName,
			ServerURL:   serverURL,
			DefaultEnv:  defaultEnv,
		}

		if err := config.SaveProjectConfig(cfg, cfgFile); err != nil {
			return fmt.Errorf("failed to initialize project: %w", err)
		}

		fmt.Printf("✔ Project '%s' (ID: %s) initialized successfully.\n", cfg.ProjectName, cfg.ProjectID)
		fmt.Printf("✔ Config file %s created.\n", config.GetConfigPath(cfgFile))
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "Name of the project")
	initCmd.Flags().StringVar(&serverURL, "server", config.DefaultServerURL, "Central API Server URL")
	initCmd.Flags().StringVarP(&defaultEnv, "env", "e", "dev", "Default environment (dev, staging, prod)")

	RootCmd.AddCommand(initCmd)
}
