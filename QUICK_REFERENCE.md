# Quick Reference Guide

## Building the Go Binary

```bash
# Build for your current system
make build
./bin/code-review --help

# Build for all platforms
make build-all
ls -lh bin/

# Clean build artifacts
make clean
```

## Running Code Reviews

### Basic Usage
```bash
# Review changed files against main branch
code-review -t main

# Review against develop branch
code-review -t develop

# Full codebase scan
code-review -t main --full-scan

# Verbose output
code-review -t main -v

# JSON output
code-review -t main --json
```

### With Email
```bash
# Send report to email
code-review -t main --email user@example.com

# Requires SMTP environment variables:
export SMTP_HOST=smtp.gmail.com
export SMTP_USER=your-email@gmail.com
export SMTP_PASSWORD=your-app-password
export FROM_EMAIL=your-email@gmail.com
```

### Custom Output
```bash
# Save to specific directory
code-review -t main -o /tmp/reports

# JSON to file
code-review -t main --json > report.json
```

## GitHub Actions Integration

### 1. Add Workflow to Your Repo
```bash
mkdir -p .github/workflows
cp templates/github-actions-workflow.yml {TARGET REPOSITORY}/.github/workflows/code-review.yml
cd {TARGET REPOSITORY}
git add .github/workflows/code-review.yml
git commit -m "Add code review workflow"
git push
```

### 2. Configure Email (Optional)
Go to repository Settings → Secrets and variables → Actions, add:
- `SMTP_HOST`
- `SMTP_USER`
- `SMTP_PASSWORD`
- `FROM_EMAIL`

### 3. Test
Create a pull request - the workflow will run automatically!

## Workflow Customization

### Change Target Branches
Edit `.github/workflows/code-review.yml`:
```yaml
on:
  pull_request:
    branches:
      - main
      - develop
      - staging
```

### Fail on High Severity Issues
Add to workflow:
```yaml
- name: Check severity
  run: |
    HIGH=$(jq '.summary.high_severity' review_report.json)
    if [ "$HIGH" -gt 0 ]; then exit 1; fi
```

### Add Email Notification
```yaml
- name: Send email
  env:
    SMTP_HOST: ${{ secrets.SMTP_HOST }}
    SMTP_USER: ${{ secrets.SMTP_USER }}
    SMTP_PASSWORD: ${{ secrets.SMTP_PASSWORD }}
    FROM_EMAIL: ${{ secrets.FROM_EMAIL }}
  run: |
    ./code-review -t ${{ github.base_ref }} \
      --email your-email@example.com
```

## Troubleshooting

### Build fails
```bash
# Ensure Go is installed
go version

# Update dependencies
go mod tidy

# Try clean build
make clean && make build
```

### Binary not found
```bash
# Check current directory
pwd

# List binaries
ls -lh bin/

# Make executable
chmod +x bin/code-review
```

### Git diff shows no changes
```bash
# Verify target branch exists
git branch -a

# Check git is initialized
git status

# Ensure fetch-depth: 0 in GitHub Actions
```

### Email not sending
```bash
# Test SMTP credentials
# Verify environment variables are set
env | grep SMTP

# Check firewall/network access
# Use app-specific passwords for Gmail
```

## Common Commands

```bash
# Show help
code-review --help

# Show version
code-review version

# Show configuration
code-review config show

# Review with all options
code-review -t main -v --json -o ./reports --email user@example.com
```

## File Locations

- **Binary**: `./bin/code-review`
- **Workflow template**: `./templates/github-actions-workflow.yml`
- **Build guide**: `./BUILD_AND_DEPLOY.md`
- **Integration guide**: `./INTEGRATION_GUIDE.md`
- **Go source**: `./cmd/` and `./internal/`

## Release Checklist

- [ ] Update version in `internal/cmd/version.go`
- [ ] Test build: `make build`
- [ ] Test all platforms: `make build-all`
- [ ] Run tests: `make test`
- [ ] Tag release: `git tag -a v1.0.0 -m "Release v1.0.0"`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] GitHub Actions creates release automatically

## Support Resources

- **Main README**: `README.md`
- **Build Guide**: `BUILD_AND_DEPLOY.md`
- **Integration Guide**: `INTEGRATION_GUIDE.md`
- **Go Summary**: `GO_IMPLEMENTATION_SUMMARY.md`
- **GitHub Repo**: https://github.com/BrandonThomas84/code-review-automation

