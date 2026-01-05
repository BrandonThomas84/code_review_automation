package review

import (
	"os"
	"path/filepath"
	"strings"
)

// checkTypeScriptQuality analyzes TypeScript files for quality and security issues
func (a *Analyzer) checkTypeScriptQuality(file string, report *Report) {
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

		// Check for console.log statements
		if strings.Contains(line, "console.log") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "console.log statement found - remove before production",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for debugger statements
		if strings.Contains(strings.TrimSpace(line), "debugger") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "debugger statement found - remove before production",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for 'any' type usage
		if strings.Contains(line, ": any") || strings.Contains(line, "<any>") || strings.Contains(line, "as any") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "Avoid using 'any' type - use specific types instead",
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

		// Check for @ts-ignore usage
		if strings.Contains(line, "@ts-ignore") || strings.Contains(line, "@ts-nocheck") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "TypeScript ignore directive found - consider fixing the type error",
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

		// SECURITY: Check for Function constructor
		if strings.Contains(line, "new Function(") || strings.Contains(line, "Function(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Function constructor usage - similar risks to eval()",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for innerHTML/dangerouslySetInnerHTML (XSS vulnerability)
		if strings.Contains(line, ".innerHTML") || strings.Contains(line, "dangerouslySetInnerHTML") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "innerHTML/dangerouslySetInnerHTML usage - potential XSS vulnerability",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for document.write (XSS vulnerability)
		if strings.Contains(line, "document.write") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "document.write usage - potential XSS vulnerability",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for child_process usage
		if strings.Contains(line, "child_process") || strings.Contains(line, "exec(") || strings.Contains(line, "execSync(") || strings.Contains(line, "spawn(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "child_process/exec usage - ensure input is sanitized to prevent command injection",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for insecure randomness
		if strings.Contains(line, "Math.random()") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Math.random() is not cryptographically secure - use crypto.randomBytes() for security-sensitive operations",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for disabled SSL verification
		if strings.Contains(line, "rejectUnauthorized: false") || strings.Contains(line, "NODE_TLS_REJECT_UNAUTHORIZED") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "SSL verification disabled - vulnerable to man-in-the-middle attacks",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for hardcoded JWT secrets
		if strings.Contains(lineLower, "jwt") && (strings.Contains(line, "secret") || strings.Contains(line, "Secret")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential hardcoded JWT secret - use environment variables",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for path traversal vulnerabilities
		if strings.Contains(line, "fs.") && (strings.Contains(line, "req.") || strings.Contains(line, "params.") || strings.Contains(line, "query.")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential path traversal - validate and sanitize file paths from user input",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for unsafe regex (ReDoS)
		if strings.Contains(line, "new RegExp(") && !strings.Contains(line, "new RegExp(\"") && !strings.Contains(line, "new RegExp('") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Non-literal RegExp - potential ReDoS vulnerability with user input",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for object injection/prototype pollution
		if strings.Contains(line, "Object.assign(") && strings.Contains(line, "req.") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Object.assign with user input - potential prototype pollution",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for non-null assertion (!)
		if strings.Contains(line, "!.") || strings.Contains(line, "!)") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Non-null assertion (!) used - consider proper null checking",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for SQL query string concatenation
		if (strings.Contains(line, "query(") || strings.Contains(line, "execute(")) && (strings.Contains(line, "+") || strings.Contains(line, "${")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential SQL injection - use parameterized queries instead of string concatenation",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for non-literal require (potential code injection)
		if strings.Contains(line, "require(") && !strings.Contains(line, "require(\"") && !strings.Contains(line, "require('") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Non-literal require() - potential arbitrary code execution",
				File:     file,
				Line:     i + 1,
			})
		}
	}
}
