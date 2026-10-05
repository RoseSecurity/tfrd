// Copyright RoseSecurity 2026

package internal

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func attr(required, optional, computed bool) *tfjson.SchemaAttribute {
	return &tfjson.SchemaAttribute{
		AttributeType: cty.String,
		Required:      required,
		Optional:      optional,
		Computed:      computed,
	}
}

// A schema with a plain attribute, a read-only attribute, and a nested block
// that itself has a settable and a required argument.
func sampleSchema() *tfjson.SchemaBlock {
	return &tfjson.SchemaBlock{
		Attributes: map[string]*tfjson.SchemaAttribute{
			"name": attr(true, false, false), // required, settable
			"tags": attr(false, true, false), // optional, settable
			"id":   attr(false, false, true), // computed-only, read-only
			"arn":  attr(false, true, true),  // optional+computed, settable
		},
		NestedBlocks: map[string]*tfjson.SchemaBlockType{
			"ingress": {
				NestingMode: tfjson.SchemaNestingModeList,
				MinItems:    1, // required block
				Block: &tfjson.SchemaBlock{
					Attributes: map[string]*tfjson.SchemaAttribute{
						"port":     attr(true, false, false),
						"protocol": attr(false, true, false),
					},
				},
			},
		},
	}
}

func report(t *testing.T, cfg *CfgBlock, includeComputed bool) ResourceReport {
	t.Helper()
	r := &Resource{Type: "x", Name: "y", Body: cfg}
	return BuildReport(r, sampleSchema(), includeComputed)
}

func findingPaths(rep ResourceReport) map[string]Finding {
	m := map[string]Finding{}
	for _, f := range rep.Unused {
		m[f.Path] = f
	}
	return m
}

func TestComputedOnlyExcludedByDefault(t *testing.T) {
	rep := report(t, NewCfgBlock(), false)
	got := findingPaths(rep)
	if _, ok := got["id"]; ok {
		t.Errorf("read-only attribute 'id' should not be reported without -include-computed")
	}
	// arn is optional+computed → settable → should be reported.
	if _, ok := got["arn"]; !ok {
		t.Errorf("optional+computed attribute 'arn' is settable and should be reported")
	}
}

func TestRequiredFlagged(t *testing.T) {
	rep := report(t, NewCfgBlock(), false)
	got := findingPaths(rep)
	if f := got["name"]; !f.Required {
		t.Errorf("'name' should be flagged required, got %+v", f)
	}
	if f := got["ingress"]; !f.Required || f.Kind != "block" {
		t.Errorf("'ingress' should be a required block, got %+v", f)
	}
}

func TestUsedAttributesNotReported(t *testing.T) {
	cfg := NewCfgBlock()
	cfg.Attrs["name"] = true
	cfg.Attrs["tags"] = true
	rep := report(t, cfg, false)
	got := findingPaths(rep)
	if _, ok := got["name"]; ok {
		t.Errorf("'name' is set and must not be reported unused")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("'tags' is set and must not be reported unused")
	}
}

func TestNestedBlockRecursion(t *testing.T) {
	// ingress present but only sets 'port'; 'protocol' should be reported as
	// ingress.protocol and the block itself counted as used.
	cfg := NewCfgBlock()
	inner := NewCfgBlock()
	inner.Attrs["port"] = true
	cfg.Blocks["ingress"] = []*CfgBlock{inner}

	rep := report(t, cfg, false)
	got := findingPaths(rep)
	if _, ok := got["ingress"]; ok {
		t.Errorf("'ingress' is present and should not be reported as an unused block")
	}
	if _, ok := got["ingress.protocol"]; !ok {
		t.Errorf("'ingress.protocol' should be reported unused; got %v", keys(got))
	}
	if _, ok := got["ingress.port"]; ok {
		t.Errorf("'ingress.port' is set and must not be reported")
	}
}

func TestUnusedBlockStillListsItsArgs(t *testing.T) {
	rep := report(t, NewCfgBlock(), false)
	got := findingPaths(rep)
	// The block is unused, but its arguments should be enumerated so the
	// author sees the whole surface.
	for _, p := range []string{"ingress", "ingress.port", "ingress.protocol"} {
		if _, ok := got[p]; !ok {
			t.Errorf("expected %q in findings; got %v", p, keys(got))
		}
	}
}

func TestCoverageCounts(t *testing.T) {
	cfg := NewCfgBlock()
	cfg.Attrs["name"] = true
	rep := report(t, cfg, false)
	// settable: name, tags, arn, ingress(block), ingress.port, ingress.protocol = 6
	if rep.Total != 6 {
		t.Errorf("expected 6 settable elements, got %d", rep.Total)
	}
	if rep.Used != 1 {
		t.Errorf("expected 1 used, got %d", rep.Used)
	}
}

func keys(m map[string]Finding) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
