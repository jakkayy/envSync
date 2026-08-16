package ui

import (
	"fmt"

	"github.com/jakkayy/envSync/pkg/env"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

// PrintDiff formatted colored terminal output for DiffResult
func PrintDiff(diff *env.DiffResult, envName string) {
	fmt.Printf("\nComparing local .env with remote (%s%s%s):\n\n", ColorCyan, envName, ColorReset)

	if !diff.HasChanges {
		fmt.Printf("%s✔ Local .env is in sync with remote (%s)!%s\n\n", ColorGreen, envName, ColorReset)
		return
	}

	for _, item := range diff.Items {
		switch item.Type {
		case env.DiffAdded:
			fmt.Printf("  %s+ ADDED:%s    %s%s%s = %s\n", ColorGreen, ColorReset, ColorBold, item.Key, ColorReset, item.LocalValue)
		case env.DiffModified:
			fmt.Printf("  %s~ MODIFIED:%s %s%s%s (%s -> %s)\n", ColorYellow, ColorReset, ColorBold, item.Key, ColorReset, item.LocalValue, item.RemoteValue)
		case env.DiffRemoved:
			fmt.Printf("  %s- REMOVED:%s  %s%s%s (present in remote: %s)\n", ColorRed, ColorReset, ColorBold, item.Key, ColorReset, item.RemoteValue)
		}
	}

	fmt.Printf("\nSummary: %s%d added%s, %s%d modified%s, %s%d removed%s, %d unchanged.\n\n",
		ColorGreen, diff.Added, ColorReset,
		ColorYellow, diff.Modified, ColorReset,
		ColorRed, diff.Removed, ColorReset,
		diff.Unchanged)
}
