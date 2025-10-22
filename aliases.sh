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
