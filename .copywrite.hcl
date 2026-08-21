schema_version = 1

project {
  copyright_year = 2026
  copyright_holder = "Sheetz"

  header_ignore = [
    # tests used within documentation (prose)
    "test/**",
     "*.md",
     "docs/**",
    # golangci-lint tooling configuration
    ".golangci.yml",
  ]
}
