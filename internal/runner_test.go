// Copyright RoseSecurity 2026

package internal

import (
	"bytes"
	"strings"
	"testing"
)

func TestMatchesFilter(t *testing.T) {
	r := &Resource{
		Type: "aws_instance",
		Name: "web",
	}

	// Filter by type
	if !MatchesFilter(r, "aws_instance") {
		t.Errorf("expected match for type 'aws_instance'")
	}
	if MatchesFilter(r, "aws_s3_bucket") {
		t.Errorf("expected no match for type 'aws_s3_bucket'")
	}

	// Filter by full address
	if !MatchesFilter(r, "aws_instance.web") {
		t.Errorf("expected match for address 'aws_instance.web'")
	}
	if MatchesFilter(r, "aws_instance.db") {
		t.Errorf("expected no match for address 'aws_instance.db'")
	}
}

func TestPrintTextAndJSON(t *testing.T) {
	reports := []ResourceReport{
		{
			Address:  "local_file.example",
			Type:     "local_file",
			Total:    5,
			Used:     2,
			Coverage: 40.0,
			Unused: []Finding{
				{
					Path:        "file_permission",
					Kind:        "attribute",
					Type:        "string",
					Description: "The permissions for the file.",
					Required:    false,
				},
				{
					Path:     "content_base64",
					Kind:     "attribute",
					Type:     "string",
					Required: false,
				},
			},
		},
		{
			Address: "aws_instance.broken",
			Type:    "aws_instance",
			Error:   "no schema found for this resource type",
		},
		{
			Address:  "local_file.full",
			Type:     "local_file",
			Total:    2,
			Used:     2,
			Coverage: 100.0,
		},
	}

	var buf bytes.Buffer
	PrintText(&buf, reports)
	output := buf.String()
	if !strings.Contains(output, "local_file.example") ||
		!strings.Contains(output, "2/5 used") ||
		!strings.Contains(output, "40% coverage") {
		t.Errorf("unexpected text output: %s", output)
	}
	if !strings.Contains(output, "no schema found for this resource type") {
		t.Errorf("expected error in text output: %s", output)
	}
	if !strings.Contains(output, "all settable arguments used") {
		t.Errorf("expected full coverage message in text output: %s", output)
	}
	if !strings.Contains(output, "The permissions for the file.") {
		t.Errorf("expected attribute description in text output: %s", output)
	}

	buf.Reset()
	if err := PrintJSON(&buf, reports); err != nil {
		t.Fatalf("PrintJSON error: %v", err)
	}
	jsonOut := buf.String()
	if !strings.Contains(jsonOut, `"address": "local_file.example"`) {
		t.Errorf("unexpected json output: %s", jsonOut)
	}
	if !strings.Contains(jsonOut, `"description": "The permissions for the file."`) {
		t.Errorf("expected attribute description in JSON output: %s", jsonOut)
	}
}

func TestRunNoTfFiles(t *testing.T) {
	dir := t.TempDir()
	err := Run(Options{Dir: dir})
	if err == nil {
		t.Errorf("expected error for directory with no .tf files")
	}
}
