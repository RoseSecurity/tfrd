// Copyright RoseSecurity 2026

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate markdown documentation for tfrd",
	Args:   cobra.NoArgs,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		docsDir := "docs"
		if err := os.MkdirAll(docsDir, 0o755); err != nil {
			return fmt.Errorf("creating docs directory: %w", err)
		}
		if err := doc.GenMarkdownTree(rootCmd, docsDir); err != nil {
			return fmt.Errorf("generating docs: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
