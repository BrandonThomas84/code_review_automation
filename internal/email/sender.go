package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/BrandonThomas84/code-review-automation/internal/review"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type Sender struct {
	config Config
}

func NewSender(config Config) *Sender {
	return &Sender{config: config}
}

// NewSenderFromEnv creates a Sender with configuration from environment variables
func NewSenderFromEnv() *Sender {
	return &Sender{config: Config{}}
}

// getEnvWithFallback tries the primary env var first, then falls back to the secondary
func getEnvWithFallback(primary, fallback string) string {
	if val := os.Getenv(primary); val != "" {
		return val
	}
	if fallback != "" {
		return os.Getenv(fallback)
	}
	return ""
}

// SendReport sends a formatted email report
func (s *Sender) SendReport(report *review.Report, toEmail string) error {
	return s.SendReportWithContext(report, toEmail, "", "", 0, "")
}

// SendReportWithContext sends a formatted email report with optional context
func (s *Sender) SendReportWithContext(report *review.Report, toEmail, repoName, branchName string, prNumber int, prTitle string) error {
	// Get config from environment if not provided (AUTOREVIEW_ prefixed for GitHub secrets)
	if s.config.SMTPHost == "" {
		s.config.SMTPHost = getEnvWithFallback("AUTOREVIEW_SMTP_HOST", "SMTP_HOST")
	}
	if s.config.SMTPPort == 0 {
		s.config.SMTPPort = 587 // Default SMTP port
	}
	if s.config.SMTPUser == "" {
		s.config.SMTPUser = getEnvWithFallback("AUTOREVIEW_SMTP_USER", "SMTP_USER")
	}
	if s.config.SMTPPassword == "" {
		s.config.SMTPPassword = getEnvWithFallback("AUTOREVIEW_SMTP_PASSWORD", "SMTP_PASSWORD")
	}
	if s.config.FromEmail == "" {
		s.config.FromEmail = getEnvWithFallback("AUTOREVIEW_FROM_EMAIL", "FROM_EMAIL")
	}
	if s.config.FromName == "" {
		s.config.FromName = getEnvWithFallback("AUTOREVIEW_FROM_NAME", "")
		if s.config.FromName == "" {
			s.config.FromName = "AutoReview Bot"
		}
	}

	if s.config.SMTPHost == "" || s.config.SMTPUser == "" {
		return fmt.Errorf("SMTP configuration not provided")
	}

	// Use the new formatter
	formatter := NewFormatter().
		WithRepo(repoName).
		WithBranch(branchName).
		WithPR(prNumber, prTitle)

	subject := formatter.FormatSubject(report)
	body := formatter.FormatHTML(report)

	// Send email
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		s.config.FromName, s.config.FromEmail, toEmail, subject, body)

	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{toEmail}, []byte(msg))
}
