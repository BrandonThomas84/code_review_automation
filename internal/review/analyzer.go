package review

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Analyzer struct {
	repoPath       string
	ignorePatterns []string
}

func NewAnalyzer(repoPath string) *Analyzer {
	analyzer := &Analyzer{
		repoPath:       repoPath,
		ignorePatterns: []string{},
	}
	// Load ignore patterns from .autoreview-ignore file
	analyzer.loadIgnorePatterns()
	return analyzer
}

// loadIgnorePatterns reads the .autoreview-ignore file and loads patterns
func (a *Analyzer) loadIgnorePatterns() {
	ignoreFilePath := filepath.Join(a.repoPath, ".autoreview-ignore")
	content, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		// File doesn't exist or can't be read, which is fine
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		// Trim whitespace
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			a.ignorePatterns = append(a.ignorePatterns, line)
		}
	}
}

// shouldIgnoreFile checks if a file matches any ignore patterns
func (a *Analyzer) shouldIgnoreFile(filePath string) bool {
	for _, pattern := range a.ignorePatterns {
		// Check for exact match
		if filePath == pattern {
			return true
		}
		// Check if pattern matches using filepath.Match (supports wildcards)
		if matched, err := filepath.Match(pattern, filePath); err == nil && matched {
			return true
		}
		// Check if the file is within an ignored directory
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(filePath, dirPattern+"/") {
				return true
			}
		}
	}
	return false
}

func (a *Analyzer) GenerateReport(targetBranch string, fullScan bool) (*Report, error) {
	report := NewReport()

	if fullScan {
		if err := a.analyzeFullCodebase(report); err != nil {
			return nil, fmt.Errorf("full codebase analysis failed: %w", err)
		}
	} else {
		if err := a.analyzeGitDiff(targetBranch, report); err != nil {
			return nil, fmt.Errorf("git diff analysis failed: %w", err)
		}
	}

	// Run checks
	a.runSecurityChecks(report)
	a.runQualityChecks(report)

	return report, nil
}

func (a *Analyzer) analyzeGitDiff(targetBranch string, report *Report) error {
	// Fetch the target branch
	cmd := exec.Command("git", "fetch", "origin", targetBranch)
	cmd.Dir = a.repoPath
	cmd.Run() // Ignore error, branch might be local

	// Get changed files
	cmd = exec.Command("git", "diff", "--name-only", fmt.Sprintf("origin/%s...HEAD", targetBranch))
	cmd.Dir = a.repoPath
	output, err := cmd.Output()
	if err != nil {
		// Fallback without origin
		cmd = exec.Command("git", "diff", "--name-only", fmt.Sprintf("%s...HEAD", targetBranch))
		cmd.Dir = a.repoPath
		output, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get changed files: %w", err)
		}
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, f := range files {
		if f != "" && !a.shouldIgnoreFile(f) {
			report.ChangedFiles = append(report.ChangedFiles, f)
		}
	}

	return nil
}

func (a *Analyzer) analyzeFullCodebase(report *Report) error {
	codeExtensions := []string{".py", ".js", ".ts", ".jsx", ".tsx", ".dart", ".rb", ".php", ".java", ".kt"}

	for _, ext := range codeExtensions {
		cmd := exec.Command("find", ".", "-name", fmt.Sprintf("*%s", ext), "-type", "f")
		cmd.Dir = a.repoPath
		output, err := cmd.Output()
		if err == nil {
			files := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, f := range files {
				if f != "" && f != "." && !a.shouldIgnoreFile(f) {
					report.ChangedFiles = append(report.ChangedFiles, f)
				}
			}
		}
	}

	return nil
}

func (a *Analyzer) runSecurityChecks(report *Report) {
	// Check for common security issues
	patterns := map[string]string{
		"password":    "Hardcoded password detected",
		"api_key":     "Hardcoded API key detected",
		"secret":      "Hardcoded secret detected",
		"private_key": "Private key in code",
		"aws_access":  "AWS credentials in code",
	}

	for _, file := range report.ChangedFiles {
		filePath := filepath.Join(a.repoPath, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		contentStr := strings.ToLower(string(content))
		for pattern, message := range patterns {
			if strings.Contains(contentStr, pattern) {
				report.AddIssue(Issue{
					Type:     "security",
					Severity: "high",
					Message:  message,
					File:     file,
				})
			}
		}
	}
}

func (a *Analyzer) runQualityChecks(report *Report) {
	// Check for code quality issues
	for _, file := range report.ChangedFiles {
		if strings.HasSuffix(file, ".py") {
			a.checkPythonQuality(file, report)
		} else if strings.HasSuffix(file, ".js") || strings.HasSuffix(file, ".ts") {
			a.checkJavaScriptQuality(file, report)
		}
	}
}

func (a *Analyzer) checkPythonQuality(file string, report *Report) {
	filePath := filepath.Join(a.repoPath, file)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for i, line := range lines {
		if len(line) > 120 {
			report.AddIssue(Issue{
				Type:     "quality",
				Severity: "low",
				Message:  "Line too long (>120 characters)",
				File:     file,
				Line:     i + 1,
			})
		}
	}
}

func (a *Analyzer) checkJavaScriptQuality(file string, report *Report) {
	// Similar checks for JavaScript
}
