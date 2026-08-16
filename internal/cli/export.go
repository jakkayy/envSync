package cli

import (
	"fmt"
	"os"

	"github.com/jakkayy/envSync/pkg/env"
	"github.com/jakkayy/envSync/pkg/exporter"
	"github.com/spf13/cobra"
)

var (
	exportName       string
	exportNamespace  string
	exportOutputFile string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export environment variables to external formats (e.g. Kubernetes Secret)",
	Long:  `Export your local or synced environment configuration into external ecosystem formats such as Kubernetes Secret manifests.`,
}

var exportK8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Export .env file as a Kubernetes Secret YAML manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile, err := env.ParseFile(".env")
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("local .env file not found")
			}
			return fmt.Errorf("failed to read .env file: %w", err)
		}

		yamlOutput := exporter.ExportK8sSecret(exportName, exportNamespace, envFile)

		if exportOutputFile != "" {
			if err := os.WriteFile(exportOutputFile, []byte(yamlOutput), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Printf("✔ Kubernetes Secret manifest successfully exported to '%s'\n", exportOutputFile)
		} else {
			fmt.Println(yamlOutput)
		}

		return nil
	},
}

func init() {
	exportK8sCmd.Flags().StringVar(&exportName, "name", "app-secret", "Kubernetes Secret metadata name")
	exportK8sCmd.Flags().StringVarP(&exportNamespace, "namespace", "n", "default", "Kubernetes Secret namespace")
	exportK8sCmd.Flags().StringVarP(&exportOutputFile, "output", "o", "", "Output file path (prints to stdout if omitted)")

	exportCmd.AddCommand(exportK8sCmd)
	RootCmd.AddCommand(exportCmd)
}
