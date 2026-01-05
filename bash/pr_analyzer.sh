#!/bin/bash

# PR Analyzer Script
# Comprehensive analysis of Pull Requests for code review

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
TARGET_BRANCH="main"
REPO_PATH="$(pwd)"  # Use current working directory
OUTPUT_FILE=""
VERBOSE=false

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Function to show usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -t, --target BRANCH    Target branch to compare against (default: main)"
    echo "  -o, --output FILE      Output file for report"
    echo "  -v, --verbose          Verbose output"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -t FERN-12345-Some-title"
    echo "  $0 --target FERN-12345-Some-title --verbose"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--target)
            TARGET_BRANCH="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not a git repository: $(pwd)"
    exit 1
fi

print_status "Starting PR analysis..."
print_status "Target branch: $TARGET_BRANCH"
print_status "Repository: $(pwd)"

# Create output file if specified
if [[ -n "$OUTPUT_FILE" ]]; then
    exec > >(tee "$OUTPUT_FILE")
fi

echo "========================================"
echo "🔍 PULL REQUEST ANALYSIS REPORT"
echo "========================================"
echo "Generated: $(date)"
echo "Target Branch: $TARGET_BRANCH"
echo "Current Branch: $(git branch --show-current)"
echo ""

# 1. Basic Git Analysis
print_status "Analyzing git changes..."

# Get changed files from git diff
CHANGED_FILES=$(git diff --name-only "$TARGET_BRANCH"..HEAD 2>/dev/null || echo "")

# Handle empty CHANGED_FILES
if [[ -z "$CHANGED_FILES" ]]; then
    print_warning "No changes detected between $TARGET_BRANCH and HEAD"
    CHANGED_FILES_COUNT=0
    ADDITIONS=0
    DELETIONS=0
else
    print_success "Changes detected between $TARGET_BRANCH and HEAD"
    CHANGED_FILES_COUNT=$(echo "$CHANGED_FILES" | wc -l)
    
    # Get git stats and parse more carefully
    STAT_OUTPUT=$(git diff --stat "$TARGET_BRANCH"..HEAD 2>/dev/null | tail -1)
    ADDITIONS=$(echo "$STAT_OUTPUT" | grep -o '[0-9]\+ insertion' | head -1 | grep -o '[0-9]\+')
    DELETIONS=$(echo "$STAT_OUTPUT" | grep -o '[0-9]\+ deletion' | head -1 | grep -o '[0-9]\+')
    
    # Ensure we have single numeric values
    ADDITIONS=${ADDITIONS:-0}
    DELETIONS=${DELETIONS:-0}
    
    # Remove any whitespace/newlines
    ADDITIONS=$(echo "$ADDITIONS" | tr -d '\n\r ')
    DELETIONS=$(echo "$DELETIONS" | tr -d '\n\r ')
fi


echo "📊 CHANGE STATISTICS"
echo "===================="
echo "Files changed: $CHANGED_FILES_COUNT"
echo "Lines added: $ADDITIONS"
echo "Lines deleted: $DELETIONS"
echo ""

if [[ $VERBOSE == true && -n "$CHANGED_FILES" ]]; then
    echo "Changed files:"
    echo "$CHANGED_FILES"
    echo ""
fi

# 2. Security Check
print_status "Running security checks..."
echo "🔒 SECURITY ANALYSIS"
echo "===================="

security_issues_found=false

if [[ -n "$CHANGED_FILES" ]]; then
    git diff "$TARGET_BRANCH"..HEAD > /tmp/pr_diff.txt
    
    # Check for common security issues
    security_patterns=(
        "password.*=.*['\"][^'\"]*['\"]"
        "api_key.*=.*['\"][^'\"]*['\"]"
        "secret.*=.*['\"][^'\"]*['\"]"
        "token.*=.*['\"][^'\"]*['\"]"
        "private_key"
        "BEGIN.*PRIVATE KEY"
    )

    for pattern in "${security_patterns[@]}"; do
        if grep -i "$pattern" /tmp/pr_diff.txt > /dev/null 2>&1; then
            echo "🚨 Potential security issue: $pattern"
            security_issues_found=true
        fi
    done
fi

if [[ $security_issues_found == false ]]; then
    echo "✅ No obvious security issues detected"
fi
echo ""

# 5. Code Quality Check
print_status "Checking code quality..."
echo "🎯 CODE QUALITY ANALYSIS"
echo "========================"

quality_issues_found=false

if [[ -n "$CHANGED_FILES" ]]; then
    # Check for common code quality issues in the diff
    quality_patterns=(
        "console\.log"
        "debugger\;"
        "TODO"
        "FIXME"
        "HACK"
        "\.catch\(\s*\)"
    )

    for pattern in "${quality_patterns[@]}"; do
        matches=$(grep -c "$pattern" /tmp/pr_diff.txt 2>/dev/null || echo "0")
        # Clean the matches variable to ensure it's a single number
        matches=$(echo "$matches" | head -1 | tr -d ' \n\r')
        matches=${matches:-0}

        if [[ "$matches" -gt 0 ]]; then
            echo "⚠️  Found $matches instances of: $pattern"
            quality_issues_found=true
        fi
    done
fi

if [[ $quality_issues_found == false ]]; then
    echo "✅ No obvious code quality issues detected"
fi
echo ""

# 6. Test Coverage Check
print_status "Checking test coverage..."
echo "🧪 TEST COVERAGE ANALYSIS"
echo "========================="

if [[ -n "$CHANGED_FILES" ]]; then
    test_files=$(echo "$CHANGED_FILES" | grep -E "(test|spec)" || true)
    code_files=$(echo "$CHANGED_FILES" | grep -E "\.(dart|js|ts|py|php|java|kt)$" | grep -v -E "(test|spec)" || true)

    test_count=$(echo "$test_files" | grep -c . 2>/dev/null || echo "0")
    code_count=$(echo "$code_files" | grep -c . 2>/dev/null || echo "0")

    # Clean whitespace from counts
    test_count=$(echo "$test_count" | tr -d ' \n\r')
    code_count=$(echo "$code_count" | tr -d ' \n\r')

    # Fix empty string counting
    if [[ -z "$test_files" ]]; then
        test_count=0
    fi
    if [[ -z "$code_files" ]]; then
        code_count=0
    fi
else
    test_count=0
    code_count=0
fi

echo "Code files changed: $code_count"
echo "Test files changed: $test_count"

if [[ "$code_count" -gt 0 && "$test_count" -eq 0 ]]; then
    echo "⚠️  Code changes detected but no test files modified"
    echo "💡 Consider adding or updating tests"
elif [[ "$test_count" -gt 0 ]]; then
    echo "✅ Test files included in changes"
else
    echo "ℹ️  No code files requiring tests detected"
fi
echo ""

# 4. Security Check
print_status "Running security checks..."
echo "🔒 SECURITY ANALYSIS"
echo "===================="

security_issues_found=false

if [[ -n "$CHANGED_FILES" ]]; then
    git diff "$TARGET_BRANCH"..HEAD > /tmp/pr_diff.txt
    
    # Check for common security issues
    security_patterns=(
        "password.*=.*['\"][^'\"]*['\"]"
        "api_key.*=.*['\"][^'\"]*['\"]"
        "secret.*=.*['\"][^'\"]*['\"]"
        "token.*=.*['\"][^'\"]*['\"]"
        "private_key"
        "BEGIN.*PRIVATE KEY"
    )

    for pattern in "${security_patterns[@]}"; do
        if grep -i "$pattern" /tmp/pr_diff.txt > /dev/null 2>&1; then
            echo "🚨 Potential security issue: $pattern"
            security_issues_found=true
        fi
    done
fi

if [[ $security_issues_found == false ]]; then
    echo "✅ No obvious security issues detected"
fi
echo ""

# 8. Documentation Check
print_status "Checking documentation..."
echo "📚 DOCUMENTATION ANALYSIS"
echo "========================="

doc_files=$(echo "$CHANGED_FILES" | grep -E "\.(md|txt|rst|adoc)$" || echo "")
doc_count=$(echo "$doc_files" | grep -c . 2>/dev/null || echo "0")

# Clean whitespace from doc_count
doc_count=$(echo "$doc_count" | tr -d ' \n\r')

if [[ $doc_count -gt 0 ]]; then
    echo "✅ Documentation files updated: $doc_count"
    if [[ $VERBOSE == true ]]; then
        echo "$doc_files"
    fi
else
    echo "ℹ️  No documentation files changed"
    if [[ "$code_count" -gt 5 ]]; then
        echo "💡 Consider updating documentation for significant code changes"
    fi
fi
echo ""

# 9. Final Summary
echo "📋 REVIEW SUMMARY"
echo "================="
echo "✅ Analysis complete"
echo "📁 Total files analyzed: $CHANGED_FILES_COUNT"
echo "➕ Lines added: $ADDITIONS"
echo "➖ Lines deleted: $DELETIONS"

# Cleanup
rm -f /tmp/pr_diff.txt

print_success "PR analysis completed!"

if [[ -n "$OUTPUT_FILE" ]]; then
    print_success "Report saved to: $OUTPUT_FILE"
fi

exit 0
