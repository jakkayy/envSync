package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "envsync",
	Short: "envSync is a secure, developer-first CLI to synchronize and protect environment variables.",
	Long: `envSync is a Platform Engineering & DevSecOps CLI tool designed to prevent 
configuration drift, encrypt environment variables client-side using AES-256-GCM, 
and seamlessly sync configs across development teams.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .envsync.json)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
