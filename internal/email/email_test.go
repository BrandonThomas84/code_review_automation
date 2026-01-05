package email

import (
	"os"
	"strings"
	"testing"

	"github.com/BrandonThomas84/code-review-automation/internal/review"
)

// ============== Formatter Tests ==============

func TestNewFormatter(t *testing.T) {
	f := NewFormatter()
	if f == nil {
		t.Fatal("NewFormatter returned nil")
	}
	if f.RepoName != "" || f.BranchName != "" || f.PRNumber != 0 {
		t.Error("NewFormatter should initialize with empty values")
	}
}

func TestFormatter_FluentBuilder(t *testing.T) {
	f := NewFormatter().
		WithRepo("test-repo").
		WithBranch("main").
		WithPR(123, "Test PR Title")

	if f.RepoName != "test-repo" {
		t.Errorf("Expected RepoName 'test-repo', got '%s'", f.RepoName)
	}
	if f.BranchName != "main" {
		t.Errorf("Expected BranchName 'main', got '%s'", f.BranchName)
	}
	if f.PRNumber != 123 {
		t.Errorf("Expected PRNumber 123, got %d", f.PRNumber)
	}
	if f.PRTitle != "Test PR Title" {
		t.Errorf("Expected PRTitle 'Test PR Title', got '%s'", f.PRTitle)
	}
}

func TestFormatter_FormatSubject_NoIssues(t *testing.T) {
	f := NewFormatter()
	report := review.NewReport()

	subject := f.FormatSubject(report)
	if !strings.Contains(subject, "✅") {
		t.Error("Expected checkmark emoji for no issues")
	}
	if !strings.Contains(subject, "0 issues") {
		t.Error("Expected '0 issues' in subject")
	}
}

func TestFormatter_FormatSubject_HighSeverity(t *testing.T) {
	f := NewFormatter()
	report := review.NewReport()
	report.AddIssue(review.Issue{Type: "security", Severity: "high", Message: "Test"})

	subject := f.FormatSubject(report)
	if !strings.Contains(subject, "🚨") {
		t.Error("Expected alert emoji for high severity")
	}
}

func TestFormatter_FormatSubject_MediumSeverity(t *testing.T) {
	f := NewFormatter()
	report := review.NewReport()
	report.AddIssue(review.Issue{Type: "quality", Severity: "medium", Message: "Test"})

	subject := f.FormatSubject(report)
	if !strings.Contains(subject, "⚠️") {
		t.Error("Expected warning emoji for medium severity")
	}
}

func TestFormatter_FormatSubject_WithRepo(t *testing.T) {
	f := NewFormatter().WithRepo("my-repo")
	report := review.NewReport()

	subject := f.FormatSubject(report)
	if !strings.Contains(subject, "[my-repo]") {
		t.Error("Expected repo name in subject")
	}
}

func TestFormatter_FormatSubject_WithPR(t *testing.T) {
	f := NewFormatter().WithPR(42, "Fix bug")
	report := review.NewReport()

	subject := f.FormatSubject(report)
	if !strings.Contains(subject, "PR #42") {
		t.Error("Expected PR number in subject")
	}
}

func TestFormatter_FormatHTML_ContainsBasicStructure(t *testing.T) {
	f := NewFormatter()
	report := review.NewReport()

	html := f.FormatHTML(report)

	// Check for HTML structure
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Expected DOCTYPE declaration")
	}
	if !strings.Contains(html, "<html") {
		t.Error("Expected html tag")
	}
	if !strings.Contains(html, "</html>") {
		t.Error("Expected closing html tag")
	}
	if !strings.Contains(html, "Code Review") {
		t.Error("Expected 'Code Review' in HTML")
	}
}

func TestFormatter_FormatHTML_NoIssues(t *testing.T) {
	f := NewFormatter()
	report := review.NewReport()

	html := f.FormatHTML(report)

	if !strings.Contains(html, "No Issues Found") {
		t.Error("Expected 'No Issues Found' message")
	}
	if !strings.Contains(html, "#4caf50") || !strings.Contains(html, "#e8f5e9") {
		t.Error("Expected green color scheme for no issues")
	}
}

func TestFormatter_FormatHTML_WithHighSeverityIssues(t *testing.T) {
	f := NewFormatter()
	report := review.NewReport()
	report.AddIssue(review.Issue{
		Type:     "security",
		Severity: "high",
		Message:  "SQL injection vulnerability",
		File:     "database.py",
		Line:     42,
	})

	html := f.FormatHTML(report)

	if !strings.Contains(html, "#f44336") {
		t.Error("Expected red color for high severity")
	}
	if !strings.Contains(html, "SQL injection vulnerability") {
		t.Error("Expected issue message in HTML")
	}
	if !strings.Contains(html, "database.py") {
		t.Error("Expected file name in HTML")
	}
}

func TestFormatter_FormatHTML_GroupsIssuesBySeverity(t *testing.T) {
	f := NewFormatter()
	report := review.NewReport()
	report.AddIssue(review.Issue{Type: "quality", Severity: "low", Message: "Low issue"})
	report.AddIssue(review.Issue{Type: "security", Severity: "high", Message: "High issue"})
	report.AddIssue(review.Issue{Type: "quality", Severity: "medium", Message: "Medium issue"})

	html := f.FormatHTML(report)

	// High should appear before medium, medium before low
	highIdx := strings.Index(html, "High Severity")
	mediumIdx := strings.Index(html, "Medium Severity")
	lowIdx := strings.Index(html, "Low Severity")

	if highIdx == -1 || mediumIdx == -1 || lowIdx == -1 {
		t.Error("Expected all severity groups in HTML")
	}
	if highIdx > mediumIdx {
		t.Error("High severity should appear before medium")
	}
	if mediumIdx > lowIdx {
		t.Error("Medium severity should appear before low")
	}
}

func TestFormatter_FormatHTML_EscapesHTML(t *testing.T) {
	f := NewFormatter().WithRepo("<script>alert('xss')</script>")
	report := review.NewReport()
	report.AddIssue(review.Issue{
		Type:     "quality",
		Severity: "low",
		Message:  "<b>Bold</b> text",
		File:     "test<>.py",
	})

	html := f.FormatHTML(report)

	if strings.Contains(html, "<script>") {
		t.Error("HTML should be escaped to prevent XSS")
	}
	if strings.Contains(html, "<b>Bold") {
		t.Error("Issue message HTML should be escaped")
	}
}

func TestFormatter_FormatHTML_WithContext(t *testing.T) {
	f := NewFormatter().
		WithRepo("my-app").
		WithBranch("feature/test").
		WithPR(99, "Add new feature")

	report := review.NewReport()
	html := f.FormatHTML(report)

	if !strings.Contains(html, "my-app") {
		t.Error("Expected repo name in HTML")
	}
	if !strings.Contains(html, "feature/test") {
		t.Error("Expected branch name in HTML")
	}
	if !strings.Contains(html, "PR #99") {
		t.Error("Expected PR number in HTML")
	}
}

// ============== Sender Tests ==============

func TestNewSender(t *testing.T) {
	config := Config{
		SMTPHost:     "smtp.test.com",
		SMTPPort:     587,
		SMTPUser:     "user@test.com",
		SMTPPassword: "password",
		FromEmail:    "from@test.com",
		FromName:     "Test Bot",
	}

	sender := NewSender(config)
	if sender == nil {
		t.Fatal("NewSender returned nil")
	}
}

func TestNewSenderFromEnv(t *testing.T) {
	sender := NewSenderFromEnv()
	if sender == nil {
		t.Fatal("NewSenderFromEnv returned nil")
	}
}

func TestGetEnvWithFallback_Primary(t *testing.T) {
	os.Setenv("TEST_PRIMARY", "primary_value")
	os.Setenv("TEST_FALLBACK", "fallback_value")
	defer os.Unsetenv("TEST_PRIMARY")
	defer os.Unsetenv("TEST_FALLBACK")

	result := getEnvWithFallback("TEST_PRIMARY", "TEST_FALLBACK")
	if result != "primary_value" {
		t.Errorf("Expected primary value, got '%s'", result)
	}
}

func TestGetEnvWithFallback_Fallback(t *testing.T) {
	os.Unsetenv("TEST_PRIMARY_MISSING")
	os.Setenv("TEST_FALLBACK", "fallback_value")
	defer os.Unsetenv("TEST_FALLBACK")

	result := getEnvWithFallback("TEST_PRIMARY_MISSING", "TEST_FALLBACK")
	if result != "fallback_value" {
		t.Errorf("Expected fallback value, got '%s'", result)
	}
}

func TestGetEnvWithFallback_NoFallback(t *testing.T) {
	os.Unsetenv("TEST_MISSING")

	result := getEnvWithFallback("TEST_MISSING", "")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestSender_SendReport_MissingConfig(t *testing.T) {
	// Clear any environment variables that might be set
	os.Unsetenv("AUTOREVIEW_SMTP_HOST")
	os.Unsetenv("AUTOREVIEW_SMTP_USER")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_USER")

	sender := NewSender(Config{})
	report := review.NewReport()

	err := sender.SendReport(report, "test@example.com")
	if err == nil {
		t.Error("Expected error when SMTP config is missing")
	}
	if !strings.Contains(err.Error(), "SMTP configuration not provided") {
		t.Errorf("Expected SMTP config error, got: %v", err)
	}
}

func TestSender_EnvVariables_AutoreviewPrefix(t *testing.T) {
	// Set AUTOREVIEW_ prefixed variables
	os.Setenv("AUTOREVIEW_SMTP_HOST", "autoreview.smtp.com")
	os.Setenv("AUTOREVIEW_SMTP_USER", "autoreview@test.com")
	defer os.Unsetenv("AUTOREVIEW_SMTP_HOST")
	defer os.Unsetenv("AUTOREVIEW_SMTP_USER")

	// Verify getEnvWithFallback uses the AUTOREVIEW_ prefix first
	host := getEnvWithFallback("AUTOREVIEW_SMTP_HOST", "SMTP_HOST")
	if host != "autoreview.smtp.com" {
		t.Errorf("Expected AUTOREVIEW_ prefixed host, got '%s'", host)
	}

	user := getEnvWithFallback("AUTOREVIEW_SMTP_USER", "SMTP_USER")
	if user != "autoreview@test.com" {
		t.Errorf("Expected AUTOREVIEW_ prefixed user, got '%s'", user)
	}
}

// ============== Filter Tests ==============

func TestFilterBySeverity(t *testing.T) {
	issues := []review.Issue{
		{Severity: "high", Message: "High 1"},
		{Severity: "medium", Message: "Medium 1"},
		{Severity: "low", Message: "Low 1"},
		{Severity: "high", Message: "High 2"},
		{Severity: "LOW", Message: "Low 2"}, // Test case insensitivity
	}

	high := filterBySeverity(issues, "high")
	if len(high) != 2 {
		t.Errorf("Expected 2 high severity issues, got %d", len(high))
	}

	medium := filterBySeverity(issues, "medium")
	if len(medium) != 1 {
		t.Errorf("Expected 1 medium severity issue, got %d", len(medium))
	}

	low := filterBySeverity(issues, "low")
	if len(low) != 2 {
		t.Errorf("Expected 2 low severity issues, got %d", len(low))
	}
}
