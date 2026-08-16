package cli

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/jakkayy/envSync/internal/config"
	"github.com/jakkayy/envSync/internal/ui"
	"github.com/jakkayy/envSync/pkg/client"
	"github.com/jakkayy/envSync/pkg/crypto"
	"github.com/spf13/cobra"
)

var (
	rollbackEnv        string
	targetVersion      string
	rollbackPassphrase string
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback local environment configuration to a previous version",
	Long:  `Fetches a specific historical version from the central server and restores your local .env file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetVersion == "" {
			return fmt.Errorf("--to flag is required (e.g. envsync rollback --to 2)")
		}

		cfg, err := config.LoadProjectConfig(cfgFile)
		if err != nil {
			return err
		}

		fmt.Printf("⏪ Rolling back '%s' (environment: %s) to version v%s...\n", cfg.ProjectName, rollbackEnv, targetVersion)

		apiClient := client.NewAPIClient(cfg.ServerURL, os.Getenv("ENVSYNC_TOKEN"))
		base64Payload, version, err := apiClient.Pull(cfg.ProjectID, rollbackEnv, targetVersion)
		if err != nil {
			return fmt.Errorf("rollback failed: %w", err)
		}

		encryptedBytes, err := base64.StdEncoding.DecodeString(base64Payload)
		if err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}

		if rollbackPassphrase == "" {
			rollbackPassphrase = os.Getenv("ENVSYNC_PASSPHRASE")
			if rollbackPassphrase == "" {
				rollbackPassphrase = "default-master-key-12345"
			}
		}

		salt := []byte(fmt.Sprintf("%-16s", cfg.ProjectID)[:16])
		secretKey := crypto.DeriveKey(rollbackPassphrase, salt)

		decryptedBytes, err := crypto.Decrypt(encryptedBytes, secretKey)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		err = os.WriteFile(".env", decryptedBytes, 0644)
		if err != nil {
			return fmt.Errorf("failed to update local .env file: %w", err)
		}

		fmt.Printf("%s✔ Local .env file successfully rolled back to version v%d!%s\n", ui.ColorGreen, version, ui.ColorReset)
		return nil
	},
}

func init() {
	rollbackCmd.Flags().StringVarP(&rollbackEnv, "env", "e", "dev", "Target environment (dev, staging, prod)")
	rollbackCmd.Flags().StringVar(&targetVersion, "to", "", "Target version number to rollback to (e.g. 1, 2, 3)")
	rollbackCmd.Flags().StringVarP(&rollbackPassphrase, "password", "p", "", "Master encryption passphrase")

	RootCmd.AddCommand(rollbackCmd)
}
