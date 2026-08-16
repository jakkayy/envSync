package security

import (
	"strings"
)

var sensitiveKeyKeywords = []string{
	"PASSWORD", "PASS", "SECRET", "KEY", "TOKEN", "AUTH", "CREDENTIAL", "PRIVATE",
}

// MaskValue masks sensitive values for safe logging
func MaskValue(key string, val string) string {
	if val == "" {
		return ""
	}

	upperKey := strings.ToUpper(key)
	isSensitive := false
	for _, kw := range sensitiveKeyKeywords {
		if strings.Contains(upperKey, kw) {
			isSensitive = true
			break
		}
	}

	if !isSensitive {
		return val
	}

	if len(val) <= 4 {
		return "******"
	}

	return val[:2] + "******" + val[len(val)-2:]
}
