package cli

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/jakkayy/envSync/internal/config"
	"github.com/jakkayy/envSync/internal/ui"
	"github.com/jakkayy/envSync/pkg/client"
	"github.com/jakkayy/envSync/pkg/crypto"
	"github.com/jakkayy/envSync/pkg/env"
	"github.com/spf13/cobra"
)

var (
	pullEnv        string
	pullPassphrase string
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Fetch and decrypt latest environment configuration from central server",
	Long:  `Downloads the latest encrypted environment payload from the central server, decrypts it using AES-256-GCM, and updates your local .env file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadProjectConfig(cfgFile)
		if err != nil {
			return err
		}

		fmt.Printf("⬇️  Fetching latest config for '%s' (environment: %s)...\n", cfg.ProjectName, pullEnv)

		apiClient := client.NewAPIClient(cfg.ServerURL, os.Getenv("ENVSYNC_TOKEN"))
		base64Payload, version, err := apiClient.Pull(cfg.ProjectID, pullEnv, "")
		if err != nil {
			return fmt.Errorf("pull failed: %w", err)
		}

		encryptedBytes, err := base64.StdEncoding.DecodeString(base64Payload)
		if err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}

		if pullPassphrase == "" {
			pullPassphrase = os.Getenv("ENVSYNC_PASSPHRASE")
			if pullPassphrase == "" {
				pullPassphrase = "default-master-key-12345"
			}
		}

		salt := []byte(fmt.Sprintf("%-16s", cfg.ProjectID)[:16])
		secretKey := crypto.DeriveKey(pullPassphrase, salt)

		fmt.Println("🔓 Decrypting environment variables...")
		decryptedBytes, err := crypto.Decrypt(encryptedBytes, secretKey)
		if err != nil {
			return fmt.Errorf("decryption failed (check passphrase): %w", err)
		}

		var localMap map[string]string
		localEnvFile, err := env.ParseFile(".env")
		if err == nil {
			localMap = localEnvFile.Map()
		} else {
			localMap = make(map[string]string)
		}

		remoteParsed, err := env.Parse(bytes.NewReader(decryptedBytes))
		if err == nil {
			diff := env.CompareMaps(localMap, remoteParsed.Map())
			ui.PrintDiff(diff, pullEnv)
		}

		err = os.WriteFile(".env", decryptedBytes, 0644)
		if err != nil {
			return fmt.Errorf("failed to write local .env file: %w", err)
		}

		fmt.Printf("✔ Local .env file updated successfully to version v%d!\n", version)
		return nil
	},
}

func init() {
	pullCmd.Flags().StringVarP(&pullEnv, "env", "e", "dev", "Target environment (dev, staging, prod)")
	pullCmd.Flags().StringVarP(&pullPassphrase, "password", "p", "", "Master encryption passphrase")

	RootCmd.AddCommand(pullCmd)
}
