#!/usr/bin/env python3
"""
Quick Code Review Automation Script
Performs rapid analysis of code changes for common issues
"""

import os
import sys
import subprocess
import json
import re
from pathlib import Path
from typing import List, Dict, Any
import argparse

class CodeReviewAnalyzer:
    def __init__(self, repo_path: str = "."):
        self.repo_path = Path(repo_path)
        self.issues = []
        
    def analyze_git_diff(self, target_branch: str) -> Dict[str, Any]:
        """Analyze git diff for common issues"""
        try:
            # First, try to fetch the target branch to ensure it's up to date
            subprocess.run(
                ["git", "fetch", "origin", target_branch],
                capture_output=True, text=True, cwd=self.repo_path
            )

            # Get list of changed files between target branch and current branch
            result = subprocess.run(
                ["git", "diff", "--name-only", f"origin/{target_branch}...HEAD"],
                capture_output=True, text=True, cwd=self.repo_path
            )

            if result.returncode != 0:
                # Fallback: try without origin prefix
                result = subprocess.run(
                    ["git", "diff", "--name-only", f"{target_branch}...HEAD"],
                    capture_output=True, text=True, cwd=self.repo_path
                )

            changed_files = [f for f in result.stdout.strip().split('\n') if f]

            # Get detailed diff
            result = subprocess.run(
                ["git", "diff", f"origin/{target_branch}...HEAD"],
                capture_output=True, text=True, cwd=self.repo_path
            )

            if result.returncode != 0:
                # Fallback: try without origin prefix
                result = subprocess.run(
                    ["git", "diff", f"{target_branch}...HEAD"],
                    capture_output=True, text=True, cwd=self.repo_path
                )

            diff_content = result.stdout

            return {
                "changed_files": changed_files,
                "diff_content": diff_content,
                "file_count": len(changed_files)
            }
        except Exception as e:
            return {"error": str(e)}

    def analyze_full_codebase(self) -> Dict[str, Any]:
        """Analyze entire codebase for issues"""
        try:
            # Get all source files
            code_extensions = ['.py', '.js', '.ts', '.jsx', '.tsx', '.dart', '.rb', '.php', '.java', '.kt']
            all_files = []

            for ext in code_extensions:
                result = subprocess.run(
                    ["find", ".", "-name", f"*{ext}", "-type", "f"],
                    capture_output=True, text=True, cwd=self.repo_path
                )
                if result.stdout.strip():
                    all_files.extend(result.stdout.strip().split('\n'))

            # Read all files content
            full_content = ""
            for file_path in all_files:
                try:
                    with open(self.repo_path / file_path, 'r', encoding='utf-8', errors='ignore') as f:
                        full_content += f.read() + "\n"
                except:
                    pass

            return {
                "changed_files": all_files,
                "diff_content": full_content,
                "file_count": len(all_files)
            }
        except Exception as e:
            return {"error": str(e)}
    
    def check_security_issues(self, diff_content: str) -> List[Dict[str, str]]:
        """Check for common security issues"""
        security_patterns = [
            (r'password\s*=\s*["\'][^"\']+["\']', "Hardcoded password detected"),
            (r'api_key\s*=\s*["\'][^"\']+["\']', "Hardcoded API key detected"),
            (r'secret\s*=\s*["\'][^"\']+["\']', "Hardcoded secret detected"),
            (r'token\s*=\s*["\'][^"\']+["\']', "Hardcoded token detected"),
            (r'eval\s*\(', "Use of eval() function - security risk"),
            (r'exec\s*\(', "Use of exec() function - security risk"),
            (r'\.innerHTML\s*=', "Direct innerHTML assignment - XSS risk"),
            (r'document\.write\s*\(', "Use of document.write - XSS risk"),
        ]
        
        issues = []
        for pattern, message in security_patterns:
            matches = re.finditer(pattern, diff_content, re.IGNORECASE)
            for match in matches:
                issues.append({
                    "type": "security",
                    "severity": "high",
                    "message": message,
                    "pattern": pattern,
                    "line_content": match.group(0)
                })
        
        return issues
    
    def check_code_quality(self, diff_content: str) -> List[Dict[str, str]]:
        """Check for code quality issues"""
        quality_patterns = [
            (r'console\.log\s*\(', "Console.log statement found - remove before production"),
            (r'print\s*\(', "Print statement found - consider using proper logging"),
            (r'debugger;', "Debugger statement found - remove before production"),
            (r'TODO|FIXME|HACK', "TODO/FIXME/HACK comment found"),
            (r'\.catch\s*\(\s*\)', "Empty catch block - handle errors properly"),
            (r'function\s+\w+\s*\([^)]*\)\s*\{[^}]{200,}', "Large function detected - consider breaking down"),
            (r'if\s*\([^)]+\)\s*\{[^}]*if\s*\([^)]+\)\s*\{[^}]*if', "Deep nesting detected - consider refactoring"),
        ]
        
        issues = []
        for pattern, message in quality_patterns:
            matches = re.finditer(pattern, diff_content, re.IGNORECASE)
            for match in matches:
                issues.append({
                    "type": "quality",
                    "severity": "medium",
                    "message": message,
                    "pattern": pattern,
                    "line_content": match.group(0)
                })
        
        return issues
    
    def check_flutter_specific(self, changed_files: List[str], diff_content: str) -> List[Dict[str, str]]:
        """Flutter/Dart specific checks"""
        dart_files = [f for f in changed_files if f.endswith('.dart')]
        if not dart_files:
            return []
        
        flutter_patterns = [
            (r'setState\s*\(\s*\(\s*\)\s*\{\s*\}\s*\)', "Empty setState call"),
            (r'build\s*\([^)]*\)\s*\{[^}]{500,}', "Large build method - consider extracting widgets"),
            (r'Container\s*\(\s*child:\s*Container', "Nested Container widgets - consider simplifying"),
            (r'Column\s*\([^)]*children:\s*\[\s*\]', "Empty Column widget"),
            (r'Row\s*\([^)]*children:\s*\[\s*\]', "Empty Row widget"),
            (r'\.of\(context\)(?!\.)', "Missing null safety check for context"),
        ]
        
        issues = []
        for pattern, message in flutter_patterns:
            matches = re.finditer(pattern, diff_content, re.IGNORECASE)
            for match in matches:
                issues.append({
                    "type": "flutter",
                    "severity": "medium",
                    "message": message,
                    "pattern": pattern,
                    "line_content": match.group(0)
                })
        
        return issues
    
    def check_test_coverage(self, changed_files: List[str]) -> List[Dict[str, str]]:
        """Check if tests are included for new features"""
        code_files = [f for f in changed_files if f.endswith(('.dart', '.js', '.ts', '.py', '.php'))]
        test_files = [f for f in changed_files if 'test' in f.lower() or 'spec' in f.lower()]
        
        issues = []
        if code_files and not test_files:
            issues.append({
                "type": "testing",
                "severity": "medium",
                "message": f"Code changes detected but no test files modified. Consider adding tests.",
                "files": code_files
            })
        
        return issues
    
    def generate_report(self, target_branch: str, full_scan: bool = False) -> Dict[str, Any]:
        """Generate comprehensive review report"""
        if full_scan:
            print("🔍 Scanning entire codebase...")
            git_analysis = self.analyze_full_codebase()
        else:
            print(f"🔍 Analyzing changes against {target_branch}...")
            git_analysis = self.analyze_git_diff(target_branch)

        if "error" in git_analysis:
            return {"error": git_analysis["error"]}

        changed_files = git_analysis["changed_files"]
        diff_content = git_analysis["diff_content"]

        print(f"📁 Found {len(changed_files)} files to analyze")
        
        # Run all checks
        security_issues = self.check_security_issues(diff_content)
        quality_issues = self.check_code_quality(diff_content)
        flutter_issues = self.check_flutter_specific(changed_files, diff_content)
        test_issues = self.check_test_coverage(changed_files)
        
        all_issues = security_issues + quality_issues + flutter_issues + test_issues
        
        # Categorize by severity
        high_severity = [i for i in all_issues if i.get("severity") == "high"]
        medium_severity = [i for i in all_issues if i.get("severity") == "medium"]
        low_severity = [i for i in all_issues if i.get("severity") == "low"]
        
        return {
            "summary": {
                "total_files": len(changed_files),
                "total_issues": len(all_issues),
                "high_severity": len(high_severity),
                "medium_severity": len(medium_severity),
                "low_severity": len(low_severity)
            },
            "changed_files": changed_files,
            "issues": {
                "high": high_severity,
                "medium": medium_severity,
                "low": low_severity
            }
        }
    
    def print_report(self, report: Dict[str, Any]):
        """Print formatted report"""
        if "error" in report:
            print(f"❌ Error: {report['error']}")
            return
        
        summary = report["summary"]
        
        print("\n" + "="*60)
        print("📋 CODE REVIEW SUMMARY")
        print("="*60)
        print(f"📁 Files changed: {summary['total_files']}")
        print(f"🚨 Total issues: {summary['total_issues']}")
        print(f"🔴 High severity: {summary['high_severity']}")
        print(f"🟡 Medium severity: {summary['medium_severity']}")
        print(f"🟢 Low severity: {summary['low_severity']}")
        
        # Print issues by severity
        for severity, color in [("high", "🔴"), ("medium", "🟡"), ("low", "🟢")]:
            issues = report["issues"][severity]
            if issues:
                print(f"\n{color} {severity.upper()} SEVERITY ISSUES:")
                print("-" * 40)
                for i, issue in enumerate(issues, 1):
                    print(f"{i}. {issue['message']}")
                    if 'line_content' in issue:
                        print(f"   Code: {issue['line_content'][:100]}...")
                    print()

def main():
    parser = argparse.ArgumentParser(description="Quick Code Review Automation")
    parser.add_argument("--target", "-t", required=True, help="Target branch to compare against (required)")
    parser.add_argument("--json", "-j", action="store_true", help="Output as JSON")
    parser.add_argument("--full-scan", action="store_true", help="Scan entire codebase instead of just changed files")

    args = parser.parse_args()

    analyzer = CodeReviewAnalyzer(".")  # Use current directory
    report = analyzer.generate_report(args.target, full_scan=args.full_scan)

    if args.json:
        print(json.dumps(report, indent=2))
    else:
        analyzer.print_report(report)

if __name__ == "__main__":
    main()
