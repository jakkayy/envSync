package env

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ParseFile reads and parses a dotenv file from filepath
func ParseFile(filepath string) (*EnvFile, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

// Parse parses dotenv content from a reader preserving line structure
func Parse(r io.Reader) (*EnvFile, error) {
	scanner := bufio.NewScanner(r)
	envFile := &EnvFile{Items: []*EnvItem{}}

	var multilineKey string
	var multilineVal strings.Builder
	inMultiline := false

	for scanner.Scan() {
		line := scanner.Text()

		if inMultiline {
			if strings.HasSuffix(line, `"`) || strings.HasSuffix(line, `'`) {
				multilineVal.WriteString("\n")
				multilineVal.WriteString(strings.TrimSuffix(line, line[len(line)-1:]))
				envFile.Items = append(envFile.Items, &EnvItem{
					Type:    ItemEntry,
					Key:     multilineKey,
					Value:   multilineVal.String(),
					RawLine: line,
				})
				inMultiline = false
				multilineKey = ""
				multilineVal.Reset()
			} else {
				multilineVal.WriteString("\n")
				multilineVal.WriteString(line)
			}
			continue
		}

		trimmed := strings.TrimSpace(line)

		// Empty line
		if trimmed == "" {
			envFile.Items = append(envFile.Items, &EnvItem{
				Type:    ItemEmptyLine,
				RawLine: line,
			})
			continue
		}

		// Comment line
		if strings.HasPrefix(trimmed, "#") {
			envFile.Items = append(envFile.Items, &EnvItem{
				Type:    ItemComment,
				Comment: trimmed,
				RawLine: line,
			})
			continue
		}

		// Key-Value entry
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// Strip export keyword
			if strings.HasPrefix(key, "export ") {
				key = strings.TrimPrefix(key, "export ")
				key = strings.TrimSpace(key)
			}

			if (strings.HasPrefix(val, `"`) && !strings.HasSuffix(val, `"`) && len(val) > 1) ||
				(strings.HasPrefix(val, `'`) && !strings.HasSuffix(val, `'`) && len(val) > 1) {
				inMultiline = true
				multilineKey = key
				multilineVal.WriteString(strings.TrimPrefix(val, val[:1]))
				continue
			}

			val = unquote(val)

			envFile.Items = append(envFile.Items, &EnvItem{
				Type:    ItemEntry,
				Key:     key,
				Value:   val,
				RawLine: line,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envFile, nil
}

func unquote(val string) string {
	if (strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`)) ||
		(strings.HasPrefix(val, `'`) && strings.HasSuffix(val, `'`)) {
		if len(val) >= 2 {
			return val[1 : len(val)-1]
		}
	}
	return val
}
