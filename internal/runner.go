// Copyright RoseSecurity 2026

package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

// Options specifies configuration for running tfrd.
type Options struct {
	Dir             string
	File            string
	Resource        string
	IncludeComputed bool
	JSON            bool
	Writer          io.Writer
}

// MatchesFilter reports whether a resource matches the filter string (resource type or full address).
func MatchesFilter(r *Resource, filter string) bool {
	if strings.Contains(filter, ".") {
		return r.Address() == filter
	}
	return r.Type == filter
}

// Run executes the diff analysis according to the specified options.
func Run(opts Options) error {
	s := spinner.New(
		spinner.CharSets[14],
		100*time.Millisecond,
		spinner.WithWriterFile(os.Stderr),
	)
	s.Suffix = " Analyzing Terraform configuration..."
	s.Color("purple", "bold")
	s.Start()
	spinnerRunning := true
	stopSpinner := func() {
		if spinnerRunning {
			s.Stop()
			spinnerRunning = false
		}
	}
	defer stopSpinner()

	dir := opts.Dir
	if dir == "" {
		dir = "."
	}
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	// 1. Figure out which files to parse.
	var files []string
	if opts.File != "" {
		files = []string{opts.File}
	} else {
		matches, err := filepath.Glob(filepath.Join(dir, "*.tf"))
		if err != nil {
			return err
		}
		files = matches
	}
	if len(files) == 0 {
		return fmt.Errorf("no .tf files found in %q", dir)
	}

	// 2. Parse the config to learn which resources exist and which attributes
	//    each one sets.
	resources, err := ParseResources(files)
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		return fmt.Errorf("no resource blocks found in %v", files)
	}

	// 3. Pull the provider schema from the initialized working directory.
	schemas, err := LoadSchemas(dir)
	if err != nil {
		return err
	}

	// 4. Diff each resource against its schema.
	var reports []ResourceReport
	for _, r := range resources {
		if opts.Resource != "" && !MatchesFilter(r, opts.Resource) {
			continue
		}
		schemaBlock := schemas.ResourceBlock(r.Type)
		if schemaBlock == nil {
			reports = append(reports, ResourceReport{
				Address: r.Address(),
				Type:    r.Type,
				Error:   "no schema found for this resource type (is the provider installed?)",
			})
			continue
		}
		reports = append(reports, BuildReport(r, schemaBlock, opts.IncludeComputed))
	}

	sort.Slice(reports, func(i, j int) bool { return reports[i].Address < reports[j].Address })

	stopSpinner()
	if opts.JSON {
		return PrintJSON(w, reports)
	}
	PrintText(w, reports)
	return nil
}
