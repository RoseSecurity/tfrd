package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

// Finding is one settable schema element the config does not use.
type Finding struct {
	Path     string `json:"path"`               // dotted path, e.g. "ebs_block_device.encrypted"
	Kind     string `json:"kind"`               // "attribute" or "block"
	Type     string `json:"type,omitempty"`     // friendly type for attributes
	Required bool   `json:"required"`           // required arguments are the loud ones
	Computed bool   `json:"computed,omitempty"` // only set when -include-computed surfaces read-only attrs
}

type ResourceReport struct {
	Address  string    `json:"address"`
	Type     string    `json:"type"`
	Total    int       `json:"settable_total"` // settable schema elements
	Used     int       `json:"settable_used"`
	Unused   []Finding `json:"unused"`
	Error    string    `json:"error,omitempty"`
	Coverage float64   `json:"coverage,omitempty"` // used/total, computed for text output
}

func BuildReport(r *Resource, schema *tfjson.SchemaBlock, includeComputed bool) ResourceReport {
	rep := ResourceReport{Address: r.Address(), Type: r.Type}
	Diff(schema, r.Body, "", includeComputed, &rep)
	if rep.Total > 0 {
		rep.Coverage = float64(rep.Used) / float64(rep.Total) * 100
	}
	sort.Slice(rep.Unused, func(i, j int) bool {
		// Required first, then alphabetical — required args are what a module
		// author most needs to notice.
		if rep.Unused[i].Required != rep.Unused[j].Required {
			return rep.Unused[i].Required
		}
		return rep.Unused[i].Path < rep.Unused[j].Path
	})
	return rep
}

// Diff walks a schema block against the presence tree of a config block,
// recording settable elements the config omits and tallying coverage.
func Diff(schema *tfjson.SchemaBlock, cfg *CfgBlock, prefix string, includeComputed bool, rep *ResourceReport) {
	for name, attr := range schema.Attributes {
		readOnly := attr.Computed && !Settable(attr)
		if readOnly && !includeComputed {
			continue
		}
		rep.Total++
		if cfg != nil && cfg.Attrs[name] {
			rep.Used++
			// A nested-attribute type could be partially used, but the config
			// assigns it as a single expression, so we treat presence as full
			// use and don't descend.
			continue
		}
		rep.Unused = append(rep.Unused, Finding{
			Path:     prefix + name,
			Kind:     "attribute",
			Type:     AttrType(attr),
			Required: attr.Required,
			Computed: readOnly,
		})
	}

	for name, blockType := range schema.NestedBlocks {
		var instances []*CfgBlock
		if cfg != nil {
			instances = cfg.Blocks[name]
		}
		rep.Total++
		if len(instances) == 0 {
			rep.Unused = append(rep.Unused, Finding{
				Path:     prefix + name,
				Kind:     "block",
				Required: blockType.MinItems > 0,
			})
			// Still descend so the author sees which arguments the block
			// exposes, all reported as unused under it.
			Diff(blockType.Block, nil, prefix+name+".", includeComputed, rep)
			continue
		}
		rep.Used++
		Diff(blockType.Block, Merge(instances), prefix+name+".", includeComputed, rep)
	}
}

func AttrType(a *tfjson.SchemaAttribute) string {
	if a.AttributeNestedType != nil {
		return "nested"
	}
	if a.AttributeType != cty.NilType {
		return a.AttributeType.FriendlyName()
	}
	return ""
}

func PrintText(w io.Writer, reports []ResourceReport) {
	for i, rep := range reports {
		if i > 0 {
			fmt.Fprintln(w)
		}
		if rep.Error != "" {
			fmt.Fprintf(w, "%s\n  ! %s\n", rep.Address, rep.Error)
			continue
		}
		fmt.Fprintf(w, "%s  (%d/%d settable args used, %.0f%% coverage)\n",
			rep.Address, rep.Used, rep.Total, rep.Coverage)
		if len(rep.Unused) == 0 {
			fmt.Fprintln(w, "  ✓ every settable argument is used")
			continue
		}
		for _, f := range rep.Unused {
			marker := " "
			label := ""
			switch {
			case f.Required:
				marker = "!"
				label = "  (required)"
			case f.Computed:
				label = "  (read-only)"
			}
			typ := ""
			if f.Kind == "block" {
				typ = " {block}"
			} else if f.Type != "" {
				typ = " : " + f.Type
			}
			fmt.Fprintf(w, "  %s %s%s%s\n", marker, f.Path, typ, label)
		}
	}
}

func PrintJSON(w io.Writer, reports []ResourceReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(reports)
}
