package review

import (
	"os"
	"path/filepath"
	"strings"
)

// checkPythonQuality analyzes Python files for quality and security issues
func (a *Analyzer) checkPythonQuality(file string, report *Report) {
	filePath := filepath.Join(a.repoPath, file)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for i, line := range lines {
		lineLower := strings.ToLower(line)
		trimmed := strings.TrimSpace(line)

		// Line length check (PEP 8 recommends 79, but 120 is common)
		if len(line) > 120 {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Line too long (>120 characters)",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for print statements (should use logging in production)
		if strings.HasPrefix(trimmed, "print(") || strings.HasPrefix(trimmed, "print (") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "print() statement found - consider using logging instead",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for pdb/debugger statements
		if strings.Contains(line, "import pdb") || strings.Contains(line, "pdb.set_trace()") || strings.Contains(line, "breakpoint()") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "Debugger statement found - remove before production",
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
		if strings.Contains(line, "eval(") || strings.Contains(line, "exec(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "eval()/exec() usage detected - potential code injection vulnerability",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for subprocess with shell=True
		if strings.Contains(line, "subprocess") && strings.Contains(line, "shell=True") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "subprocess with shell=True - potential command injection risk",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for os.system
		if strings.Contains(line, "os.system(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "os.system() usage - consider using subprocess with proper escaping",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for bare except clauses
		if trimmed == "except:" {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "Bare except clause - specify the exception type",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for type: ignore comments
		if strings.Contains(line, "# type: ignore") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Type ignore comment found - consider fixing the type error",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for pickle (unsafe deserialization)
		if strings.Contains(line, "pickle.load") || strings.Contains(line, "pickle.loads") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "pickle.load() is unsafe - can execute arbitrary code during deserialization",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for yaml.load without Loader
		if strings.Contains(line, "yaml.load(") && !strings.Contains(line, "Loader=") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "yaml.load() without safe Loader - use yaml.safe_load() or specify Loader=yaml.SafeLoader",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for SQL string formatting
		if (strings.Contains(line, "execute(") || strings.Contains(line, "executemany(")) && (strings.Contains(line, "%") || strings.Contains(line, ".format(") || strings.Contains(line, "f\"") || strings.Contains(line, "f'")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential SQL injection - use parameterized queries instead of string formatting",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for hardcoded passwords/secrets
		if strings.Contains(lineLower, "password") && strings.Contains(line, "=") && (strings.Contains(line, "\"") || strings.Contains(line, "'")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential hardcoded password - use environment variables",
				File:     file,
				Line:     i + 1,
			})
		}
	}
}

