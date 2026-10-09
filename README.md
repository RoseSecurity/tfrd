# tfrd

<p align="center">
  <em>Terraform resource diffs for identifying unused attributes and parameters.</em>
</p>

## What I wrote this tool

Reports the **settable schema attributes a resource does not use**. It's a
coverage tool for Terraform/OpenTofu module authors: point it at a module and it
tells you which arguments of the underlying resource you haven't exposed yet.

It works by diffing the provider's published schema
(`terraform providers schema -json`) against the attributes actually assigned in
your configuration.

## How it works

1. Parses your `.tf` files with HCL to collect, per resource, the set of
   attributes and nested blocks that are assigned (`dynamic` blocks are
   unwrapped so their generated arguments count).
2. Reads the provider schema from the **initialized** working directory.
3. Diffs the two, recursively into nested blocks, and reports what's missing.

Read-only (computed-only) attributes are excluded by default — you can't set
them, so they aren't "unused". Meta-arguments (`count`, `for_each`, `lifecycle`,
…) are ignored because they aren't part of the resource schema.

## Installation

## Usage

