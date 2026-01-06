package review

import (
	"os"
	"path/filepath"
	"strings"
)

// checkRubyQuality analyzes Ruby files for quality and security issues
func (a *Analyzer) checkRubyQuality(file string, report *Report) {
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

		// Line length check (Ruby style guide recommends 80, but 120 is common)
		if len(line) > 120 {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Line too long (>120 characters)",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for puts/p debug statements
		if strings.HasPrefix(trimmed, "puts ") || strings.HasPrefix(trimmed, "p ") || strings.HasPrefix(trimmed, "pp ") {
			// Avoid false positives for method definitions
			if !strings.Contains(line, "def ") {
				report.AddIssue(Issue{
					Type:     "quality",
					Severity: "low",
					Message:  "Debug output (puts/p/pp) found - remove before production",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// Check for binding.pry or byebug (debugger)
		if strings.Contains(line, "binding.pry") || strings.Contains(line, "byebug") || strings.Contains(line, "debugger") {
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
		if strings.Contains(line, "eval(") || strings.Contains(line, "instance_eval") || strings.Contains(line, "class_eval") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "eval() usage detected - potential code injection vulnerability",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for system/exec calls (Command Injection)
		if strings.Contains(line, "system(") || strings.Contains(line, "exec(") || strings.Contains(line, "`") || strings.Contains(line, "%x(") || strings.Contains(line, "Open3.") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Shell command execution detected - ensure input is sanitized to prevent command injection",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for SQL injection (raw SQL with interpolation)
		if strings.Contains(line, ".where(\"") || strings.Contains(line, ".find_by_sql(") || strings.Contains(line, ".execute(") {
			if strings.Contains(line, "#{") {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  "Potential SQL injection - use parameterized queries instead of string interpolation",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// SECURITY: Check for mass assignment vulnerabilities
		if strings.Contains(line, ".update_attributes(") || strings.Contains(line, ".update(params") || strings.Contains(line, ".create(params") || strings.Contains(line, ".new(params") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential mass assignment vulnerability - use strong parameters",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for XSS vulnerabilities (raw HTML output)
		if strings.Contains(line, ".html_safe") || strings.Contains(line, "raw(") || strings.Contains(line, "<%==") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential XSS vulnerability - html_safe/raw bypasses HTML escaping",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for unsafe YAML loading
		if strings.Contains(line, "YAML.load(") && !strings.Contains(line, "YAML.safe_load(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Unsafe YAML.load - use YAML.safe_load to prevent code execution",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for unsafe deserialization
		if strings.Contains(line, "Marshal.load(") || strings.Contains(line, "Marshal.restore(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Unsafe deserialization with Marshal - can lead to remote code execution",
				File:     file,
				Line:     i + 1,
			})
		}

		// Rescue without specific exception
		if strings.Contains(line, "rescue StandardError") || strings.Contains(line, "rescue =>") {
			report.AddIssue(Issue{
				Type:     "error_handling",
				Severity: "medium",
				Message:  "Generic rescue clause",
				File:     file,
				Line:     i + 1,
			})
		}

		// Empty rescue blocks
		if strings.Contains(line, "rescue") && strings.Contains(line, "end") {
			report.AddIssue(Issue{
				Type:     "error_handling",
				Severity: "medium",
				Message:  "Empty rescue block",
				File:     file,
				Line:     i + 1,
			})
		}
	}

	// Continue with more security checks in a helper function
	a.checkRubySecurityExtended(file, contentStr, lines, report)
}

// checkRubySecurityExtended contains additional Ruby security checks
func (a *Analyzer) checkRubySecurityExtended(file string, contentStr string, lines []string, report *Report) {
	for i, line := range lines {
		lineLower := strings.ToLower(line)

		// SECURITY: Check for open redirect vulnerabilities
		if strings.Contains(line, "redirect_to") && (strings.Contains(line, "params[") || strings.Contains(line, "request.")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Potential open redirect - validate redirect URLs",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for file access with user input
		if (strings.Contains(line, "File.read(") || strings.Contains(line, "File.open(") || strings.Contains(line, "IO.read(")) && strings.Contains(line, "params[") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential path traversal - validate file paths from user input",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for send with user input (dangerous send)
		if strings.Contains(line, ".send(") && (strings.Contains(line, "params[") || strings.Contains(line, "#{")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Dangerous send with user input - can call arbitrary methods",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for constantize with user input
		if strings.Contains(line, ".constantize") && (strings.Contains(line, "params[") || strings.Contains(line, "#{")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Dangerous constantize with user input - can instantiate arbitrary classes",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for render with user input (dynamic render path)
		if strings.Contains(line, "render") && strings.Contains(line, "params[") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Dynamic render path with user input - potential information disclosure",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for weak cryptography
		if strings.Contains(line, "MD5.") || strings.Contains(line, "Digest::MD5") || strings.Contains(line, "SHA1.") || strings.Contains(line, "Digest::SHA1") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Weak hash algorithm (MD5/SHA1) - use SHA256 or stronger",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for SSL verification bypass
		if strings.Contains(line, "verify_mode") && strings.Contains(line, "VERIFY_NONE") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "SSL verification disabled - vulnerable to man-in-the-middle attacks",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for session manipulation
		if strings.Contains(line, "session[") && strings.Contains(line, "params[") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Session manipulation with user input - validate before storing",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for unscoped find
		if strings.Contains(line, ".find(params[") || strings.Contains(line, ".find_by(") {
			if !strings.Contains(contentStr, "current_user") && !strings.Contains(line, "current_user") {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "medium",
					Message:  "Unscoped find - consider scoping to current user to prevent unauthorized access",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// SECURITY: Check for basic authentication credentials
		if strings.Contains(lineLower, "basic_auth") || (strings.Contains(lineLower, "authorization") && strings.Contains(lineLower, "basic")) {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Basic authentication detected - ensure credentials are not hardcoded",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for CSRF protection disabled
		if strings.Contains(line, "skip_before_action :verify_authenticity_token") || strings.Contains(line, "protect_from_forgery except:") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "CSRF protection disabled - ensure this is intentional and properly secured",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Missing strong parameters
		if strings.Contains(line, ".params[") && !strings.Contains(line, ".permit(") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Open parameters detected - use strong parameters to whitelist allowed attributes",
				File:     file,
				Line:     i + 1,
			})
		}

		// N+1 query patterns
		if strings.Contains(line, ".each") && strings.Contains(line, ".find") {
			report.AddIssue(Issue{
				Type:     "performance",
				Severity: "high",
				Message:  "Potential N+1 query detected",
				File:     file,
				Line:     i + 1,
			})
		}

		// Missing validations in models
		if strings.Contains(file, "model") && strings.Contains(line, "class") && !strings.Contains(line, "validates") {
			report.AddIssue(Issue{
				Type:     "rails_structure",
				Severity: "medium",
				Message:  "Model without validations",
				File:     file,
				Line:     i + 1,
			})
		}

		// Callback hell
		callbackCount := strings.Count(contentStr, "before_") + strings.Count(contentStr, "after_") + strings.Count(contentStr, "around_")
		if callbackCount > 5 {
			report.AddIssue(Issue{
				Type:     "rails_structure",
				Severity: "medium",
				Message:  "Too many callbacks detected",
				File:     file,
				Line:     i + 1,
			})
		}

		// Inefficient queries in loops
		if strings.Contains(line, ".each") && (strings.Contains(line, ".find") || strings.Contains(line, ".where") || strings.Contains(line, ".create") || strings.Contains(line, ".update")) {
			report.AddIssue(Issue{
				Type:     "performance",
				Severity: "medium",
				Message:  "Database query inside loop",
				File:     file,
				Line:     i + 1,
			})
		}

		// Inefficient string concatenation
		if strings.Contains(line, "+=") && (strings.Contains(line, "\"") || strings.Contains(line, "'")) {
			report.AddIssue(Issue{
				Type:     "performance",
				Severity: "low",
				Message:  "String concatenation with +=",
				File:     file,
				Line:     i + 1,
			})
		}
	}
}
