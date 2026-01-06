package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BrandonThomas84/code-review-automation/internal/review"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	targetBranch string
	outputDir    string
	jsonOutput   bool
	fullScan     bool
	emailTo      string
	verbose      bool
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code-review",
		Short: "Automated code review tool for multiple languages",
		Long: `Code Review Automation - A comprehensive code review tool that analyzes
code changes across multiple languages including Python, JavaScript, TypeScript,
Dart, Ruby, PHP, and Java.`,
		RunE: runReview,
	}

	cmd.Flags().StringVarP(&targetBranch, "target", "t", "", "Target branch to compare against (required)")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "review_reports", "Output directory for reports")
	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
	cmd.Flags().BoolVar(&fullScan, "full-scan", false, "Scan entire codebase instead of just changed files")
	cmd.Flags().StringVar(&emailTo, "email", "", "Email address to send report to")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	cmd.MarkFlagRequired("target")

	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewConfigCommand())

	return cmd
}

func runReview(cmd *cobra.Command, args []string) error {
	if verbose {
		color.Blue("[INFO] Starting code review analysis...")
		color.Blue("[INFO] Target branch: %s", targetBranch)
		color.Blue("[INFO] Full scan: %v", fullScan)
		color.Blue("[INFO] Output directory: %s", outputDir)
		color.Blue("[INFO] JSON output: %v", jsonOutput)
		color.Blue("[INFO] Email: %s", emailTo)

		color.Blue("[INFO] creating output directory: %s", outputDir)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if verbose {
		color.Blue("[INFO] Getting current working directory...")
	}

	// Get current working directory
	repoPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if verbose {
		color.Blue("[INFO] Repository path: %s", repoPath)
	}

	// Run the review
	analyzer := review.NewAnalyzer(repoPath, verbose)
	report, err := analyzer.GenerateReport(targetBranch, fullScan)
	if err != nil {
		return fmt.Errorf("review failed: %w", err)
	}

	if verbose {
		color.Blue("[INFO] Review complete")
	}

	// Output results
	if jsonOutput {
		if verbose {
			color.Blue("[INFO] Outputting JSON...")
		}

		if err := report.OutputJSON(os.Stdout); err != nil {
			return fmt.Errorf("failed to output JSON: %w", err)
		}
	} else {
		if verbose {
			color.Blue("[INFO] Outputting report...")
		}

		report.PrintReport()
	}

	if verbose {
		color.Blue("[INFO] Saving report to file...")
	}

	// Save report to file
	reportPath := filepath.Join(outputDir, "review_report.json")
	if err := report.SaveToFile(reportPath); err != nil {
		color.Yellow("[WARNING] Failed to save report: %v", err)
	} else if verbose {
		color.Green("[SUCCESS] Report saved to: %s", reportPath)
	}

	if verbose {
		color.Blue("[INFO] Sending email...")
	}

	// Send email if requested
	if emailTo != "" {
		if err := sendEmailReport(report, emailTo); err != nil {
			color.Yellow("[WARNING] Failed to send email: %v", err)
		} else if verbose {
			color.Green("[SUCCESS] Email sent to: %s", emailTo)
		}
	} else if verbose {
		color.Blue("[INFO] No email requested")
	}

	return nil
}

func sendEmailReport(report *review.Report, emailTo string) error {
	// Email functionality will be implemented in a separate module
	color.Blue("[INFO] Email functionality coming soon")
	return nil
}
