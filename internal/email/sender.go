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

func (s *Sender) SendReport(report *review.Report, toEmail string) error {
	// Get config from environment if not provided
	if s.config.SMTPHost == "" {
		s.config.SMTPHost = os.Getenv("SMTP_HOST")
	}
	if s.config.SMTPUser == "" {
		s.config.SMTPUser = os.Getenv("SMTP_USER")
	}
	if s.config.SMTPPassword == "" {
		s.config.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	}
	if s.config.FromEmail == "" {
		s.config.FromEmail = os.Getenv("FROM_EMAIL")
	}

	if s.config.SMTPHost == "" || s.config.SMTPUser == "" {
		return fmt.Errorf("SMTP configuration not provided")
	}

	// Build email content
	subject := fmt.Sprintf("Code Review Report - %d Issues Found", report.Summary.TotalIssues)
	body := s.buildHTMLBody(report)

	// Send email
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		s.config.FromName, s.config.FromEmail, toEmail, subject, body)

	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{toEmail}, []byte(msg))
}

func (s *Sender) buildHTMLBody(report *review.Report) string {
	html := `
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; }
			.summary { background-color: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0; }
			.high { color: #d32f2f; font-weight: bold; }
			.medium { color: #f57c00; font-weight: bold; }
			.low { color: #388e3c; font-weight: bold; }
			.issue { border-left: 4px solid #ddd; padding: 10px; margin: 10px 0; }
			.issue.high { border-left-color: #d32f2f; }
			.issue.medium { border-left-color: #f57c00; }
			.issue.low { border-left-color: #388e3c; }
		</style>
	</head>
	<body>
		<h2>Code Review Report</h2>
		<div class="summary">
			<p><strong>Files Changed:</strong> ` + fmt.Sprintf("%d", report.Summary.TotalFiles) + `</p>
			<p><strong>Total Issues:</strong> ` + fmt.Sprintf("%d", report.Summary.TotalIssues) + `</p>
			<p><span class="high">🔴 High: ` + fmt.Sprintf("%d", report.Summary.HighSeverity) + `</span> | 
			   <span class="medium">🟡 Medium: ` + fmt.Sprintf("%d", report.Summary.MediumSeverity) + `</span> | 
			   <span class="low">🟢 Low: ` + fmt.Sprintf("%d", report.Summary.LowSeverity) + `</span></p>
		</div>
	`

	if len(report.Issues) > 0 {
		html += `<h3>Issues Found:</h3>`
		for _, issue := range report.Issues {
			html += fmt.Sprintf(`
			<div class="issue %s">
				<strong>[%s]</strong> %s<br>
				<small>File: <code>%s</code>`, issue.Severity, issue.Severity, issue.Message, issue.File)
			if issue.Line > 0 {
				html += fmt.Sprintf(` (line %d)`, issue.Line)
			}
			html += `</small></div>`
		}
	} else {
		html += `<p style="color: #388e3c;"><strong>✅ No issues found!</strong></p>`
	}

	html += `
	</body>
	</html>
	`
	return html
}

