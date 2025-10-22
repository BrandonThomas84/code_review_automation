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
