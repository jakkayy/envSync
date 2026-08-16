package exporter

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jakkayy/envSync/pkg/env"
)

// ExportK8sSecret converts an EnvFile into a Kubernetes Secret YAML manifest
func ExportK8sSecret(secretName string, namespace string, envFile *env.EnvFile) string {
	if secretName == "" {
		secretName = "envsync-secret"
	}
	if namespace == "" {
		namespace = "default"
	}

	var builder strings.Builder
	builder.WriteString("apiVersion: v1\n")
	builder.WriteString("kind: Secret\n")
	builder.WriteString("metadata:\n")
	builder.WriteString(fmt.Sprintf("  name: %s\n", secretName))
	builder.WriteString(fmt.Sprintf("  namespace: %s\n", namespace))
	builder.WriteString("type: Opaque\n")
	builder.WriteString("data:\n")

	for _, item := range envFile.Items {
		if item.Type == env.ItemEntry {
			encodedVal := base64.StdEncoding.EncodeToString([]byte(item.Value))
			builder.WriteString(fmt.Sprintf("  %s: %s\n", item.Key, encodedVal))
		}
	}

	return builder.String()
}
