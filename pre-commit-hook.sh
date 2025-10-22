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
