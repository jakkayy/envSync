package env

import (
	"fmt"
	"os"
	"strings"
)

// Format serializes EnvFile back to dotenv string format preserving comments
func Format(envFile *EnvFile) string {
	var builder strings.Builder

	for i, item := range envFile.Items {
		switch item.Type {
		case ItemComment:
			builder.WriteString(item.Comment)
		case ItemEmptyLine:
			// empty line
		case ItemEntry:
			val := item.Value
			if strings.Contains(val, " ") || strings.Contains(val, "\n") {
				val = fmt.Sprintf(`"%s"`, val)
			}
			builder.WriteString(fmt.Sprintf("%s=%s", item.Key, val))
		}

		if i < len(envFile.Items)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// SaveFile writes formatted EnvFile content to disk
func SaveFile(filepath string, envFile *EnvFile) error {
	content := Format(envFile)
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(filepath, []byte(content), 0644)
}
