// Copyright RoseSecurity 2026

package cmd

import (
	"log/slog"
	"os"

	"github.com/RoseSecurity/tfrd/internal"
	"github.com/spf13/cobra"
)

var (
	// Version is injected during build via ldflags.
	Version = "dev"

	file            string
	resource        string
	includeComputed bool
	jsonOut         bool
)

var rootCmd = &cobra.Command{
	Use:   "tfrd [flags] [dir]",
	Short: "Terraform resource diffs for identifying unused attributes and parameters",
	Long: `Reports unused settable attributes for the resources in a Terraform
configuration. The directory must already be initialized (terraform init)
so that the provider schema is available.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       Version,
	Args:          cobra.MaximumNArgs(1),
	RunE:          runRoot,
}

func init() {
	// Disable auto-generated string from documentation so that documentation is cleanly built and updated
	rootCmd.DisableAutoGenTag = true

	rootCmd.Flags().StringVarP(&file, "file", "f", "", "analyze a single .tf file instead of every .tf in the directory")
	rootCmd.Flags().StringVarP(&resource, "resource", "r", "", "only analyze resources of this type (e.g. aws_instance) or address (aws_instance.web)")
	rootCmd.Flags().BoolVarP(&includeComputed, "include-computed", "c", false, "also report read-only (computed-only) attributes")
	rootCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "emit JSON instead of a text report")
}

func runRoot(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
	}

	if len(args) > 0 {
		dir = args[0]
	}

	return internal.Run(internal.Options{
		Dir:             dir,
		File:            file,
		Resource:        resource,
		IncludeComputed: includeComputed,
		JSON:            jsonOut,
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
