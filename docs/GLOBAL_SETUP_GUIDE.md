# Global Code Review Automation Setup Guide

## 🌍 Global Access Configuration

Your code review automation system is now designed to work from **any directory** and automatically detect **monorepos** with multiple languages.

## 📋 Setup Instructions

### 1. Run the Setup Script
```bash
cd ~/code-review-automation
./setup.sh
```

### 2. Add to Your Shell Configuration

**For Bash users (~/.bashrc):**
```bash
echo 'export PATH="~/code-review-automation:$PATH"' >> ~/.bashrc
echo 'alias autoreview="~/code-review-automation/autoreview"' >> ~/.bashrc
source ~/.bashrc
```

**For Zsh users (~/.zshrc):**
```bash
echo 'export PATH="~/code-review-automation:$PATH"' >> ~/.zshrc
echo 'alias autoreview="~/code-review-automation/autoreview"' >> ~/.zshrc
source ~/.zshrc
```

**For Fish users (~/.config/fish/config.fish):**
```bash
echo 'set -gx PATH ~/code-review-automation $PATH' >> ~/.config/fish/config.fish
echo 'alias autoreview="~/code-review-automation/autoreview"' >> ~/.config/fish/config.fish
```

### 3. Verify Installation
```bash
# Test from any directory
cd /tmp
autoreview --help
```

## 🚀 Usage Examples

### Basic Usage
```bash
# Navigate to any project directory
cd /path/to/your/project

# Run comprehensive review
autoreview

# Review against specific branch
autoreview -t develop

# Verbose output
autoreview -v
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

## 🏗️ Monorepo Support

The system automatically detects multiple languages in monorepos:

### Example Output
```
Project type: monorepo (ruby,javascript,php,java)
```

### Supported Detection
- **Ruby**: Gemfile, *.rb files
- **JavaScript/TypeScript**: package.json, *.js, *.ts, *.jsx, *.tsx files
- **Flutter/Dart**: pubspec.yaml, *.dart files
- **Python**: requirements.txt, pyproject.toml, *.py files
- **PHP**: composer.json, *.php files
- **Java**: pom.xml, build.gradle, *.java files

### Automatic Analysis
When multiple languages are detected, the system will:
1. Run general security and quality checks
2. Execute language-specific analyzers for each detected language
3. Generate separate reports for each language
4. Combine everything into a comprehensive summary

## 📊 Directory Structure

The system works with any directory structure:

```
your-project/
├── frontend/           # JavaScript/TypeScript
│   ├── package.json
│   └── src/
├── backend/           # Ruby/Rails
│   ├── Gemfile
│   └── app/
├── mobile/           # Flutter
│   ├── pubspec.yaml
│   └── lib/
└── api/             # Python
    ├── requirements.txt
    └── src/
```

Running `autoreview` from `your-project/` will detect and analyze all languages.

## 🎯 Key Benefits

### 1. **Global Accessibility**
- Run from any directory
- No need to navigate to script location
- Consistent command across all projects

### 2. **Intelligent Detection**
- Automatically identifies project languages
- Handles monorepos with multiple languages
- No manual configuration required

### 3. **Context Awareness**
- Analyzes the current working directory
- Respects git repository boundaries
- Generates reports in the current location

### 4. **Flexible Usage**
- Override automatic detection with flags
- Customize output locations
- Skip specific analysis types

## 🔧 Advanced Configuration

### Custom Aliases
Add these to your shell configuration for even faster access:

```bash
# Quick aliases
alias review="autoreview"
alias review-pr="autoreview -t main"
alias review-dev="autoreview -t develop"
alias review-ruby="autoreview -r"
alias review-js="autoreview -j"
alias review-flutter="autoreview -f"
```

### CI/CD Integration
```yaml
# GitHub Actions example
name: Code Review
on: [pull_request]
jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Code Review Tools
        run: |
          git clone https://github.com/your-org/code-review-automation.git
          chmod +x code-review-automation/autoreview
      - name: Run Code Review
        run: ./code-review-automation/autoreview --no-interactive
```

### Pre-commit Hook
```bash
# Install pre-commit hook
cp ~/code-review-automation/pre-commit-hook.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

## 📈 Workflow Integration

### Daily Development
```bash
# Before committing
cd my-project
autoreview

# Quick check during development
autoreview --no-security --no-quality  # Just language-specific checks
```

### Code Review Process
```bash
# Reviewer workflow
cd project-directory
autoreview -t feature-branch  # Compare against feature branch
autoreview -v                 # Verbose output for detailed review
```

### Team Standards
```bash
# Enforce team standards
autoreview --no-quality  # Focus on security and language-specific issues
autoreview -r -j         # Only Ruby and JavaScript in mixed repos
```

## 🎉 Success Indicators

After setup, you should be able to:

1. ✅ Run `autoreview` from any directory
2. ✅ See automatic language detection
3. ✅ Get comprehensive reports in current directory
4. ✅ Use language-specific flags to override detection
5. ✅ Generate timestamped reports in `review_reports/`

## 🆘 Troubleshooting

### Command Not Found
```bash
# Check PATH
echo $PATH | grep code-review-automation

# Re-source shell configuration
source ~/.bashrc  # or ~/.zshrc
```

### Permission Issues
```bash
# Make scripts executable
chmod +x ~/code-review-automation/*
```

### Git Repository Issues
```bash
# Ensure you're in a git repository
git status

# Initialize if needed
git init
```

The system is now ready for global use across all your projects! 🚀
