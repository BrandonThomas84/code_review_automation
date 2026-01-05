# Getting Started with Code Review Automation

## 5-Minute Quick Start

### Step 1: Build the Binary (2 minutes)

```bash
# Clone or navigate to the repository
cd code-review-automation

# Build the Go binary
make build

# Verify it works
./bin/code-review --help
```

### Step 2: Test Locally (2 minutes)

```bash
# Navigate to any git repository
cd /path/to/your/project

# Run a code review
/path/to/code-review-automation/bin/code-review -t main

# Or with JSON output
/path/to/code-review-automation/bin/code-review -t main --json
```

### Step 3: Add to GitHub Actions (1 minute)

```bash
# In your target repository
mkdir -p .github/workflows

# Copy the workflow template
cp /path/to/code-review-automation/templates/github-actions-workflow.yml \
   .github/workflows/code-review.yml

# Commit and push
git add .github/workflows/code-review.yml
git commit -m "Add automated code review"
git push
```

## Detailed Setup

### Prerequisites

- Go 1.21+ (for building)
- Git
- Make (optional, but recommended)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/BrandonThomas84/code-review-automation.git
cd code-review-automation
```

2. **Build the binary**
```bash
# For your current platform
make build

# For all platforms (Linux, macOS, Windows)
make build-all
```

3. **Verify installation**
```bash
./bin/code-review --version
./bin/code-review --help
```

## Using in Your Projects

### Option A: GitHub Actions (Recommended)

1. **Copy workflow to your repository**
```bash
mkdir -p .github/workflows
cp templates/github-actions-workflow.yml .github/workflows/code-review.yml
```

2. **Customize the workflow** (optional)
   - Edit `.github/workflows/code-review.yml`
   - Change branch names if needed
   - Add email configuration

3. **Commit and push**
```bash
git add .github/workflows/code-review.yml
git commit -m "Add code review workflow"
git push
```

4. **Create a pull request**
   - The workflow will run automatically
   - Results will appear as a PR comment

### Option B: Local Usage

```bash
# Review changed files
code-review -t main

# Full codebase scan
code-review -t main --full-scan

# JSON output for parsing
code-review -t main --json > report.json

# Verbose output
code-review -t main -v
```

### Option C: CI/CD Integration

Use in any CI/CD system:

```bash
# Download the binary
wget https://github.com/BrandonThomas84/code-review-automation/releases/download/v1.0.0/code-review-linux-amd64
chmod +x code-review-linux-amd64

# Run the review
./code-review-linux-amd64 -t main --json > report.json

# Process the report
cat report.json | jq '.summary'
```

## Email Configuration (Optional)

### Set Environment Variables

```bash
export SMTP_HOST="smtp.gmail.com"
export SMTP_USER="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export FROM_EMAIL="your-email@gmail.com"
```

### In GitHub Actions

1. Go to repository Settings → Secrets and variables → Actions
2. Add these secrets:
   - `SMTP_HOST`
   - `SMTP_USER`
   - `SMTP_PASSWORD`
   - `FROM_EMAIL`

3. Update workflow to use secrets:
```yaml
- name: Run code review with email
  env:
    SMTP_HOST: ${{ secrets.SMTP_HOST }}
    SMTP_USER: ${{ secrets.SMTP_USER }}
    SMTP_PASSWORD: ${{ secrets.SMTP_PASSWORD }}
    FROM_EMAIL: ${{ secrets.FROM_EMAIL }}
  run: |
    ./code-review -t ${{ github.base_ref }} \
      --email your-email@example.com
```

## Common Tasks

### Build for Distribution

```bash
# Build all platforms
make build-all

# Binaries are in ./bin/
ls -lh bin/

# Create a release on GitHub
gh release create v1.0.0 bin/* --title "v1.0.0"
```

### Run Tests

```bash
make test
```

### Clean Build Artifacts

```bash
make clean
```

### Update Dependencies

```bash
go mod tidy
```

## Troubleshooting

### "command not found: make"
Install Make:
- **Ubuntu/Debian**: `sudo apt-get install make`
- **macOS**: `brew install make`
- **Windows**: Use `go build` directly

### "go: command not found"
Install Go from https://golang.org/dl/

### Binary doesn't work
```bash
# Make it executable
chmod +x ./bin/code-review

# Run with full path
./bin/code-review -t main
```

### Git diff shows no changes
```bash
# Ensure you're in a git repository
git status

# Verify target branch exists
git branch -a

# Check git is installed
git --version
```

## Next Steps

1. ✅ Build the binary: `make build`
2. ✅ Test locally: `./bin/code-review -t main`
3. ✅ Add to a repository: Copy workflow template
4. ✅ Create a PR to test
5. ✅ Configure email (optional)
6. ✅ Share with your team!

## Documentation

- **Quick Reference**: `QUICK_REFERENCE.md`
- **Build Guide**: `BUILD_AND_DEPLOY.md`
- **Integration Guide**: `INTEGRATION_GUIDE.md`
- **Go Implementation**: `GO_IMPLEMENTATION_SUMMARY.md`
- **Main README**: `README.md`

## Support

- GitHub Issues: https://github.com/BrandonThomas84/code-review-automation/issues
- Discussions: https://github.com/BrandonThomas84/code-review-automation/discussions

