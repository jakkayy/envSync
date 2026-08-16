package cli

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/jakkayy/envSync/internal/config"
	"github.com/jakkayy/envSync/pkg/client"
	"github.com/jakkayy/envSync/pkg/crypto"
	"github.com/jakkayy/envSync/pkg/env"
	"github.com/spf13/cobra"
)

var (
	runEnv        string
	runPassphrase string
)

var runCmd = &cobra.Command{
	Use:   "run -- <command> [args...]",
	Short: "Inject decrypted environment variables directly into process memory",
	Long:  `Downloads and decrypts environment variables directly into process RAM and executes child command without leaving plain .env files on disk.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadProjectConfig(cfgFile)
		if err != nil {
			return err
		}

		apiClient := client.NewAPIClient(cfg.ServerURL, os.Getenv("ENVSYNC_TOKEN"))
		base64Payload, _, err := apiClient.Pull(cfg.ProjectID, runEnv, "")
		if err != nil {
			return fmt.Errorf("failed to fetch secrets from server: %w", err)
		}

		encryptedBytes, err := base64.StdEncoding.DecodeString(base64Payload)
		if err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}

		if runPassphrase == "" {
			runPassphrase = os.Getenv("ENVSYNC_PASSPHRASE")
			if runPassphrase == "" {
				runPassphrase = "default-master-key-12345"
			}
		}

		salt := []byte(fmt.Sprintf("%-16s", cfg.ProjectID)[:16])
		secretKey := crypto.DeriveKey(runPassphrase, salt)

		decryptedBytes, err := crypto.Decrypt(encryptedBytes, secretKey)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		parsedEnv, err := env.Parse(bytes.NewReader(decryptedBytes))
		if err != nil {
			return fmt.Errorf("failed to parse decrypted environment variables: %w", err)
		}

		envSlice := os.Environ()
		for k, v := range parsedEnv.Map() {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
		}

		childCmd := exec.Command(args[0], args[1:]...)
		childCmd.Env = envSlice
		childCmd.Stdin = os.Stdin
		childCmd.Stdout = os.Stdout
		childCmd.Stderr = os.Stderr

		fmt.Printf("🚀 Executing command [%s] with injected secrets (environment: %s)...\n", args[0], runEnv)
		return childCmd.Run()
	},
}

func init() {
	runCmd.Flags().StringVarP(&runEnv, "env", "e", "dev", "Environment to inject (dev, staging, prod)")
	runCmd.Flags().StringVarP(&runPassphrase, "password", "p", "", "Master encryption passphrase")

	RootCmd.AddCommand(runCmd)
}
