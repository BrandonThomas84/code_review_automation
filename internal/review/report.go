package review

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Issue struct {
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	File     string `json:"file"`
	Line     int    `json:"line,omitempty"`
}

type Report struct {
	Timestamp    time.Time `json:"timestamp"`
	ChangedFiles []string  `json:"changed_files"`
	Issues       []Issue   `json:"issues"`
	Summary      Summary   `json:"summary"`
}

type Summary struct {
	TotalFiles     int `json:"total_files"`
	TotalIssues    int `json:"total_issues"`
	HighSeverity   int `json:"high_severity"`
	MediumSeverity int `json:"medium_severity"`
	LowSeverity    int `json:"low_severity"`
}

func NewReport() *Report {
	return &Report{
		Timestamp:    time.Now(),
		ChangedFiles: []string{},
		Issues:       []Issue{},
	}
}

func (r *Report) AddIssue(issue Issue) {
	r.Issues = append(r.Issues, issue)
	r.updateSummary()
}

func (r *Report) updateSummary() {
	r.Summary.TotalFiles = len(r.ChangedFiles)
	r.Summary.TotalIssues = len(r.Issues)
	r.Summary.HighSeverity = 0
	r.Summary.MediumSeverity = 0
	r.Summary.LowSeverity = 0

	for _, issue := range r.Issues {
		switch issue.Severity {
		case "high":
			r.Summary.HighSeverity++
		case "medium":
			r.Summary.MediumSeverity++
		case "low":
			r.Summary.LowSeverity++
		}
	}
}

func (r *Report) PrintReport() {
	// create separator string
	equal_separator := strings.Repeat("=", 60)
	color.Blue("\n" + equal_separator)
	color.Blue("📋 CODE REVIEW SUMMARY")
	color.Blue(equal_separator)
	fmt.Printf("📁 Files changed: %d\n", r.Summary.TotalFiles)
	fmt.Printf("🚨 Total issues: %d\n", r.Summary.TotalIssues)
	color.Red("🔴 High severity: %d\n", r.Summary.HighSeverity)
	color.Yellow("🟡 Medium severity: %d\n", r.Summary.MediumSeverity)
	color.Green("🟢 Low severity: %d\n", r.Summary.LowSeverity)

	if len(r.Issues) > 0 {
		line_separator := strings.Repeat("-", 60)
		fmt.Println("\n" + line_separator)
		fmt.Println("ISSUES FOUND:")
		for i, issue := range r.Issues {
			fmt.Printf("%d. [%s] %s\n", i+1, issue.Severity, issue.Message)
			fmt.Printf("   File: %s", issue.File)
			if issue.Line > 0 {
				fmt.Printf(" (line %d)", issue.Line)
			}
			fmt.Println()
		}
	}
}

func (r *Report) OutputJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r)
}

func (r *Report) SaveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return r.OutputJSON(file)
}
