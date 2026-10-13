// Copyright RoseSecurity 2026

package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSchemasRunsInitWhenRequested(t *testing.T) {
	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "terraform.log")
	terraform := filepath.Join(binDir, "terraform")
	script := `#!/bin/sh
printf '%s\n' "$*" >> "$TF_LOG"
if [ "$1" = "init" ]; then
  exit 0
fi
printf '{"format_version":"1.0","provider_schemas":{}}\n'
`
	if err := os.WriteFile(terraform, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(terraform, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("TF_LOG", logPath)

	if _, err := LoadSchemas(t.TempDir(), true); err != nil {
		t.Fatalf("LoadSchemas() error = %v", err)
	}

	log, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("reading terraform invocation log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(log)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected init and schema invocations, got %q", log)
	}
	if !strings.HasPrefix(lines[0], "init -backend=false") {
		t.Errorf("first terraform invocation should initialize without a backend, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "providers schema -json") {
		t.Errorf("second terraform invocation should load provider schemas, got %q", lines[1])
	}
}
