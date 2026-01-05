package review

import (
	"os"
	"path/filepath"
	"strings"
)

// checkJavaKotlinQuality analyzes Java and Kotlin files for quality and security issues
func (a *Analyzer) checkJavaKotlinQuality(file string, report *Report) {
	filePath := filepath.Join(a.repoPath, file)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	isKotlin := strings.HasSuffix(file, ".kt")

	for i, line := range lines {
		lineLower := strings.ToLower(line)
		trimmed := strings.TrimSpace(line)

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

		// Check for System.out.println (Java) or println (Kotlin)
		if strings.Contains(line, "System.out.println") || strings.Contains(line, "System.err.println") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "System.out.println found - use proper logging instead",
				File:     file,
				Line:     i + 1,
			})
		}

		// Check for e.printStackTrace()
		if strings.Contains(line, ".printStackTrace()") {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "medium",
				Message:  "printStackTrace() found - use proper logging instead",
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

		// Check for empty catch blocks
		if trimmed == "catch" || strings.Contains(line, "catch (") {
			// Look ahead for empty catch block
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine == "}" || nextLine == "{ }" || nextLine == "{}" {
					report.AddIssue(Issue{
						Type:     "quality",
						Severity: "medium",
						Message:  "Empty catch block - handle or log the exception",
						File:     file,
						Line:     i + 1,
					})
				}
			}
		}

		// SECURITY: Check for Runtime.exec
		if strings.Contains(line, "Runtime.getRuntime().exec") || strings.Contains(line, "ProcessBuilder") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Process execution detected - ensure input is sanitized",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for SQL injection
		if strings.Contains(line, "Statement") && strings.Contains(line, "execute") {
			if strings.Contains(line, "+") || strings.Contains(line, "concat") {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  "Potential SQL injection - use PreparedStatement with parameterized queries",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// SECURITY: Check for hardcoded credentials
		if strings.Contains(lineLower, "password") && strings.Contains(line, "=") && strings.Contains(line, "\"") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "Potential hardcoded password - use secure configuration",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for weak cryptography
		if strings.Contains(line, "MD5") || strings.Contains(line, "SHA1") || strings.Contains(line, "DES") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "medium",
				Message:  "Weak cryptographic algorithm - use SHA-256 or stronger",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for disabled SSL verification
		if strings.Contains(line, "TrustAllCerts") || strings.Contains(line, "ALLOW_ALL_HOSTNAME_VERIFIER") {
			report.AddIssue(Issue{
				Type:     "security",
				Severity: "high",
				Message:  "SSL verification disabled - vulnerable to man-in-the-middle attacks",
				File:     file,
				Line:     i + 1,
			})
		}

		// SECURITY: Check for XXE vulnerability
		if strings.Contains(line, "XMLInputFactory") || strings.Contains(line, "DocumentBuilderFactory") {
			if !strings.Contains(contentStr, "setFeature") {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  "XML parser without secure features - potential XXE vulnerability",
					File:     file,
					Line:     i + 1,
				})
			}
		}

		// Kotlin-specific checks
		if isKotlin {
			a.checkKotlinSpecific(file, line, i, report)
		}
	}
}

// checkKotlinSpecific contains Kotlin-specific quality checks
func (a *Analyzer) checkKotlinSpecific(file string, line string, lineNum int, report *Report) {
	// Check for !! (force unwrap) which can cause NullPointerException
	if strings.Contains(line, "!!") {
		report.AddIssue(Issue{
			Type:     "quality",
			Severity: "medium",
			Message:  "Force unwrap (!!) used - consider safe call (?.) or null check",
			File:     file,
			Line:     lineNum + 1,
		})
	}

	// Check for println in Kotlin
	if strings.Contains(line, "println(") && !strings.Contains(line, "System.out") {
		report.AddIssue(Issue{
			Type:     "quality",
			Severity: "low",
			Message:  "println() found - use proper logging instead",
			File:     file,
			Line:     lineNum + 1,
		})
	}
}

