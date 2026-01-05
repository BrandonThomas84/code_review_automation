#!/bin/bash

# Setup script for code review automation tools
# Makes all scripts executable and sets up the environment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

echo "🚀 Setting up Code Review Automation Tools"
echo "=========================================="

# Make scripts executable
print_status "Making scripts executable..."
chmod +x "$SCRIPT_DIR"/*.sh
chmod +x "$SCRIPT_DIR"/*.py
chmod +x "$SCRIPT_DIR"/*.dart

print_success "Scripts are now executable"

# Check dependencies
print_status "Checking dependencies..."

# Check for git
if command -v git &> /dev/null; then
    print_success "Git is available"
else
    print_warning "Git not found - required for code analysis"
fi

# Check for Python
if command -v python3 &> /dev/null; then
    print_success "Python 3 is available"
else
    print_warning "Python 3 not found - some features will be limited"
fi

# Check for Dart (for Flutter projects)
if command -v dart &> /dev/null; then
    print_success "Dart is available"
else
    print_warning "Dart not found - Flutter analysis will be skipped"
fi

# Check for Ruby
if command -v ruby &> /dev/null; then
    print_success "Ruby is available"
else
    print_warning "Ruby not found - Ruby analysis will be skipped"
fi

# Check for Node.js (for JavaScript projects)
if command -v node &> /dev/null; then
    print_success "Node.js is available"
else
    print_warning "Node.js not found - JavaScript/TypeScript analysis will be skipped"
fi

# Check for npm (for JavaScript dependency checks)
if command -v npm &> /dev/null; then
    print_success "npm is available"
else
    print_warning "npm not found - JavaScript dependency checks will be limited"
fi

# Create alias suggestions
print_status "Creating convenience aliases..."

ALIAS_FILE="$SCRIPT_DIR/aliases.sh"
cat > "$ALIAS_FILE" << 'EOF'
#!/bin/bash
# Code Review Automation Aliases
# Source this file to add convenient aliases to your shell

# Get the directory where this script is located
REVIEW_TOOLS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Main review command (global)
alias autoreview="$REVIEW_TOOLS_DIR/autoreview"
alias code-review="$REVIEW_TOOLS_DIR/autoreview"

# Quick review
alias quick-review="$REVIEW_TOOLS_DIR/pr_analyzer.sh"

# Python detailed analysis
alias deep-review="python3 $REVIEW_TOOLS_DIR/quick_review.py"

# Language specific (work from current directory)
alias flutter-review="dart $REVIEW_TOOLS_DIR/flutter_review.dart ."
alias ruby-review="ruby $REVIEW_TOOLS_DIR/ruby_review.rb ."
alias js-review="node $REVIEW_TOOLS_DIR/js_ts_review.js ."

# Common review scenarios
alias review-pr="$REVIEW_TOOLS_DIR/autoreview -t main"
alias review-develop="$REVIEW_TOOLS_DIR/autoreview -t develop"
alias review-flutter="$REVIEW_TOOLS_DIR/autoreview -f"
alias review-ruby="$REVIEW_TOOLS_DIR/autoreview -r"
alias review-js="$REVIEW_TOOLS_DIR/autoreview -j"

echo "Code review aliases loaded! Available commands:"
echo "  autoreview       - Global comprehensive review (works from any directory)"
echo "  code-review      - Same as autoreview"
echo "  quick-review     - Fast PR analysis"
echo "  deep-review      - Detailed Python analysis"
echo "  flutter-review   - Flutter-specific analysis"
echo "  ruby-review      - Ruby-specific analysis"
echo "  js-review        - JavaScript/TypeScript analysis"
echo "  review-pr        - Review against main branch"
echo "  review-develop   - Review against develop branch"
echo "  review-flutter   - Review with Flutter checks"
echo "  review-ruby      - Review with Ruby checks"
echo "  review-js        - Review with JS/TS checks"
EOF

chmod +x "$ALIAS_FILE"
print_success "Aliases created in $ALIAS_FILE"

# Create a simple configuration file
CONFIG_FILE="$SCRIPT_DIR/config.sh"
cat > "$CONFIG_FILE" << 'EOF'
#!/bin/bash
# Code Review Automation Configuration

# Default target branch
DEFAULT_TARGET_BRANCH="main"

# Default output directory
DEFAULT_OUTPUT_DIR="review_reports"

# Enable/disable specific checks
ENABLE_SECURITY_CHECKS=true
ENABLE_QUALITY_CHECKS=true
ENABLE_FLUTTER_CHECKS=false

# Severity thresholds
HIGH_SEVERITY_THRESHOLD=5
MEDIUM_SEVERITY_THRESHOLD=10

# File size thresholds (lines)
LARGE_FILE_THRESHOLD=500
HUGE_FILE_THRESHOLD=1000

# Custom patterns to check (add your own)
CUSTOM_SECURITY_PATTERNS=(
    "your_secret_pattern_here"
)

CUSTOM_QUALITY_PATTERNS=(
    "your_quality_pattern_here"
)
EOF

print_success "Configuration file created at $CONFIG_FILE"

# Create a sample pre-commit hook
HOOK_FILE="$SCRIPT_DIR/pre-commit-hook.sh"
cat > "$HOOK_FILE" << 'EOF'
#!/bin/bash
# Sample pre-commit hook for code review automation
# Copy this to .git/hooks/pre-commit and make it executable

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REVIEW_TOOLS_DIR="$SCRIPT_DIR/../../code-review-automation"

echo "🔍 Running pre-commit code review..."

# Run quick analysis on staged files
if [[ -f "$REVIEW_TOOLS_DIR/quick_review.py" ]]; then
    python3 "$REVIEW_TOOLS_DIR/quick_review.py" --target HEAD~1
    
    # You can add logic here to fail the commit if critical issues are found
    # For example:
    # if [[ $? -ne 0 ]]; then
    #     echo "❌ Code review failed - commit aborted"
    #     exit 1
    # fi
fi

echo "✅ Pre-commit review completed"
EOF

chmod +x "$HOOK_FILE"
print_success "Sample pre-commit hook created at $HOOK_FILE"

#
# Final instructions
echo ""
echo "🎉 Setup Complete!"
echo "=================="
echo ""
echo "Next steps:"
echo "1. Add the CODE_REVIEW_TOOLS_DIR environment variable to your shell configuration (e.g., ~/.bashrc)"
echo "      export CODE_REVIEW_TOOLS_DIR=\"$SCRIPT_DIR\""
echo "2. Source aliases: source $ALIAS_FILE"
echo "3. Run a test review: ./review_master.sh"
echo "4. Customize config: edit $CONFIG_FILE"
echo "5. Read documentation: cat README.md"
echo ""
print_success "Code review automation is ready to use!"

echo ""
echo "🌍 GLOBAL ACCESS SETUP"
echo "======================"
echo ""
echo "To make 'autoreview' available from any directory, add this to your shell configuration:"
echo ""
echo "For ~/.bashrc or ~/.zshrc:"
echo "export PATH=\"$SCRIPT_DIR:\$PATH\""
echo "alias autoreview=\"$SCRIPT_DIR/autoreview\""
echo ""
echo "Or run these commands now:"
echo "echo 'export PATH=\"$SCRIPT_DIR:\$PATH\"' >> ~/.bashrc"
echo "echo 'alias autoreview=\"$SCRIPT_DIR/autoreview\"' >> ~/.bashrc"
echo "source ~/.bashrc"
echo ""
echo "After setup, you can run 'autoreview' from any directory!"
echo ""
print_success "Setup complete! Add the PATH and alias to your shell configuration."
