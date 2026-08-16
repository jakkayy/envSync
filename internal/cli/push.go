package cli

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/jakkayy/envSync/internal/config"
	"github.com/jakkayy/envSync/pkg/client"
	"github.com/jakkayy/envSync/pkg/crypto"
	"github.com/jakkayy/envSync/pkg/env"
	"github.com/spf13/cobra"
)

var (
	pushEnv     string
	pushMessage string
	passphrase  string
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Encrypt and push local .env file to central server",
	Long:  `Encrypts local environment variables using AES-256-GCM and uploads the encrypted payload to the envSync central API server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadProjectConfig(cfgFile)
		if err != nil {
			return err
		}

		envFile, err := env.ParseFile(".env")
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("local .env file not found")
			}
			return fmt.Errorf("failed to read .env file: %w", err)
		}

		rawEnvContent := []byte(env.Format(envFile))

		if passphrase == "" {
			passphrase = os.Getenv("ENVSYNC_PASSPHRASE")
			if passphrase == "" {
				passphrase = "default-master-key-12345"
			}
		}

		salt := []byte(fmt.Sprintf("%-16s", cfg.ProjectID)[:16])
		secretKey := crypto.DeriveKey(passphrase, salt)

		fmt.Println("🔒 Encrypting environment variables (AES-256-GCM)...")
		encryptedBytes, err := crypto.Encrypt(rawEnvContent, secretKey)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		encodedPayload := base64.StdEncoding.EncodeToString(encryptedBytes)

		apiClient := client.NewAPIClient(cfg.ServerURL, os.Getenv("ENVSYNC_TOKEN"))

		fmt.Printf("⬆️  Pushing %d variables to project '%s' (environment: %s)...\n",
			len(envFile.Map()), cfg.ProjectName, pushEnv)

		version, err := apiClient.Push(cfg.ProjectID, pushEnv, encodedPayload, pushMessage)
		if err != nil {
			return fmt.Errorf("push failed: %w", err)
		}

		fmt.Printf("✅ Successfully updated! (Version: v%d)\n", version)
		return nil
	},
}

func init() {
	pushCmd.Flags().StringVarP(&pushEnv, "env", "e", "dev", "Target environment (dev, staging, prod)")
	pushCmd.Flags().StringVarP(&pushMessage, "message", "m", "", "Commit message describing the changes")
	pushCmd.Flags().StringVarP(&passphrase, "password", "p", "", "Master encryption passphrase")

	RootCmd.AddCommand(pushCmd)
}
