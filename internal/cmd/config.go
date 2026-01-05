package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Configuration:")
			fmt.Println("  Target Branch: (set via -t flag)")
			fmt.Println("  Output Directory: review_reports")
			fmt.Println("  Full Scan: false (set via --full-scan flag)")
			fmt.Println("  Email: (set via --email flag)")
			return nil
		},
	})

	return cmd
}

