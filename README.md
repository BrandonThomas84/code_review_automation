# AutoReview - Code Review Automation

[![Build and Release](https://github.com/BrandonThomas84/code_review_automation/actions/workflows/build-release.yml/badge.svg)](https://github.com/BrandonThomas84/code_review_automation/actions/workflows/build-release.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**AutoReview** is a standalone, multi-language code review automation tool that integrates with GitHub Actions to automatically analyze pull requests for security vulnerabilities, code quality issues, and best practice violations.

## ✨ Features

- **Multi-Language Support**: Python, JavaScript, TypeScript, Ruby, Dart/Flutter, PHP, Java, and Kotlin
- **Security Analysis**: Detects SQL injection, XSS, eval usage, hardcoded credentials, and more
- **Code Quality Checks**: Finds debug statements, TODO comments, overly long lines, and anti-patterns
- **GitHub Actions Integration**: Automatic PR comments with detailed results
- **Email Notifications**: Optional HTML-formatted email reports via SMTP
- **Cross-Platform Binaries**: Pre-built for Linux, macOS, and Windows
- **Zero Dependencies**: Single binary, no runtime requirements

## 🚀 Quick Start

### Option 1: Download Pre-Built Binary

Download the latest release for your platform:

```bash
# Linux (amd64)
curl -L https://github.com/BrandonThomas84/code_review_automation/releases/latest/download/code-review-linux-amd64 -o code-review
chmod +x code-review

# macOS (Apple Silicon)
curl -L https://github.com/BrandonThomas84/code_review_automation/releases/latest/download/code-review-darwin-arm64 -o code-review
chmod +x code-review

# macOS (Intel)
curl -L https://github.com/BrandonThomas84/code_review_automation/releases/latest/download/code-review-darwin-amd64 -o code-review
chmod +x code-review
```

### Option 2: Build from Source

```bash
git clone https://github.com/BrandonThomas84/code_review_automation.git
cd code_review_automation
make build
./bin/code-review --help
```

## 📖 Basic Usage

```bash
# Review changes against main branch
./code-review -t main

# Review with JSON output (for CI/CD)
./code-review -t main --json

# Full codebase scan (not just changed files)
./code-review -t main --full-scan

# Review and send email notification
./code-review -t main --email team@example.com

# Verbose output
./code-review -t main -v
```

### Command Reference

| Flag | Description |
| ------ | ------------- |
| `-t, --target` | **Required.** Target branch to compare against |
| `-o, --output` | Output directory for reports (default: `review_reports`) |
| `-j, --json` | Output results as JSON |
| `--full-scan` | Scan entire codebase, not just changed files |
| `--email` | Email address to send report to |
| `-v, --verbose` | Enable verbose output |

## 🔧 GitHub Actions Integration

Add automated code reviews to any repository by creating `.github/workflows/code-review.yml`:

```yaml
name: Code Review

on:
  pull_request:
    branches: [main, develop]

jobs:
  code-review:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write

    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Download AutoReview
        run: |
          curl -L https://github.com/BrandonThomas84/code_review_automation/releases/latest/download/code-review-linux-amd64 -o code-review
          chmod +x code-review

      - name: Run Code Review
        run: ./code-review -t ${{ github.base_ref }} --json > review_report.json

      - name: Comment on PR
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const report = JSON.parse(fs.readFileSync('review_report.json', 'utf8'));

            let comment = '## 📋 AutoReview Results\n\n';
            comment += `🔴 High: ${report.summary.high_severity} | `;
            comment += `🟡 Medium: ${report.summary.medium_severity} | `;
            comment += `🟢 Low: ${report.summary.low_severity}\n\n`;

            if (report.issues.length > 0) {
              report.issues.slice(0, 10).forEach((issue, i) => {
                comment += `${i + 1}. **[${issue.severity.toUpperCase()}]** ${issue.message}\n`;
                comment += `   📁 \`${issue.file}\`${issue.line ? ` (line ${issue.line})` : ''}\n`;
              });
            } else {
              comment += '✅ No issues found!';
            }

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });
```

> 💡 **Tip:** A complete workflow template with additional features is available at [`templates/github-actions-workflow.yml`](templates/github-actions-workflow.yml)

## 📧 Email Notifications

Send HTML-formatted review reports via email by setting these environment variables:

```bash
export AUTOREVIEW_SMTP_HOST="smtp.gmail.com"
export AUTOREVIEW_SMTP_USER="your-email@gmail.com"
export AUTOREVIEW_SMTP_PASSWORD="your-app-password"
export AUTOREVIEW_FROM_EMAIL="your-email@gmail.com"
export AUTOREVIEW_FROM_NAME="AutoReview Bot"  # Optional
```

### GitHub Actions Secrets

Add these secrets to your repository (Settings → Secrets and variables → Actions):

| Secret | Description |
| -------- | ------------- |
| `AUTOREVIEW_SMTP_HOST` | SMTP server hostname (e.g., `smtp.gmail.com`) |
| `AUTOREVIEW_SMTP_USER` | SMTP username/email |
| `AUTOREVIEW_SMTP_PASSWORD` | SMTP password or app password |
| `AUTOREVIEW_FROM_EMAIL` | Sender email address |

Then update your workflow:

```yaml
- name: Run Code Review with Email
  env:
    AUTOREVIEW_SMTP_HOST: ${{ secrets.AUTOREVIEW_SMTP_HOST }}
    AUTOREVIEW_SMTP_USER: ${{ secrets.AUTOREVIEW_SMTP_USER }}
    AUTOREVIEW_SMTP_PASSWORD: ${{ secrets.AUTOREVIEW_SMTP_PASSWORD }}
    AUTOREVIEW_FROM_EMAIL: ${{ secrets.AUTOREVIEW_FROM_EMAIL }}
  run: ./code-review -t ${{ github.base_ref }} --email team@example.com
```

> **Note:** Legacy variable names without the `AUTOREVIEW_` prefix are supported for backward compatibility.

## 🚫 Ignoring Files and Patterns

Create a `.autoreviewignore` file in your repository root (syntax similar to `.gitignore`):

```gitignore
# Ignore test files
*_test.go
**/__tests__/**

# Ignore generated files
*.generated.go
*.pb.go

# Ignore directories
vendor/
node_modules/
dist/
```

See the [AutoReview Ignore Guide](docs/AUTOREVIEW_IGNORE_GUIDE.md) for more details.

## 🏗️ Building from Source

### Prerequisites

- Go 1.25 or later
- Git

### Build Commands

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make test           # Run tests
make clean          # Clean build artifacts
```

### Output Binaries

| Platform | Binary |
| ---------- | -------- |
| Linux (amd64) | `bin/code-review-linux-amd64` |
| Linux (arm64) | `bin/code-review-linux-arm64` |
| macOS (Intel) | `bin/code-review-darwin-amd64` |
| macOS (Apple Silicon) | `bin/code-review-darwin-arm64` |
| Windows | `bin/code-review-windows-amd64.exe` |

## 📋 Customizing the Workflow

### Change Target Branches

```yaml
on:
  pull_request:
    branches:
      - main
      - develop
      - release/*
```

### Fail on High Severity Issues

```yaml
- name: Check for Critical Issues
  run: |
    HIGH=$(jq '.summary.high_severity' review_report.json)
    if [ "$HIGH" -gt 0 ]; then
      echo "❌ High severity issues found!"
      exit 1
    fi
```

### Filter by File Types

```yaml
on:
  pull_request:
    paths:
      - '**.py'
      - '**.js'
      - '**.ts'
```

## 🔍 Supported Languages & Checks

| Language | Security Checks | Quality Checks |
| ---------- | ----------------- | ---------------- |
| **Python** | SQL injection, eval(), exec(), pickle | print statements, debugger, TODO/FIXME |
| **JavaScript/TypeScript** | eval(), innerHTML, dangerouslySetInnerHTML | console.log, debugger, any type |
| **Ruby** | eval(), html_safe, YAML.load | debugger, binding.pry, puts |
| **Dart/Flutter** | Hardcoded credentials, HTTP URLs | print statements, dynamic type |
| **PHP** | SQL injection, eval(), shell_exec | var_dump, print_r, die |
| **Java** | Runtime.exec(), weak crypto | System.out.println, printStackTrace |
| **Kotlin** | Force unwrap (!!) | println, TODO |

## 📚 Documentation

| Document | Description |
| ---------- | ------------- |
| [Getting Started](docs/GETTING_STARTED.md) | 5-minute quick start guide |
| [Quick Reference](docs/QUICK_REFERENCE.md) | Command reference and common tasks |
| [Build & Deploy](docs/BUILD_AND_DEPLOY.md) | Detailed build and deployment instructions |
| [Integration Guide](docs/INTEGRATION_GUIDE.md) | Adding AutoReview to other repositories |
| [AutoReview Ignore Guide](docs/AUTOREVIEW_IGNORE_GUIDE.md) | Configuring file ignore patterns |
| [Deployment Checklist](docs/DEPLOYMENT_CHECKLIST.md) | Step-by-step deployment checklist |
| [Solution Summary](docs/SOLUTION_SUMMARY.md) | Technical architecture overview |

## 🔄 Creating Releases

Releases are automated via GitHub Actions. To create a new release:

1. Tag your commit with a version number:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create a GitHub Release
   - Attach the binaries as release assets

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes and add tests
4. Run tests: `make test`
5. Commit and push: `git push origin feature/my-feature`
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Made with ❤️ for faster, more consistent code reviews.**
