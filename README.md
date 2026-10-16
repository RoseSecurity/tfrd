<p align="center">
  <img src="https://github.com/user-attachments/assets/f4ca875c-c687-442f-ad2e-0d3b28ec28f1" alt="tfrd logo" width="550">
</p>

---

<p align="center">
  <em>Terraform resource diffs for identifying unused attributes and parameters.</em>
</p>

## Why I wrote this tool

In my limited free time, I enjoy maintaining community Terraform modules.
Digging through the Terraform Registry documentation and switching between the
browser and terminal can be tedious, so I built `tfrd` to compare the resources
in a Terraform configuration with the provider schema and report unused
attributes and parameters. As a community contributor, I want my modules to
reflect the full range of functionality providers support, and this tool helps
me find gaps.

## Demo

<img width="1200" height="550" alt="tfrd-demo" src="https://github.com/user-attachments/assets/fa40eab2-cb4b-4ea3-a645-c6180caa0f21" />

## Install

### Go

If you have a functional Go environment, you can install with:

```sh
go install github.com/RoseSecurity/tfrd@latest
```

### Apt

To install packages, you can quickly setup the repository automatically:

```sh
curl -1sLf \
  'https://dl.cloudsmith.io/public/rosesecurity/tfrd/setup.deb.sh' \
  | sudo -E bash
```

Once the repository is configured, you can install with:

```sh
apt install tfrd
```

### Source

```sh
git clone https://github.com/RoseSecurity/tfrd.git
cd tfrd
make build
```

## Usage

Run `tfrd` from an initialized Terraform directory, or pass the directory as an
argument. Use `--init` to run `terraform init -backend=false` automatically
before loading the provider schema:

```sh
tfrd
tfrd ./examples
tfrd --init ./examples
```

By default, `tfrd` analyzes `.tf` files directly inside the selected directory.
Use `--file` to analyze one file, `--resource` to filter by resource type or
address, and `--json` for machine-readable output:

```sh
tfrd --file ./examples/main.tf ./examples
tfrd --resource local_file ./examples
tfrd --json ./examples
```

Read-only computed attributes are excluded by default. Add
`--include-computed` to include them in the report.
Add `--include-descriptions` to include provider-supplied descriptions for
unused attributes in text and JSON output.

Run `tfrd --help` for all available options.

## Contributing

Bug reports and feature requests are welcome in the
[issue tracker](https://github.com/RoseSecurity/tfrd/issues). Pull requests are
welcome; please fork the repository, make your changes on a branch, and open a
pull request.
