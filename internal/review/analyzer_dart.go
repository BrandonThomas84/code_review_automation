package review

import (
	"os"
	"path/filepath"
	"strings"
)

// checkDartQuality analyzes Dart files for quality and security issues
func (a *Analyzer) checkDartQuality(file string, report *Report) {
	filePath := filepath.Join(a.repoPath, file)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for i, line := range lines {
		lineLower := strings.ToLower(line)

		// Line length check (Dart style guide recommends 80, but 120 is common)
		if len(line) > 120 {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Line too long (>120 characters)",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for print statements
		if strings.Contains(line, "print(") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "print() statement found - remove before production",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for debugPrint statements
		if strings.Contains(line, "debugPrint(") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "debugPrint() statement found - remove before production",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for TODO/FIXME comments
		if strings.Contains(lineLower, "todo") || strings.Contains(lineLower, "fixme") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "TODO/FIXME comment found",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for dynamic type usage
		if strings.Contains(line, ": dynamic") || strings.Contains(line, "<dynamic>") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "Avoid using 'dynamic' type - use specific types instead",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for ignore directives
		if strings.Contains(line, "// ignore:") || strings.Contains(line, "// ignore_for_file:") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "Dart ignore directive found - consider fixing the issue",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for hardcoded URLs/API endpoints
		if (strings.Contains(line, "http://") || strings.Contains(line, "https://")) && strings.Contains(lineLower, "api") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Hardcoded API URL - consider using environment configuration",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for hardcoded credentials
		if strings.Contains(lineLower, "password") || strings.Contains(lineLower, "apikey") || strings.Contains(lineLower, "api_key") {
			if strings.Contains(line, "=") && (strings.Contains(line, "\"") || strings.Contains(line, "'")) {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  "Potential hardcoded credential - use secure storage",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// SECURITY: Check for insecure HTTP usage (non-HTTPS)
		if strings.Contains(line, "http://") && !strings.Contains(line, "localhost") && !strings.Contains(line, "127.0.0.1") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Insecure HTTP URL - use HTTPS for production",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for disabled SSL certificate verification
		if strings.Contains(line, "badCertificateCallback") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Custom certificate callback - ensure SSL verification is not disabled",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for force unwrap (!) which can cause runtime crashes
		if strings.Contains(line, "!") && !strings.Contains(line, "!=") && !strings.Contains(line, "//") {
			// Simple heuristic - might have false positives
			if strings.Contains(line, "!.") || strings.Contains(line, "!)") || strings.Contains(line, "!;") {
				report.AddIssue(Issue{
					Type:     "quality",
					Severity: "medium",
					Message:  "Force unwrap (!) used - consider null safety patterns",
					File:     file,
					Line:     i + 1,
				})
			}
		}
	}
}

