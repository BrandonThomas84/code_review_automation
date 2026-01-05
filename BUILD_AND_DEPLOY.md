# Build and Deployment Guide

This guide explains how to build the Go-based code review tool and integrate it into your GitHub workflows.

## Building the Binary

### Prerequisites
- Go 1.21 or later
- Make (optional, but recommended)
- Git

### Build for Your Platform

```bash
# Build for current OS/architecture
make build

# Binary will be in ./bin/code-review
./bin/code-review --help
```

### Build for All Platforms

```bash
# Build for Linux, macOS, and Windows
make build-all

# Binaries will be in ./bin/
ls -lh bin/
```

### Manual Build (without Make)

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o code-review-linux-amd64 ./cmd/code-review

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o code-review-darwin-amd64 ./cmd/code-review

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o code-review-darwin-arm64 ./cmd/code-review

# Windows
GOOS=windows GOARCH=amd64 go build -o code-review-windows-amd64.exe ./cmd/code-review
```

## Using in Other Repositories

### Option 1: GitHub Actions (Recommended)

1. Copy the workflow template to your repository:
```bash
mkdir -p .github/workflows
cp templates/github-actions-workflow.yml .github/workflows/code-review.yml
```

2. Customize the workflow:
   - Change branch names if needed
   - Add email configuration (see Email Setup below)
   - Adjust failure conditions

3. Commit and push:
```bash
git add .github/workflows/code-review.yml
git commit -m "Add automated code review workflow"
git push
```

### Option 2: Manual Integration

1. Download the binary for your platform:
```bash
# From releases page or build locally
wget https://github.com/BrandonThomas84/code-review-automation/releases/download/v1.0.0/code-review-linux-amd64
chmod +x code-review-linux-amd64
```

2. Run on a PR:
```bash
./code-review-linux-amd64 -t main --json > report.json
```

## Email Configuration

### Environment Variables

Set these in your GitHub Actions secrets or local environment:

```bash
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
export SMTP_USER="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export FROM_EMAIL="your-email@gmail.com"
```

### GitHub Actions Setup

1. Add secrets to your repository:
   - Go to Settings → Secrets and variables → Actions
   - Add: `SMTP_HOST`, `SMTP_USER`, `SMTP_PASSWORD`, `FROM_EMAIL`

2. Update workflow to use secrets:
```yaml
- name: Run code review with email
  env:
    SMTP_HOST: ${{ secrets.SMTP_HOST }}
    SMTP_USER: ${{ secrets.SMTP_USER }}
    SMTP_PASSWORD: ${{ secrets.SMTP_PASSWORD }}
    FROM_EMAIL: ${{ secrets.FROM_EMAIL }}
  run: |
    ./code-review -t main --email your-email@example.com
```

## Release Process

### Creating a Release

1. Tag the version:
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

2. Build all binaries:
```bash
make build-all
```

3. Create GitHub release with binaries:
```bash
# Using GitHub CLI
gh release create v1.0.0 bin/* --title "v1.0.0" --notes "Release notes here"
```

## Continuous Integration

The repository includes a GitHub Actions workflow that:
- Builds binaries for all platforms on each release
- Runs tests
- Creates release artifacts

See `.github/workflows/` for the CI configuration.

## Troubleshooting

### Binary not found
- Ensure you're in the correct directory
- Check that the binary has execute permissions: `chmod +x code-review`

### Git diff fails
- Ensure you're in a git repository
- Verify the target branch exists: `git branch -a`
- Check git is installed: `git --version`

### Email not sending
- Verify SMTP credentials
- Check firewall/network access to SMTP server
- Enable "Less secure app access" for Gmail
- Use app-specific passwords for Gmail

## Next Steps

1. Build the binary: `make build`
2. Test locally: `./bin/code-review -t main`
3. Create a release on GitHub
4. Add workflow to your target repositories
5. Configure email notifications

