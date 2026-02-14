package review

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Analyzer struct {
	repoPath       string
	ignorePatterns []string
	verbose        bool
	targetBranch   string // Store for use in security checks
}

func NewAnalyzer(repoPath string, verbose bool) *Analyzer {
	analyzer := &Analyzer{
		repoPath:       repoPath,
		ignorePatterns: []string{},
		verbose:        verbose,
	}
	// Load ignore patterns from .autoreview-ignore file
	analyzer.loadIgnorePatterns()
	return analyzer
}

// loadIgnorePatterns reads the .autoreview-ignore file and loads patterns
func (a *Analyzer) loadIgnorePatterns() {
	if a.verbose {
		color.Blue("[INFO] Loading ignore patterns...")
	}

	ignoreFilePath := filepath.Join(a.repoPath, ".autoreview-ignore")
	content, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		// File doesn't exist or can't be read, which is fine
		return
	}

	if a.verbose {
		color.Blue("[INFO] Found ignore file")
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
	if a.verbose {
		color.Blue("[INFO] Checking if file should be ignored: %s", filePath)
	}

	for _, pattern := range a.ignorePatterns {
		// Check for exact match
		if filePath == pattern {
			if a.verbose {
				color.Blue("[INFO] File matches ignore pattern: %s", pattern)
			}
			return true
		}
		// Check if pattern matches using filepath.Match (supports wildcards)
		if matched, err := filepath.Match(pattern, filePath); err == nil && matched {
			if a.verbose {
				color.Blue("[INFO] File matches ignore pattern: %s", pattern)
			}
			return true
		}
		// Check if the file is within an ignored directory
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(filePath, dirPattern+"/") {
				if a.verbose {
					color.Blue("[INFO] File is within ignored directory:", pattern)
				}
				return true
			}
		}
	}

	if a.verbose {
		color.Blue("[INFO] File should NOT be ignored")
	}

	return false
}

func (a *Analyzer) GenerateReport(targetBranch string, fullScan bool) (*Report, error) {
	if a.verbose {
		color.Blue("[INFO] Generating report...")
	}

	// Store target branch for use in security checks
	a.targetBranch = targetBranch

	report := NewReport()

	if fullScan {
		if a.verbose {
			color.Blue("[INFO] Full scan requested")
		}

		if err := a.analyzeFullCodebase(report); err != nil {
			return nil, fmt.Errorf("full codebase analysis failed: %w", err)
		}
		// Full scan uses old security checks (scans whole files)
		a.runSecurityChecks(report)
	} else {
		if a.verbose {
			color.Blue("[INFO] Analyzing git diff")
		}

		if err := a.analyzeGitDiff(targetBranch, report); err != nil {
			return nil, fmt.Errorf("git diff analysis failed: %w", err)
		}
		// Diff mode uses improved security checks (changed lines only)
		a.RunSecurityChecksV2(report, targetBranch)
	}

	// Run quality checks
	a.runQualityChecks(report)

	return report, nil
}

func (a *Analyzer) analyzeGitDiff(targetBranch string, report *Report) error {
	// Fetch the target branch
	cmd := exec.Command("git", "fetch", "origin", targetBranch)
	cmd.Dir = a.repoPath
	cmd.Run() // Ignore error, branch might be local

	if a.verbose {
		color.Blue("[INFO] Getting changed files...")
	}

	// Get changed files
	cmd = exec.Command("git", "diff", "--name-only", fmt.Sprintf("origin/%s..HEAD", targetBranch))

	if a.verbose {
		color.Blue("[INFO] Git command: %s\n", cmd.String())
	}

	cmd.Dir = a.repoPath
	output, err := cmd.Output()
	if err != nil {
		// Fallback without origin
		cmd = exec.Command("git", "diff", "--name-only", fmt.Sprintf("%s..HEAD", targetBranch))
		cmd.Dir = a.repoPath
		output, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get changed files: %w", err)
		}
	}

	if a.verbose {
		color.Blue("[INFO] Found changed files")
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, f := range files {
		if f != "" && !a.shouldIgnoreFile(f) {
			report.ChangedFiles = append(report.ChangedFiles, f)
		}
	}

	if a.verbose {
		color.Blue("[INFO] Done analyzing git diff")
	}

	return nil
}

func (a *Analyzer) analyzeFullCodebase(report *Report) error {
	codeExtensions := []string{".py", ".js", ".ts", ".jsx", ".tsx", ".dart", ".rb", ".php", ".java", ".kt"}

	if a.verbose {
		color.Blue("[INFO] Analyzing full codebase")
		color.Blue("[INFO] Searching for files with extensions:", codeExtensions)
	}

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

	if a.verbose {
		color.Blue("[INFO] Done analyzing full codebase")
	}

	return nil
}

func (a *Analyzer) runSecurityChecks(report *Report) {
	if a.verbose {
		color.Blue("[INFO] Running security checks")
	}

	// Check for common security issues
	patterns := map[string]string{
		"password":    "Hardcoded password detected",
		"api_key":     "Hardcoded API key detected",
		"secret":      "Hardcoded secret detected",
		"private_key": "Private key in code",
		"aws_access":  "AWS credentials in code",
	}

	if a.verbose {
		color.Blue("[INFO] Checking for security issues...")
	}

	for _, file := range report.ChangedFiles {
		if a.verbose {
			color.Blue("[INFO] Checking file for security issues: %s", file)
		}

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

		if a.verbose {
			color.Blue("[INFO] Done checking for security issues in file: %s", file)
		}
	}

	if a.verbose {
		color.Blue("[INFO] Done running security checks")
	}
}

func (a *Analyzer) runQualityChecks(report *Report) {
	if a.verbose {
		color.Blue("[INFO] Running quality checks")
	}

	// Check for code quality issues
	for _, file := range report.ChangedFiles {
		switch {
		case strings.HasSuffix(file, ".py"):
			a.checkPythonQuality(file, report)
		case strings.HasSuffix(file, ".js"), strings.HasSuffix(file, ".jsx"):
			a.checkJavaScriptQuality(file, report)
		case strings.HasSuffix(file, ".ts"), strings.HasSuffix(file, ".tsx"):
			a.checkTypeScriptQuality(file, report)
		case strings.HasSuffix(file, ".rb"):
			a.checkRubyQuality(file, report)
		case strings.HasSuffix(file, ".dart"):
			a.checkDartQuality(file, report)
		case strings.HasSuffix(file, ".php"):
			a.checkPHPQuality(file, report)
		case strings.HasSuffix(file, ".java"), strings.HasSuffix(file, ".kt"):
			a.checkJavaKotlinQuality(file, report)
		}
	}
}
