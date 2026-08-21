# tfrd — terraform resource diff

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

## Usage

```bash
# The directory must be initialized so the provider is on disk.
cd path/to/module
terraform init

tfrd                       # analyze every .tf in the current directory
tfrd ./example             # analyze a specific directory
tfrd -resource aws_instance ./example
tfrd -file main.tf ./example
tfrd -include-computed     # also list read-only attributes
tfrd -json                 # machine-readable output
```

### Example output

```
local_file.example  (2/7 settable args used, 29% coverage)
    content_base64 : string
    directory_permission : string
    file_permission : string
    sensitive_content : string
    source : string
```

Required arguments that are missing are marked with `!`, and unused nested
blocks are listed with their inner arguments so you see the full surface.

## Build

```bash
go build -o tfrd .
go test ./...
```

## Known limitations

- The working directory must be initialized (`terraform init`) so the schema is
  available.
- Nested *attribute* types (map/object attributes with structure) are treated as
  used once the attribute is assigned; tfrd does not descend into the assigned
  expression.
- Coverage counts each nested block as a single settable element plus its inner
  arguments.
