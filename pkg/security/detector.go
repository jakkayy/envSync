package security

import (
	"regexp"
	"strings"
)

type SecretFinding struct {
	Key         string `json:"key"`
	PatternType string `json:"pattern_type"`
	Description string `json:"description"`
}

var (
	awsKeyRegex     = regexp.MustCompile(`^(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}$`)
	jwtRegex        = regexp.MustCompile(`^eyJ[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`)
	privateKeyRegex = regexp.MustCompile(`-----BEGIN [A-Z]+ PRIVATE KEY-----`)
	dbURLRegex      = regexp.MustCompile(`(?i)(postgres|mysql|mongodb|redis):\/\/[^:]+:[^@]+@`)
	slackTokenRegex = regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`)
)

// DetectSecrets scans environment key-value pairs for high-risk plaintext secrets
func DetectSecrets(envMap map[string]string) []SecretFinding {
	var findings []SecretFinding

	for k, v := range envMap {
		val := strings.TrimSpace(v)
		if val == "" {
			continue
		}

		if awsKeyRegex.MatchString(val) {
			findings = append(findings, SecretFinding{
				Key:         k,
				PatternType: "AWS_ACCESS_KEY",
				Description: "High risk: Hardcoded AWS Access Key ID detected",
			})
		} else if jwtRegex.MatchString(val) {
			findings = append(findings, SecretFinding{
				Key:         k,
				PatternType: "JWT_TOKEN",
				Description: "Medium risk: Plaintext JWT Token detected",
			})
		} else if privateKeyRegex.MatchString(val) {
			findings = append(findings, SecretFinding{
				Key:         k,
				PatternType: "PRIVATE_KEY",
				Description: "Critical risk: Raw Private Key block detected",
			})
		} else if dbURLRegex.MatchString(val) {
			findings = append(findings, SecretFinding{
				Key:         k,
				PatternType: "DB_CONNECTION_STRING",
				Description: "High risk: Database Connection String contains plaintext credentials",
			})
		} else if slackTokenRegex.MatchString(val) {
			findings = append(findings, SecretFinding{
				Key:         k,
				PatternType: "SLACK_TOKEN",
				Description: "High risk: Slack API Token detected",
			})
		}
	}

	return findings
}
