package review

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// SecurityPattern defines a pattern to check with exclusions
type SecurityPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Exclusions  []*regexp.Regexp
	Message     string
	Severity    string
}

// Files to always skip for security scanning
var securityIgnoreFiles = []string{
	"package-lock.json",
	"yarn.lock",
	"pnpm-lock.yaml",
	"Gemfile.lock",
	"poetry.lock",
	"go.sum",
	"Cargo.lock",
	"composer.lock",
}

// Files/patterns to skip
var securityIgnorePatterns = []string{
	"*.min.js",
	"*.min.css",
	"*.map",
	"*.snap",
	"__snapshots__/*",
	"*.generated.*",
	"vendor/*",
	"node_modules/*",
}

// GetSecurityPatterns returns the patterns to check for security issues
func GetSecurityPatterns() []SecurityPattern {
	return []SecurityPattern{
		{
			Name: "hardcoded_password",
			// Match: password = "value" or password: "value" with actual content (8+ chars)
			Pattern: regexp.MustCompile(`(?i)password\s*[:=]\s*["']([^"']{8,})["']`),
			Exclusions: []*regexp.Regexp{
				regexp.MustCompile(`(?i)type\s*[:=]\s*["']password["']`),           // HTML input type
				regexp.MustCompile(`(?i)autocomplete\s*[:=]\s*["'].*password.*["']`), // autocomplete attr
				regexp.MustCompile(`(?i)password\s*[:=]\s*["']["']`),                // empty string
				regexp.MustCompile(`(?i)placeholder.*password`),                     // placeholder text
				regexp.MustCompile(`(?i)label.*password`),                           // label text
				regexp.MustCompile(`(?i)message.*password`),                         // error messages
				regexp.MustCompile(`(?i)name\s*[:=]\s*["'].*password.*["']`),        // form field names
				regexp.MustCompile(`(?i)required.*password`),                        // validation rules
				regexp.MustCompile(`(?i)password.*required`),                        // validation rules
			},
			Message:  "Potential hardcoded password detected",
			Severity: "high",
		},
		{
			Name: "hardcoded_api_key",
			// Match: api_key = "value" with actual key-like content
			Pattern: regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*["']([A-Za-z0-9_\-]{16,})["']`),
			Exclusions: []*regexp.Regexp{
				regexp.MustCompile(`(?i)process\.env`),                     // env var reference
				regexp.MustCompile(`(?i)ENV\[`),                            // Ruby env
				regexp.MustCompile(`(?i)os\.environ`),                      // Python env
				regexp.MustCompile(`(?i)getenv`),                           // getenv calls
				regexp.MustCompile(`(?i)api_key.*\(\)`),                    // method calls
				regexp.MustCompile(`(?i)def\s+api_key`),                    // method definitions
				regexp.MustCompile(`(?i)function\s+api_key`),               // function definitions
				regexp.MustCompile(`(?i)api_key_authorized`),               // method names
			},
			Message:  "Potential hardcoded API key detected",
			Severity: "high",
		},
		{
			Name: "hardcoded_secret",
			// Match: secret = "value" with actual content
			Pattern: regexp.MustCompile(`(?i)(secret|secret_key|client_secret)\s*[:=]\s*["']([A-Za-z0-9_\-]{16,})["']`),
			Exclusions: []*regexp.Regexp{
				regexp.MustCompile(`(?i)process\.env`),
				regexp.MustCompile(`(?i)ENV\[`),
				regexp.MustCompile(`(?i)os\.environ`),
				regexp.MustCompile(`(?i)getenv`),
				regexp.MustCompile(`(?i)\{\{.*secret.*\}\}`),              // template vars
				regexp.MustCompile(`(?i)\$\{.*secret.*\}`),                // interpolation
				regexp.MustCompile(`(?i)otp_secret`),                      // OTP display (from var)
				regexp.MustCompile(`(?i)secret.*data\[`),                  // accessing data
				regexp.MustCompile(`(?i)data\..*secret`),                  // accessing data
			},
			Message:  "Potential hardcoded secret detected",
			Severity: "high",
		},
		{
			Name: "private_key",
			// Match: actual private key content
			Pattern: regexp.MustCompile(`-----BEGIN\s+(RSA|EC|DSA|OPENSSH|PGP)?\s*PRIVATE KEY-----`),
			Exclusions: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\.example`),
				regexp.MustCompile(`(?i)template`),
				regexp.MustCompile(`(?i)\.sample`),
			},
			Message:  "Private key detected in code",
			Severity: "high",
		},
		{
			Name: "aws_credentials",
			// Match: AWS access key ID pattern (starts with AKIA, ABIA, ACCA, ASIA)
			Pattern: regexp.MustCompile(`(A3T[A-Z0-9]|AKIA|ABIA|ACCA|ASIA)[A-Z0-9]{16}`),
			Exclusions: []*regexp.Regexp{
				regexp.MustCompile(`(?i)example`),
				regexp.MustCompile(`(?i)placeholder`),
				regexp.MustCompile(`(?i)your.?access.?key`),
			},
			Message:  "AWS access key detected",
			Severity: "high",
		},
		{
			Name: "generic_token",
			// Match: token = "value" with JWT-like or long random string
			Pattern: regexp.MustCompile(`(?i)(auth_token|access_token|bearer)\s*[:=]\s*["']([A-Za-z0-9_\-\.]{32,})["']`),
			Exclusions: []*regexp.Regexp{
				regexp.MustCompile(`(?i)process\.env`),
				regexp.MustCompile(`(?i)ENV\[`),
				regexp.MustCompile(`(?i)getenv`),
				regexp.MustCompile(`(?i)localStorage`),
				regexp.MustCompile(`(?i)sessionStorage`),
				regexp.MustCompile(`(?i)cookie`),
			},
			Message:  "Potential hardcoded token detected",
			Severity: "high",
		},
	}
}

// shouldSkipFileForSecurity checks if a file should be skipped for security scanning
func (a *Analyzer) shouldSkipFileForSecurity(filePath string) bool {
	baseName := filepath.Base(filePath)
	
	// Check exact matches
	for _, ignore := range securityIgnoreFiles {
		if baseName == ignore {
			if a.verbose {
				color.Blue("[INFO] Skipping security scan for lockfile: %s", filePath)
			}
			return true
		}
	}
	
	// Check patterns
	for _, pattern := range securityIgnorePatterns {
		if matched, _ := filepath.Match(pattern, filePath); matched {
			if a.verbose {
				color.Blue("[INFO] Skipping security scan for pattern match: %s", filePath)
			}
			return true
		}
		if matched, _ := filepath.Match(pattern, baseName); matched {
			if a.verbose {
				color.Blue("[INFO] Skipping security scan for pattern match: %s", filePath)
			}
			return true
		}
	}
	
	return false
}

// getChangedLines returns only the added/modified lines from a file in the diff
func (a *Analyzer) getChangedLines(targetBranch, filePath string) ([]struct {
	LineNum int
	Content string
}, error) {
	// Get diff for specific file showing only added lines
	cmd := exec.Command("git", "diff", "-U0", 
		"--diff-filter=AM",  // Added or Modified
		"origin/"+targetBranch+"..HEAD",
		"--", filePath)
	cmd.Dir = a.repoPath
	
	output, err := cmd.Output()
	if err != nil {
		// Fallback: try without origin
		cmd = exec.Command("git", "diff", "-U0", 
			"--diff-filter=AM",
			targetBranch+"..HEAD",
			"--", filePath)
		cmd.Dir = a.repoPath
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}
	}
	
	var changedLines []struct {
		LineNum int
		Content string
	}
	
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	currentLine := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Parse @@ -X,Y +A,B @@ to get line numbers
		if strings.HasPrefix(line, "@@") {
			// Extract the +A part (new file line number)
			parts := strings.Split(line, "+")
			if len(parts) >= 2 {
				numPart := strings.Split(parts[1], ",")[0]
				numPart = strings.Split(numPart, " ")[0]
				var startLine int
				if _, err := fmt.Sscanf(numPart, "%d", &startLine); err == nil {
					currentLine = startLine - 1 // Will be incremented on first line
				}
			}
			continue
		}
		
		// Only process added lines (starting with +, but not +++)
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			currentLine++
			content := strings.TrimPrefix(line, "+")
			changedLines = append(changedLines, struct {
				LineNum int
				Content string
			}{
				LineNum: currentLine,
				Content: content,
			})
		} else if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "\\") {
			// Context line or header, increment counter
			if len(line) > 0 && line[0] != '-' && line[0] != '\\' && line[0] != '@' && line[0] != 'd' {
				currentLine++
			}
		}
	}
	
	return changedLines, nil
}

// RunSecurityChecksV2 runs improved security checks on changed lines only
func (a *Analyzer) RunSecurityChecksV2(report *Report, targetBranch string) {
	if a.verbose {
		color.Blue("[INFO] Running improved security checks (changed lines only)")
	}
	
	patterns := GetSecurityPatterns()
	
	for _, file := range report.ChangedFiles {
		// Skip files that shouldn't be security scanned
		if a.shouldSkipFileForSecurity(file) {
			continue
		}
		
		if a.verbose {
			color.Blue("[INFO] Security scanning changed lines in: %s", file)
		}
		
		// Get only changed lines
		changedLines, err := a.getChangedLines(targetBranch, file)
		if err != nil {
			if a.verbose {
				color.Yellow("[WARN] Could not get changed lines for %s: %v", file, err)
			}
			continue
		}
		
		if a.verbose {
			color.Blue("[INFO] Found %d changed lines in %s", len(changedLines), file)
		}
		
		// Check each changed line against patterns
		for _, line := range changedLines {
			for _, sp := range patterns {
				// Check if line matches the pattern
				if !sp.Pattern.MatchString(line.Content) {
					continue
				}
				
				// Check exclusions
				excluded := false
				for _, exc := range sp.Exclusions {
					if exc.MatchString(line.Content) {
						excluded = true
						if a.verbose {
							color.Blue("[INFO] Line excluded by pattern: %s", exc.String())
						}
						break
					}
				}
				
				if !excluded {
					report.AddIssue(Issue{
						Type:     "security",
						Severity: sp.Severity,
						Message:  sp.Message,
						File:     file,
						Line:     line.LineNum,
					})
					if a.verbose {
						color.Yellow("[WARN] Security issue found: %s at %s:%d", sp.Message, file, line.LineNum)
					}
				}
			}
		}
	}
	
	if a.verbose {
		color.Blue("[INFO] Done running improved security checks")
	}
}
