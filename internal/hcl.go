package internal

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Resource is a single resource block found in the configuration, reduced to
// the structure tfrd cares about: its type, name, and a tree of which
// attributes and nested blocks it sets.
type Resource struct {
	Type string
	Name string
	Body *CfgBlock
}

func (r *Resource) Address() string { return r.Type + "." + r.Name }

// CfgBlock captures which attributes are assigned and which nested blocks are
// present within a block body, recursively. We only track presence (the set of
// names), never the expression values, because coverage only depends on whether
// an argument appears at all.
type CfgBlock struct {
	Attrs  map[string]bool
	Blocks map[string][]*CfgBlock
}

func NewCfgBlock() *CfgBlock {
	return &CfgBlock{Attrs: map[string]bool{}, Blocks: map[string][]*CfgBlock{}}
}

func ParseResources(files []string) ([]*Resource, error) {
	parser := hclparse.NewParser()
	var out []*Resource
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		file, diags := parser.ParseHCL(src, f)
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing %s: %s", f, diags.Error())
		}
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}
		for _, b := range body.Blocks {
			if b.Type != "resource" || len(b.Labels) != 2 {
				continue
			}
			out = append(out, &Resource{
				Type: b.Labels[0],
				Name: b.Labels[1],
				Body: BuildCfg(b.Body),
			})
		}
	}
	return out, nil
}

// BuildCfg turns an HCL block body into a CfgBlock, descending into nested
// blocks. dynamic blocks are unwrapped so the arguments they generate are
// counted against the schema of the block they stand in for.
func BuildCfg(body *hclsyntax.Body) *CfgBlock {
	c := NewCfgBlock()
	for name := range body.Attributes {
		c.Attrs[name] = true
	}
	for _, b := range body.Blocks {
		if b.Type == "dynamic" {
			if len(b.Labels) != 1 {
				continue
			}
			name := b.Labels[0]
			for _, inner := range b.Body.Blocks {
				if inner.Type == "content" {
					c.Blocks[name] = append(c.Blocks[name], BuildCfg(inner.Body))
				}
			}
			continue
		}
		c.Blocks[b.Type] = append(c.Blocks[b.Type], BuildCfg(b.Body))
	}
	return c
}

// Merge unions several block instances into one presence tree. When a nested
// block appears multiple times (e.g. several ingress blocks), an argument is
// "used" if any instance sets it — the right semantic for "does the module ever
// expose this argument?".
func Merge(instances []*CfgBlock) *CfgBlock {
	m := NewCfgBlock()
	for _, inst := range instances {
		for name := range inst.Attrs {
			m.Attrs[name] = true
		}
		for name, blks := range inst.Blocks {
			m.Blocks[name] = append(m.Blocks[name], blks...)
		}
	}
	return m
}
