package email

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/BrandonThomas84/code-review-automation/internal/review"
)

// Formatter creates formatted HTML email content from review reports
type Formatter struct {
	RepoName   string
	BranchName string
	PRNumber   int
	PRTitle    string
}

// NewFormatter creates a new email formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// WithRepo sets repository context
func (f *Formatter) WithRepo(repoName string) *Formatter {
	f.RepoName = repoName
	return f
}

// WithBranch sets branch context
func (f *Formatter) WithBranch(branchName string) *Formatter {
	f.BranchName = branchName
	return f
}

// WithPR sets pull request context
func (f *Formatter) WithPR(number int, title string) *Formatter {
	f.PRNumber = number
	f.PRTitle = title
	return f
}

// FormatHTML generates a complete HTML email from the report
func (f *Formatter) FormatHTML(report *review.Report) string {
	var buf bytes.Buffer

	// Write HTML header with styles
	buf.WriteString(f.htmlHeader())

	// Write email body
	buf.WriteString(`<body style="margin: 0; padding: 0; background-color: #f4f4f4;">`)
	buf.WriteString(`<table width="100%" cellpadding="0" cellspacing="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">`)

	// Header banner
	buf.WriteString(f.headerBanner(report))

	// Summary section
	buf.WriteString(f.summarySection(report))

	// Issues section
	if len(report.Issues) > 0 {
		buf.WriteString(f.issuesSection(report))
	} else {
		buf.WriteString(f.noIssuesSection())
	}

	// Footer
	buf.WriteString(f.footer())

	buf.WriteString(`</table></body></html>`)

	return buf.String()
}

func (f *Formatter) htmlHeader() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Code Review Report</title>
</head>`
}

func (f *Formatter) headerBanner(report *review.Report) string {
	bgColor := "#4caf50" // green for no issues
	emoji := "✅"
	status := "All Clear"

	if report.Summary.HighSeverity > 0 {
		bgColor = "#f44336" // red for high severity
		emoji = "🚨"
		status = "Action Required"
	} else if report.Summary.MediumSeverity > 0 {
		bgColor = "#ff9800" // orange for medium
		emoji = "⚠️"
		status = "Review Recommended"
	} else if report.Summary.LowSeverity > 0 {
		bgColor = "#2196f3" // blue for low
		emoji = "ℹ️"
		status = "Minor Issues"
	}

	title := "Code Review Report"
	if f.RepoName != "" {
		title = fmt.Sprintf("Code Review: %s", f.RepoName)
	}

	return fmt.Sprintf(`
<tr>
    <td style="background-color: %s; padding: 30px; text-align: center;">
        <h1 style="color: #ffffff; margin: 0; font-family: Arial, sans-serif; font-size: 24px;">
            %s %s
        </h1>
        <p style="color: #ffffff; margin: 10px 0 0 0; font-family: Arial, sans-serif; font-size: 16px;">
            %s
        </p>
    </td>
</tr>`, bgColor, emoji, html.EscapeString(title), status)
}

func (f *Formatter) summarySection(report *review.Report) string {
	var context string
	if f.BranchName != "" {
		context = fmt.Sprintf("<p style=\"margin: 5px 0; color: #666;\">Branch: <strong>%s</strong></p>", html.EscapeString(f.BranchName))
	}
	if f.PRNumber > 0 {
		context += fmt.Sprintf("<p style=\"margin: 5px 0; color: #666;\">PR #%d: %s</p>", f.PRNumber, html.EscapeString(f.PRTitle))
	}

	return fmt.Sprintf(`
<tr>
    <td style="padding: 20px; font-family: Arial, sans-serif;">
        <h2 style="color: #333; margin: 0 0 15px 0; font-size: 18px;">📊 Summary</h2>
        %s
        <table width="100%%" cellpadding="10" cellspacing="0" style="background-color: #f9f9f9; border-radius: 8px; margin-top: 10px;">
            <tr>
                <td style="text-align: center; border-right: 1px solid #ddd;">
                    <div style="font-size: 28px; font-weight: bold; color: #333;">%d</div>
                    <div style="font-size: 12px; color: #666;">Files Changed</div>
                </td>
                <td style="text-align: center; border-right: 1px solid #ddd;">
                    <div style="font-size: 28px; font-weight: bold; color: #f44336;">%d</div>
                    <div style="font-size: 12px; color: #666;">High</div>
                </td>
                <td style="text-align: center; border-right: 1px solid #ddd;">
                    <div style="font-size: 28px; font-weight: bold; color: #ff9800;">%d</div>
                    <div style="font-size: 12px; color: #666;">Medium</div>
                </td>
                <td style="text-align: center;">
                    <div style="font-size: 28px; font-weight: bold; color: #4caf50;">%d</div>
                    <div style="font-size: 12px; color: #666;">Low</div>
                </td>
            </tr>
        </table>
    </td>
</tr>`, context, report.Summary.TotalFiles, report.Summary.HighSeverity, report.Summary.MediumSeverity, report.Summary.LowSeverity)
}

func (f *Formatter) issuesSection(report *review.Report) string {
	var buf bytes.Buffer

	buf.WriteString(`
<tr>
    <td style="padding: 0 20px 20px 20px; font-family: Arial, sans-serif;">
        <h2 style="color: #333; margin: 0 0 15px 0; font-size: 18px;">🔍 Issues Found</h2>`)

	// Group issues by severity
	highIssues := filterBySeverity(report.Issues, "high")
	mediumIssues := filterBySeverity(report.Issues, "medium")
	lowIssues := filterBySeverity(report.Issues, "low")

	// Render high severity first
	if len(highIssues) > 0 {
		buf.WriteString(f.issueGroup("High Severity", "#f44336", highIssues))
	}
	if len(mediumIssues) > 0 {
		buf.WriteString(f.issueGroup("Medium Severity", "#ff9800", mediumIssues))
	}
	if len(lowIssues) > 0 {
		buf.WriteString(f.issueGroup("Low Severity", "#4caf50", lowIssues))
	}

	buf.WriteString(`</td></tr>`)
	return buf.String()
}

func filterBySeverity(issues []review.Issue, severity string) []review.Issue {
	var filtered []review.Issue
	for _, issue := range issues {
		if strings.ToLower(issue.Severity) == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func (f *Formatter) issueGroup(title, color string, issues []review.Issue) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`
        <div style="margin-bottom: 15px;">
            <div style="background-color: %s; color: white; padding: 8px 12px; border-radius: 4px 4px 0 0; font-weight: bold; font-size: 14px;">
                %s (%d)
            </div>
            <div style="border: 1px solid #ddd; border-top: none; border-radius: 0 0 4px 4px;">`, color, title, len(issues)))

	maxIssues := 10 // Limit per group to keep email manageable
	displayIssues := issues
	if len(issues) > maxIssues {
		displayIssues = issues[:maxIssues]
	}

	for _, issue := range displayIssues {
		location := html.EscapeString(issue.File)
		if issue.Line > 0 {
			location += fmt.Sprintf(":%d", issue.Line)
		}

		buf.WriteString(fmt.Sprintf(`
                <div style="padding: 12px; border-bottom: 1px solid #eee;">
                    <div style="font-size: 14px; color: #333; margin-bottom: 5px;">%s</div>
                    <div style="font-size: 12px; color: #666;">
                        <code style="background-color: #f5f5f5; padding: 2px 6px; border-radius: 3px;">%s</code>
                    </div>
                </div>`, html.EscapeString(issue.Message), location))
	}

	if len(issues) > maxIssues {
		buf.WriteString(fmt.Sprintf(`
                <div style="padding: 12px; background-color: #f9f9f9; text-align: center; color: #666; font-size: 12px;">
                    ...and %d more issues
                </div>`, len(issues)-maxIssues))
	}

	buf.WriteString(`</div></div>`)
	return buf.String()
}

func (f *Formatter) noIssuesSection() string {
	return `
<tr>
    <td style="padding: 20px; text-align: center; font-family: Arial, sans-serif;">
        <div style="background-color: #e8f5e9; border-radius: 8px; padding: 30px;">
            <div style="font-size: 48px; margin-bottom: 10px;">✅</div>
            <h3 style="color: #2e7d32; margin: 0;">No Issues Found!</h3>
            <p style="color: #666; margin: 10px 0 0 0;">Great job! Your code passed all quality and security checks.</p>
        </div>
    </td>
</tr>`
}

func (f *Formatter) footer() string {
	timestamp := time.Now().Format("January 2, 2006 at 3:04 PM")
	return fmt.Sprintf(`
<tr>
    <td style="padding: 20px; background-color: #f9f9f9; text-align: center; font-family: Arial, sans-serif;">
        <p style="color: #999; font-size: 12px; margin: 0;">
            Generated on %s<br>
            <a href="https://github.com/BrandonThomas84/code_review_automation" style="color: #2196f3;">Code Review Automation</a>
        </p>
    </td>
</tr>`, timestamp)
}

// FormatSubject generates an appropriate email subject line
func (f *Formatter) FormatSubject(report *review.Report) string {
	var prefix string
	if report.Summary.HighSeverity > 0 {
		prefix = "🚨 "
	} else if report.Summary.MediumSeverity > 0 {
		prefix = "⚠️ "
	} else if report.Summary.TotalIssues > 0 {
		prefix = "ℹ️ "
	} else {
		prefix = "✅ "
	}

	subject := fmt.Sprintf("%sCode Review: %d issues found", prefix, report.Summary.TotalIssues)

	if f.RepoName != "" {
		subject = fmt.Sprintf("%sCode Review [%s]: %d issues found", prefix, f.RepoName, report.Summary.TotalIssues)
	}

	if f.PRNumber > 0 {
		subject = fmt.Sprintf("%sCode Review PR #%d: %d issues found", prefix, f.PRNumber, report.Summary.TotalIssues)
	}

	return subject
}
