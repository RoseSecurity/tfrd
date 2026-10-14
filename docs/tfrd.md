## tfrd

Terraform resource diffs for identifying unused attributes and parameters

### Synopsis

Reports unused settable attributes for the resources in a Terraform
configuration. The directory must be initialized (terraform init) so that
the provider schema is available, or pass --init to initialize it
automatically before loading the provider schema.

```
tfrd [flags] [dir]
```

### Options

```
  -f, --file string            analyze a single .tf file instead of every .tf in the directory
  -h, --help                   help for tfrd
  -c, --include-computed       also report read-only (computed-only) attributes
  -d, --include-descriptions   include attribute descriptions in the report
  -i, --init                   run 'terraform init' before loading provider schemas
  -j, --json                   emit JSON instead of a text report
  -r, --resource string        only analyze resources of this type (e.g. aws_instance) or address (aws_instance.web)
```

### SEE ALSO

* [tfrd completion](tfrd_completion.md)	 - Generate the autocompletion script for the specified shell

