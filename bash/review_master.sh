#!/bin/bash

# Master Code Review Script
# Orchestrates all review tools for comprehensive analysis

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

# Check if CODE_REVIEW_TOOLS_DIR environment variable is set
if [[ -z "$CODE_REVIEW_TOOLS_DIR" ]]; then
    echo -e "${RED}❌ Error: CODE_REVIEW_TOOLS_DIR environment variable is not set${NC}"
    echo ""
    echo "Please add the following to your shell configuration (~/.bashrc, ~/.zshrc, etc.):"
    echo "  export CODE_REVIEW_TOOLS_DIR=\"/path/to/code-review-automation\""
    echo ""
    echo "Then reload your shell:"
    echo "  source ~/.bashrc  # or source ~/.zshrc"
    exit 1
fi

# Verify the directory exists
if [[ ! -d "$CODE_REVIEW_TOOLS_DIR" ]]; then
    echo -e "${RED}❌ Error: CODE_REVIEW_TOOLS_DIR points to a non-existent directory: $CODE_REVIEW_TOOLS_DIR${NC}"
    exit 1
fi

# Configuration
SCRIPT_DIR="$CODE_REVIEW_TOOLS_DIR"
TARGET_BRANCH="main"
REPO_PATH="$(pwd)"  # Use current working directory
OUTPUT_DIR="review_reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
VERBOSE=false
RUN_FLUTTER=false
RUN_RUBY=false
RUN_JS=false
RUN_SECURITY=true
RUN_QUALITY=true
FULL_SCAN=false

# Function to print colored output
print_header() {
    echo -e "${PURPLE}$1${NC}"
}

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Usage function
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Master code review automation script"
    echo ""
    echo "Options:"
    echo "  -t, --target BRANCH     Target branch to compare against (default: main)"
    echo "  -o, --output DIR        Output directory for reports (default: review_reports)"
    echo "  -f, --flutter           Run Flutter-specific analysis"
    echo "  -r, --ruby              Run Ruby-specific analysis"
    echo "  -j, --javascript        Run JavaScript/TypeScript analysis"
    echo "  --full-scan             Scan entire codebase (default: only changed files)"
    echo "  --no-security           Skip security analysis"
    echo "  --no-quality            Skip quality analysis"
    echo "  --clean                 Clean up old report files and exit"
    echo "  -v, --verbose           Verbose output"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -t develop -f                    # Analyze changed files against develop with Flutter checks"
    echo "  $0 --target main --full-scan        # Full codebase scan against main"
    echo "  $0 -o /tmp/reports                  # Custom output directory"
    echo "  $0 --clean                          # Clean up old reports"
}

# Cleanup function
cleanup_reports() {
    print_header "🧹 CLEANING UP REPORT FILES"
    print_header "============================"
    
    if [[ ! -d "$OUTPUT_DIR" ]]; then
        print_warning "Output directory $OUTPUT_DIR does not exist"
        return 0
    fi
    
    # Count files before cleanup
    REPORT_COUNT=$(find "$OUTPUT_DIR" -name "*.txt" -o -name "*.json" | wc -l)
    TEMP_COUNT=$(find /tmp -name "pr_diff.txt" 2>/dev/null | wc -l)
    
    if [[ "$REPORT_COUNT" -eq 0 && "$TEMP_COUNT" -eq 0 ]]; then
        print_status "No report files found to clean up"
        return 0
    fi
    
    print_status "Found $REPORT_COUNT report files in $OUTPUT_DIR"
    if [[ "$TEMP_COUNT" -gt 0 ]]; then
        print_status "Found $TEMP_COUNT temporary files in /tmp"
    fi
    
    # Ask for confirmation unless in non-interactive mode
    if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then  # Only prompt if run directly
        echo -n "Are you sure you want to delete all report files? (y/N): "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_status "Cleanup cancelled"
            return 0
        fi
    fi
    
    # Remove report files
    if [[ "$REPORT_COUNT" -gt 0 ]]; then
        find "$OUTPUT_DIR" -name "*.txt" -delete
        find "$OUTPUT_DIR" -name "*.json" -delete
        print_success "Removed $REPORT_COUNT report files from $OUTPUT_DIR"
    fi
    
    # Remove temporary files
    if [[ "$TEMP_COUNT" -gt 0 ]]; then
        rm -f /tmp/pr_diff.txt
        print_success "Removed temporary files from /tmp"
    fi

    
    # Remove empty output directory if it exists
    if [[ -d "$OUTPUT_DIR" ]] && [[ -z "$(ls -A "$OUTPUT_DIR")" ]]; then
        echo -n "Would you like to remove the empty directory? (y/N): "
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            rmdir "$OUTPUT_DIR"
            print_success "Removed empty directory $OUTPUT_DIR"
        fi
    fi
    
    print_success "Cleanup completed!"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--target)
            TARGET_BRANCH="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -f|--flutter)
            RUN_FLUTTER=true
            shift
            ;;
        -r|--ruby)
            RUN_RUBY=true
            shift
            ;;
        -j|--javascript)
            RUN_JS=true
            shift
            ;;
        --full-scan)
            FULL_SCAN=true
            shift
            ;;
        --no-security)
            RUN_SECURITY=false
            shift
            ;;
        --no-quality)
            RUN_QUALITY=false
            shift
            ;;
        --clean)
            cleanup_reports
            exit 0
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

# Setup - work in current directory
mkdir -p "$OUTPUT_DIR"

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not a git repository: $(pwd)"
    exit 1
fi

# Detect project types (supports monorepos with multiple languages)
detect_project_types() {
    local types=()

    # Check for Flutter/Dart
    if [[ -f "pubspec.yaml" ]] || find . -name "pubspec.yaml" -type f 2>/dev/null | head -1 | grep -q .; then
        types+=("flutter")
        RUN_FLUTTER=true
    fi

    # Check for Ruby
    if [[ -f "Gemfile" ]] || find . -name "Gemfile" -type f 2>/dev/null | head -1 | grep -q .; then
        types+=("ruby")
        RUN_RUBY=true
    fi

    # Check for JavaScript/TypeScript
    if [[ -f "package.json" ]] || find . -name "package.json" -type f 2>/dev/null | head -1 | grep -q .; then
        types+=("javascript")
        RUN_JS=true
    fi

    # Check for Python
    if [[ -f "requirements.txt" ]] || [[ -f "pyproject.toml" ]] || find . -name "requirements.txt" -o -name "pyproject.toml" -type f 2>/dev/null | head -1 | grep -q .; then
        types+=("python")
    fi

    # Check for PHP
    if [[ -f "composer.json" ]] || find . -name "composer.json" -type f 2>/dev/null | head -1 | grep -q .; then
        types+=("php")
    fi

    # Check for Java
    if [[ -f "pom.xml" ]] || [[ -f "build.gradle" ]] || find . -name "pom.xml" -o -name "build.gradle" -type f 2>/dev/null | head -1 | grep -q .; then
        types+=("java")
    fi

    if [[ ${#types[@]} -eq 0 ]]; then
        echo "unknown"
    elif [[ ${#types[@]} -eq 1 ]]; then
        echo "${types[0]}"
    else
        # Monorepo with multiple languages
        echo "monorepo ($(IFS=,; echo "${types[*]}"))"
    fi
}

PROJECT_TYPE=$(detect_project_types)

# Main execution
main() {
    print_header "🚀 STARTING COMPREHENSIVE CODE REVIEW"
    print_header "======================================"
    
    print_status "Repository: $(pwd)"
    print_status "Target branch: $TARGET_BRANCH"
    print_status "Project type: $PROJECT_TYPE"
    print_status "Output directory: $OUTPUT_DIR"
    print_status "Timestamp: $TIMESTAMP"
    echo ""

    # Create main report file
    MAIN_REPORT="$OUTPUT_DIR/review_report_$TIMESTAMP.txt"
    
    {
        echo "COMPREHENSIVE CODE REVIEW REPORT"
        echo "================================"
        echo "Generated: $(date)"
        echo "Repository: $(pwd)"
        echo "Target Branch: $TARGET_BRANCH"
        echo "Project Type: $PROJECT_TYPE"
        echo ""
    } > "$MAIN_REPORT"

    # 1. Basic PR Analysis
    print_header "📊 RUNNING BASIC PR ANALYSIS in $SCRIPT_DIR"
    if [[ -f "$SCRIPT_DIR/pr_analyzer.sh" ]]; then
        bash "$SCRIPT_DIR/pr_analyzer.sh" -t "$TARGET_BRANCH" >> "$MAIN_REPORT" 2>&1
        print_success "Basic PR analysis completed"
    else
        print_warning "PR analyzer script not found"
    fi
    echo ""

    # 2. Python-based Quick Review
    print_header "🔍 RUNNING DETAILED CODE ANALYSIS"
    if [[ -f "$SCRIPT_DIR/quick_review.py" ]] && command -v python3 &> /dev/null; then
        PYTHON_REPORT="$OUTPUT_DIR/python_analysis_$TIMESTAMP.json"

        # Build command with optional --full-scan flag
        PYTHON_CMD="python3 \"$SCRIPT_DIR/quick_review.py\" -t \"$TARGET_BRANCH\" --json"
        if [[ "$FULL_SCAN" == true ]]; then
            PYTHON_CMD="$PYTHON_CMD --full-scan"
        fi

        eval "$PYTHON_CMD" > "$PYTHON_REPORT" 2>&1

        # Add summary to main report
        {
            echo ""
            echo "DETAILED CODE ANALYSIS SUMMARY"
            echo "=============================="
            SUMMARY_CMD="python3 \"$SCRIPT_DIR/quick_review.py\" -t \"$TARGET_BRANCH\""
            if [[ "$FULL_SCAN" == true ]]; then
                SUMMARY_CMD="$SUMMARY_CMD --full-scan"
            fi
            eval "$SUMMARY_CMD"
        } >> "$MAIN_REPORT" 2>&1

        print_success "Detailed analysis completed - JSON report: $PYTHON_REPORT"
    else
        print_warning "Python quick review not available"
    fi
    echo ""

    # 3. Flutter-specific Analysis
    if [[ $RUN_FLUTTER == true ]]; then
        print_header "🎯 RUNNING FLUTTER-SPECIFIC ANALYSIS"
        if [[ -f "$SCRIPT_DIR/flutter_review.dart" ]] && command -v dart &> /dev/null; then
            FLUTTER_REPORT="$OUTPUT_DIR/flutter_analysis_$TIMESTAMP.txt"
            dart "$SCRIPT_DIR/flutter_review.dart" "." > "$FLUTTER_REPORT" 2>&1

            # Add to main report
            {
                echo ""
                echo "FLUTTER-SPECIFIC ANALYSIS"
                echo "========================="
                cat "$FLUTTER_REPORT"
            } >> "$MAIN_REPORT"

            print_success "Flutter analysis completed - Report: $FLUTTER_REPORT"
        else
            print_warning "Flutter analysis not available (dart command not found)"
        fi
        echo ""
    fi

    # 4. Ruby-specific Analysis
    if [[ "$RUN_RUBY" == true ]]; then
        print_header "🔴 RUNNING RUBY ANALYSIS"
        if [[ -f "$SCRIPT_DIR/ruby_review.rb" ]] && command -v ruby &> /dev/null; then
            RUBY_REPORT="$OUTPUT_DIR/ruby_analysis_$TIMESTAMP.txt"
            ruby "$SCRIPT_DIR/ruby_review.rb" "$REPO_PATH" "$TARGET_BRANCH" > "$RUBY_REPORT" 2>&1
            
            # Add summary to main report
            {
                echo ""
                echo "RUBY ANALYSIS SUMMARY"
                echo "===================="
                ruby "$SCRIPT_DIR/ruby_review.rb" "$REPO_PATH" "$TARGET_BRANCH"
            } >> "$MAIN_REPORT" 2>&1
            
            print_success "Ruby analysis completed - Report: $RUBY_REPORT"
        else
            print_warning "Ruby analysis not available"
        fi
        echo ""
    fi

    # 5. JavaScript/TypeScript Analysis
    if [[ "$RUN_JS" == true ]]; then
        print_header "🟨 RUNNING JAVASCRIPT/TYPESCRIPT ANALYSIS"
        if [[ -f "$SCRIPT_DIR/js_ts_review.js" ]] && command -v node &> /dev/null; then
            JS_REPORT="$OUTPUT_DIR/js_analysis_$TIMESTAMP.txt"
            node "$SCRIPT_DIR/js_ts_review.js" "$REPO_PATH" "$TARGET_BRANCH" > "$JS_REPORT" 2>&1
            
            # Add summary to main report
            {
                echo ""
                echo "JAVASCRIPT/TYPESCRIPT ANALYSIS SUMMARY"
                echo "======================================"
                node "$SCRIPT_DIR/js_ts_review.js" "$REPO_PATH" "$TARGET_BRANCH"
            } >> "$MAIN_REPORT" 2>&1
            
            print_success "JavaScript/TypeScript analysis completed - Report: $JS_REPORT"
        else
            print_warning "JavaScript/TypeScript analysis not available"
        fi
        echo ""
    fi

    # 6. Security Scan (if tools available)
    if [[ $RUN_SECURITY == true ]]; then
        print_header "🔒 RUNNING SECURITY ANALYSIS"
        SECURITY_REPORT="$OUTPUT_DIR/security_analysis_$TIMESTAMP.txt"
        
        {
            echo "SECURITY ANALYSIS REPORT"
            echo "======================="
            echo "Generated: $(date)"
            echo ""
            
            # Check for common security issues
            echo "🔍 Scanning for hardcoded secrets..."
            git diff "$TARGET_BRANCH"..HEAD | grep -i -E "(password|api_key|secret|token|private_key)" || echo "✅ No obvious secrets found"
            echo ""
            
            echo "🔍 Scanning for dangerous functions..."
            git diff "$TARGET_BRANCH"..HEAD | grep -i -E "(eval|exec|system|shell_exec)" || echo "✅ No dangerous functions found"
            echo ""
            
            # Check dependencies for known vulnerabilities (if tools available)
            if command -v npm &> /dev/null && [[ -f "package.json" ]]; then
                echo "🔍 Checking npm dependencies..."
                npm audit --audit-level=moderate 2>/dev/null || echo "⚠️  npm audit not available or no issues found"
            fi
            
            if command -v pip &> /dev/null && [[ -f "requirements.txt" ]]; then
                echo "🔍 Checking Python dependencies..."
                pip-audit 2>/dev/null || echo "⚠️  pip-audit not available"
            fi
            
        } > "$SECURITY_REPORT"
        
        # Add to main report
        {
            echo ""
            cat "$SECURITY_REPORT"
        } >> "$MAIN_REPORT"
        
        print_success "Security analysis completed - Report: $SECURITY_REPORT"
        echo ""
    fi

    # 7. Generate Summary
    print_header "📋 GENERATING FINAL SUMMARY"
    
    # Count issues from all reports
    TOTAL_FILES=$(git diff --name-only "$TARGET_BRANCH"..HEAD | wc -l)
    TOTAL_ADDITIONS=$(git diff --stat "$TARGET_BRANCH"..HEAD | tail -1 | grep -o '[0-9]\+ insertion' | grep -o '[0-9]\+' || echo "0")
    TOTAL_DELETIONS=$(git diff --stat "$TARGET_BRANCH"..HEAD | tail -1 | grep -o '[0-9]\+ deletion' | grep -o '[0-9]\+' || echo "0")
    
    SUMMARY_REPORT="$OUTPUT_DIR/summary_$TIMESTAMP.txt"
    {
        echo "REVIEW SUMMARY"
        echo "=============="
        echo "📁 Files changed: $TOTAL_FILES"
        echo "➕ Lines added: $TOTAL_ADDITIONS"
        echo "➖ Lines deleted: $TOTAL_DELETIONS"
        echo "🎯 Project type: $PROJECT_TYPE"
        echo ""
        echo "📊 REPORTS GENERATED:"
        echo "- Main report: $MAIN_REPORT"
        [[ -f "$PYTHON_REPORT" ]] && echo "- Detailed analysis: $PYTHON_REPORT"
        [[ -f "$FLUTTER_REPORT" ]] && echo "- Flutter analysis: $FLUTTER_REPORT"
        [[ -f "$RUBY_REPORT" ]] && echo "- Ruby analysis: $RUBY_REPORT"
        [[ -f "$JS_REPORT" ]] && echo "- JavaScript/TypeScript analysis: $JS_REPORT"
        [[ -f "$SECURITY_REPORT" ]] && echo "- Security analysis: $SECURITY_REPORT"
        echo ""
        echo "🎉 Review completed at: $(date)"
    } > "$SUMMARY_REPORT"
    
    # Display summary
    cat "$SUMMARY_REPORT"
    
    print_success "All reports saved to: $OUTPUT_DIR"
    print_success "Main report: $MAIN_REPORT"
    
    # Open main report if on macOS/Linux with GUI
    if [[ "$OSTYPE" == "darwin"* ]] && command -v open &> /dev/null; then
        open "$MAIN_REPORT"
    elif command -v xdg-open &> /dev/null; then
        xdg-open "$MAIN_REPORT" 2>/dev/null &
    fi
}

# Run main function
main

print_header "✅ CODE REVIEW AUTOMATION COMPLETE!"
