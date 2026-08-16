package env

import (
	"strings"
	"testing"
)

func TestParseAndFormat(t *testing.T) {
	input := `# Server Config
DB_HOST=localhost
DB_PORT=5432
export API_KEY="secret-123"

# Feature Flags
ENABLE_FEATURE=true
MULTILINE_VAL="hello
world"
`

	reader := strings.NewReader(input)
	envFile, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	envMap := envFile.Map()

	if envMap["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", envMap["DB_HOST"])
	}
	if envMap["API_KEY"] != "secret-123" {
		t.Errorf("expected API_KEY=secret-123, got %s", envMap["API_KEY"])
	}
	if envMap["ENABLE_FEATURE"] != "true" {
		t.Errorf("expected ENABLE_FEATURE=true, got %s", envMap["ENABLE_FEATURE"])
	}

	formatted := Format(envFile)
	if !strings.Contains(formatted, "# Server Config") {
		t.Errorf("formatter lost comments")
	}
}

func TestCompareEnvs(t *testing.T) {
	local := map[string]string{
		"DB_HOST": "localhost",
		"NEW_VAR": "new_val",
	}

	remote := map[string]string{
		"DB_HOST":     "dev-db.com",
		"REMOVED_VAR": "old_val",
	}

	diff := CompareMaps(local, remote)

	if !diff.HasChanges {
		t.Errorf("expected changes, got none")
	}
	if diff.Added != 1 {
		t.Errorf("expected 1 added item, got %d", diff.Added)
	}
	if diff.Modified != 1 {
		t.Errorf("expected 1 modified item, got %d", diff.Modified)
	}
	if diff.Removed != 1 {
		t.Errorf("expected 1 removed item, got %d", diff.Removed)
	}
}
