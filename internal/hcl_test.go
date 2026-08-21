package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTF(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "main.tf")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseResourceAttrsAndBlocks(t *testing.T) {
	path := writeTF(t, `
resource "aws_instance" "web" {
  ami           = "ami-123"
  instance_type = "t3.micro"

  root_block_device {
    volume_size = 20
  }
}
`)
	resources, err := ParseResources([]string{path})
	if err != nil {
		t.Fatal(err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	r := resources[0]
	if r.Address() != "aws_instance.web" {
		t.Errorf("address = %q", r.Address())
	}
	if !r.Body.Attrs["ami"] || !r.Body.Attrs["instance_type"] {
		t.Errorf("expected ami and instance_type attrs, got %v", r.Body.Attrs)
	}
	rbd := r.Body.Blocks["root_block_device"]
	if len(rbd) != 1 || !rbd[0].Attrs["volume_size"] {
		t.Errorf("expected root_block_device.volume_size, got %+v", rbd)
	}
}

func TestParseDynamicBlockUnwrapped(t *testing.T) {
	// A dynamic block should be counted as an instance of the block it
	// generates, and the args inside its content{} should register as used.
	path := writeTF(t, `
resource "aws_security_group" "sg" {
  dynamic "ingress" {
    for_each = var.rules
    content {
      from_port = ingress.value.from
      to_port   = ingress.value.to
    }
  }
}
`)
	resources, err := ParseResources([]string{path})
	if err != nil {
		t.Fatal(err)
	}
	ing := resources[0].Body.Blocks["ingress"]
	if len(ing) != 1 {
		t.Fatalf("expected dynamic ingress to yield 1 block instance, got %d", len(ing))
	}
	if !ing[0].Attrs["from_port"] || !ing[0].Attrs["to_port"] {
		t.Errorf("dynamic content args not captured: %v", ing[0].Attrs)
	}
}

func TestMergeUnionsInstances(t *testing.T) {
	a := NewCfgBlock()
	a.Attrs["x"] = true
	b := NewCfgBlock()
	b.Attrs["y"] = true
	m := Merge([]*CfgBlock{a, b})
	if !m.Attrs["x"] || !m.Attrs["y"] {
		t.Errorf("merge should union attrs, got %v", m.Attrs)
	}
}
