# Code Review Automation Tools

Comprehensive code review automation for your projects. Run code reviews from any directory with automatic language detection and monorepo support.

## Quick Start

### 1. Clone the Repository

Clone this repository to your desired location:

```bash
git clone https://github.com/your-org/code-review-automation.git
cd code-review-automation
```

### 2. Run Initial Setup

```bash
./setup.sh
```

### 3. Add Environment Variable to Your Shell Configuration

Add the `CODE_REVIEW_TOOLS_DIR` environment variable pointing to your cloned repository location.

**For Bash (~/.bashrc):**

```bash
echo 'export CODE_REVIEW_TOOLS_DIR="/path/to/code-review-automation"' >> ~/.bashrc
source ~/.bashrc
```

**For Zsh (~/.zshrc):**

```bash
echo 'export CODE_REVIEW_TOOLS_DIR="/path/to/code-review-automation"' >> ~/.zshrc
source ~/.zshrc
```

**For Fish (~/.config/fish/config.fish):**

```bash
echo 'set -gx CODE_REVIEW_TOOLS_DIR "/path/to/code-review-automation"' >> ~/.config/fish/config.fish
```

Replace `/path/to/code-review-automation` with the actual path where you cloned the repository.

### Important: Target Branch is Required

The `-t` or `--target` flag is **required** when running autoreview. This specifies which branch to compare your changes against.

### 4. Verify Installation

Test from any directory:

```bash
cd /tmp
autoreview --help
```

### 5. Basic Usage

Navigate to any project and run:

```bash
autoreview -t FERN-12345-Some-title
```

That's it! The tool will automatically detect your project languages and run comprehensive analysis.

## Usage Examples

### Basic Review (Changed Files Only)

```bash
# Review changed files with ticket/branch name
autoreview -t FERN-12345-Some-title

# Review changed files against specific branch
autoreview -t develop

# Verbose output for detailed review
autoreview -v
```

### Full Codebase Scan

```bash
# Scan entire codebase instead of just changed files
autoreview --full-scan

# Full scan against specific branch
autoreview -t develop --full-scan

# Full scan with verbose output
autoreview --full-scan -v
```

### Language-Specific Analysis

```bash
# Force Ruby analysis
autoreview -r

# Force JavaScript/TypeScript analysis
autoreview -j

# Force Flutter analysis
autoreview -f

# Multiple languages (for monorepos)
autoreview -r -j -f
```

### Output Control

```bash
# Custom output directory
autoreview -o /tmp/review-reports

# Skip security checks
autoreview --no-security

# Skip quality checks
autoreview --no-quality
```

## Features

### Automatic Language Detection

- **Ruby**: Gemfile, *.rb files
- **JavaScript/TypeScript**: package.json, *.js,*.ts, *.jsx,*.tsx files
- **Flutter/Dart**: pubspec.yaml, *.dart files
- **Python**: requirements.txt, pyproject.toml, *.py files
- **PHP**: composer.json, *.php files
- **Java**: pom.xml, build.gradle, *.java files

### Analysis Capabilities

- **Security**: Hardcoded secrets, dangerous functions, dependency vulnerabilities
- **Code Quality**: Large files, complexity analysis, TODO/FIXME tracking, empty catch blocks
- **Language-Specific**: Widget analysis (Flutter), state management, performance optimizations
- **Monorepo Support**: Automatically detects and analyzes multiple languages in one project

## Output

All reports are saved to `review_reports/` directory with timestamps:

- `review_report_TIMESTAMP.txt` - Main comprehensive report
- Language-specific reports for each detected language
- Timestamped for easy tracking and comparison

## Dependencies

- **Git** (required)
- **Python 3** (for detailed analysis)
- **Dart** (for Flutter analysis)
- **Node.js/npm** (for JavaScript dependency checks)
- **Ruby** (for Ruby analysis)

## Advanced Usage

### Pre-commit Hook

```bash
cp $CODE_REVIEW_TOOLS_DIR/pre-commit-hook.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### CI/CD Integration

**GitHub Actions:**

```yaml
- name: Setup Code Review Tools
  run: |
    git clone https://github.com/your-org/code-review-automation.git
    echo "export CODE_REVIEW_TOOLS_DIR=$(pwd)/code-review-automation" >> $GITHUB_ENV

- name: Code Review
  run: $CODE_REVIEW_TOOLS_DIR/autoreview --no-interactive
```

### Custom Aliases

Add to your shell configuration (after setting `CODE_REVIEW_TOOLS_DIR`):

```bash
alias review="autoreview"
alias review-pr="autoreview -t main"
alias review-dev="autoreview -t develop"
```

## Troubleshooting

**Command not found?**

```bash
# Check if environment variable is set
echo $CODE_REVIEW_TOOLS_DIR

# If empty, add to your shell configuration and reload
export CODE_REVIEW_TOOLS_DIR="/path/to/code-review-automation"
source ~/.bashrc  # or source ~/.zshrc
```

**Permission issues?**

```bash
chmod +x $CODE_REVIEW_TOOLS_DIR/*
```

**Scripts not executable?**

```bash
# Make all scripts executable
chmod +x $CODE_REVIEW_TOOLS_DIR/*.sh
chmod +x $CODE_REVIEW_TOOLS_DIR/autoreview
```

For more detailed setup information, see [GLOBAL_SETUP_GUIDE.md](GLOBAL_SETUP_GUIDE.md).
