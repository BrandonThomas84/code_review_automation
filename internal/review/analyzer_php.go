package review

import (
	"os"
	"path/filepath"
	"strings"
)

// checkPHPQuality analyzes PHP files for quality and security issues
func (a *Analyzer) checkPHPQuality(file string, report *Report) {
	filePath := filepath.Join(a.repoPath, file)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for i, line := range lines {
		lineLower := strings.ToLower(line)

		// Line length check
		if len(line) > 120 {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Line too long (>120 characters)",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for var_dump/print_r debug statements
		if strings.Contains(line, "var_dump(") || strings.Contains(line, "print_r(") || strings.Contains(line, "var_export(") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Debug output (var_dump/print_r) found - remove before production",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for die/exit statements
		if strings.Contains(line, "die(") || strings.Contains(line, "exit(") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "die()/exit() statement found - consider proper error handling",
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

		// SECURITY: Check for eval usage
		if strings.Contains(line, "eval(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "eval() usage detected - potential code injection vulnerability",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for shell_exec/exec/system
		if strings.Contains(line, "shell_exec(") || strings.Contains(line, "exec(") || strings.Contains(line, "system(") || strings.Contains(line, "passthru(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Shell command execution detected - ensure input is sanitized",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for SQL injection vulnerabilities
		if strings.Contains(line, "$_GET") || strings.Contains(line, "$_POST") || strings.Contains(line, "$_REQUEST") {
			if strings.Contains(line, "mysql_query") || strings.Contains(line, "mysqli_query") || strings.Contains(line, "->query(") {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  "Potential SQL injection - use prepared statements",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// Check for deprecated mysql_* functions
		if strings.Contains(line, "mysql_connect") || strings.Contains(line, "mysql_query") || strings.Contains(line, "mysql_fetch") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "Deprecated mysql_* function - use mysqli or PDO instead",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for include/require with user input
		if (strings.Contains(line, "include(") || strings.Contains(line, "require(") || strings.Contains(line, "include_once(") || strings.Contains(line, "require_once(")) &&
			(strings.Contains(line, "$_GET") || strings.Contains(line, "$_POST") || strings.Contains(line, "$_REQUEST")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "File inclusion with user input - potential LFI/RFI vulnerability",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for unserialize with user input
		if strings.Contains(line, "unserialize(") && (strings.Contains(line, "$_GET") || strings.Contains(line, "$_POST") || strings.Contains(line, "$_REQUEST")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Unsafe unserialize with user input - potential object injection",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for XSS vulnerabilities (echo without htmlspecialchars)
		if strings.Contains(line, "echo") && (strings.Contains(line, "$_GET") || strings.Contains(line, "$_POST") || strings.Contains(line, "$_REQUEST")) {
			if !strings.Contains(line, "htmlspecialchars") && !strings.Contains(line, "htmlentities") {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  "Potential XSS - escape output with htmlspecialchars()",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// SECURITY: Check for weak hashing
		if strings.Contains(line, "md5(") || strings.Contains(line, "sha1(") {
			if strings.Contains(lineLower, "password") {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  "Weak password hashing - use password_hash() instead",
					File:     file,
					Line:     i + 1,
				})
			}
		}
	}
}

