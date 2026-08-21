package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

// SchemaSet wraps the provider schemas dumped by Terraform and lets callers
// look up a resource type without caring which provider owns it.
type SchemaSet struct {
	Providers *tfjson.ProviderSchemas
}

// LoadSchemas runs `terraform providers schema -json` in dir and decodes it.
// The directory must be initialized so the providers are on disk.
func LoadSchemas(dir string) (*SchemaSet, error) {
	cmd := exec.Command("terraform", fmt.Sprintf("-chdir=%s", dir), "providers", "schema", "-json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if strings.Contains(msg, "terraform init") || strings.Contains(msg, "Missing") {
			return nil, fmt.Errorf("could not read provider schema — run `terraform init` in %q first\n\n%s", dir, msg)
		}
		return nil, fmt.Errorf("terraform providers schema failed: %v\n%s", err, msg)
	}

	var ps tfjson.ProviderSchemas
	if err := json.Unmarshal(stdout.Bytes(), &ps); err != nil {
		return nil, fmt.Errorf("decoding provider schema: %w", err)
	}
	return &SchemaSet{Providers: &ps}, nil
}

// ResourceBlock finds the schema block for a resource type across all installed
// providers. Resource type prefixes are effectively unique, so the first match
// wins.
func (s *SchemaSet) ResourceBlock(resourceType string) *tfjson.SchemaBlock {
	for _, prov := range s.Providers.Schemas {
		if rs, ok := prov.ResourceSchemas[resourceType]; ok && rs.Block != nil {
			return rs.Block
		}
	}
	return nil
}

// Settable reports whether an attribute can be assigned in configuration.
// Computed-only attributes are read-only outputs and must not be flagged as
// "unused" — you cannot set them.
func Settable(a *tfjson.SchemaAttribute) bool {
	return a.Required || a.Optional
}
