#!/usr/bin/env node

/**
 * JavaScript/TypeScript code review automation
 * Analyzes JS/TS code for common issues and best practices
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

class JSReviewAnalyzer {
  constructor(repoPath = '.', targetBranch = 'main') {
    this.repoPath = repoPath;
    this.targetBranch = targetBranch;
    this.issues = [];
  }

  async analyze() {
    console.log('🔍 Starting JavaScript/TypeScript code analysis...');
    
    const changedFiles = this.getChangedFiles();
    const jsFiles = changedFiles.filter(f => 
      f.endsWith('.js') || f.endsWith('.ts') || f.endsWith('.jsx') || f.endsWith('.tsx')
    );
    
    if (jsFiles.length === 0) {
      console.log('No JavaScript/TypeScript files found to analyze', {jsFiles, changedFiles});
      return {
        summary: { total_files: 0, js_files: 0, issues: 0 },
        issues: []
      };
    }

    console.log(`📁 Found ${jsFiles.length} JavaScript/TypeScript files to analyze`);

    for (const file of jsFiles) {
      await this.analyzeFile(file);
    }

    return {
      summary: {
        total_files: changedFiles.length,
        js_files: jsFiles.length,
        issues: this.issues.length,
        high_severity: this.issues.filter(i => i.severity === 'high').length,
        medium_severity: this.issues.filter(i => i.severity === 'medium').length,
        low_severity: this.issues.filter(i => i.severity === 'low').length,
      },
      issues: this.issues,
      files_analyzed: jsFiles,
    };
  }

  getChangedFiles() {
    try {
      const result = execSync(`git diff --name-only ${this.targetBranch}..HEAD`, { 
        encoding: 'utf8',
        cwd: this.repoPath 
      });
      return result.trim().split('\n').filter(line => line.length > 0);
    } catch (error) {
      console.warn('⚠️  Warning: Could not get git diff, analyzing all JS/TS files');
      return this.getAllJSFiles();
    }
  }

  getAllJSFiles() {
    const files = [];
    const extensions = ['.js', '.ts', '.jsx', '.tsx'];
    
    const walkDir = (dir) => {
      const items = fs.readdirSync(dir);
      for (const item of items) {
        const fullPath = path.join(dir, item);
        const stat = fs.statSync(fullPath);
        
        if (stat.isDirectory() && !item.startsWith('.') && item !== 'node_modules') {
          walkDir(fullPath);
        } else if (stat.isFile() && extensions.some(ext => item.endsWith(ext))) {
          files.push(path.relative(this.repoPath, fullPath));
        }
      }
    };
    
    walkDir(this.repoPath);
    return files;
  }

  async analyzeFile(filePath) {
    const fullPath = path.join(this.repoPath, filePath);
    
    if (!fs.existsSync(fullPath)) {
      return;
    }

    const content = fs.readFileSync(fullPath, 'utf8');
    const lines = content.split('\n');

    // Run various checks
    this.checkSecurityIssues(filePath, content, lines);
    this.checkPerformanceIssues(filePath, content, lines);
    this.checkCodeQuality(filePath, content, lines);
    this.checkTypeScriptPatterns(filePath, content, lines);
    this.checkReactPatterns(filePath, content, lines);
    this.checkAsyncPatterns(filePath, content, lines);
    this.checkErrorHandling(filePath, content, lines);
    this.checkTestingPatterns(filePath, content, lines);
  }

  checkSecurityIssues(file, content, lines) {
    // XSS vulnerabilities
    if (content.includes('innerHTML =') || content.includes('outerHTML =')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /innerHTML\s*=|outerHTML\s*=/),
        type: 'security',
        severity: 'high',
        message: 'Direct HTML assignment - XSS risk',
        suggestion: 'Use textContent, createElement, or sanitize HTML input'
      });
    }

    // Dangerous eval usage
    if (content.includes('eval(') || content.includes('Function(')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /eval\(|Function\(/),
        type: 'security',
        severity: 'high',
        message: 'Use of eval() or Function() - security risk',
        suggestion: 'Avoid eval() and Function() constructor'
      });
    }

    // Hardcoded secrets
    const secretPatterns = [
      /password\s*[:=]\s*['"][^'"]+['"]/i,
      /api_key\s*[:=]\s*['"][^'"]+['"]/i,
      /secret\s*[:=]\s*['"][^'"]+['"]/i,
      /token\s*[:=]\s*['"][^'"]+['"]/i
    ];

    secretPatterns.forEach(pattern => {
      if (pattern.test(content)) {
        this.addIssue({
          file,
          line: this.findLineNumber(content, pattern),
          type: 'security',
          severity: 'high',
          message: 'Hardcoded secret detected',
          suggestion: 'Use environment variables or secure configuration'
        });
      }
    });

    // Unsafe regex patterns
    if (content.includes('RegExp(') && content.includes('+')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /RegExp\(/),
        type: 'security',
        severity: 'medium',
        message: 'Dynamic regex construction - ReDoS risk',
        suggestion: 'Use static regex patterns when possible'
      });
    }
  }

  checkPerformanceIssues(file, content, lines) {
    // Inefficient DOM queries
    if (content.includes('document.getElementById') && content.split('document.getElementById').length > 3) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /document\.getElementById/),
        type: 'performance',
        severity: 'medium',
        message: 'Multiple DOM queries - consider caching',
        suggestion: 'Cache DOM elements in variables'
      });
    }

    // Inefficient array operations
    if (content.includes('.forEach(') && content.includes('.push(')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /\.forEach.*\.push/),
        type: 'performance',
        severity: 'medium',
        message: 'forEach with push - consider using map',
        suggestion: 'Use .map() instead of .forEach() with .push()'
      });
    }

    // Memory leaks - event listeners
    if (content.includes('addEventListener') && !content.includes('removeEventListener')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /addEventListener/),
        type: 'performance',
        severity: 'medium',
        message: 'Event listener without cleanup',
        suggestion: 'Add removeEventListener in cleanup/unmount'
      });
    }

    // Large bundle imports
    if (content.includes('import *') || content.includes('require(')) {
      const matches = content.match(/import \* as .* from ['"][^'"]+['"]|require\(['"][^'"]+['"]\)/g);
      if (matches && matches.length > 0) {
        this.addIssue({
          file,
          line: this.findLineNumber(content, /import \*|require\(/),
          type: 'performance',
          severity: 'low',
          message: 'Wildcard import detected',
          suggestion: 'Import only needed functions to reduce bundle size'
        });
      }
    }
  }

  checkCodeQuality(file, content, lines) {
    // Console statements
    if (content.includes('console.log') || content.includes('console.warn')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /console\.(log|warn|error|debug)/),
        type: 'quality',
        severity: 'low',
        message: 'Console statement found',
        suggestion: 'Remove console statements or use proper logging'
      });
    }

    // Debugger statements
    if (content.includes('debugger;')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /debugger;/),
        type: 'quality',
        severity: 'medium',
        message: 'Debugger statement found',
        suggestion: 'Remove debugger statements before production'
      });
    }

    // TODO/FIXME comments
    lines.forEach((line, index) => {
      if (/TODO|FIXME|HACK/i.test(line)) {
        this.addIssue({
          file,
          line: index + 1,
          type: 'quality',
          severity: 'low',
          message: 'TODO/FIXME comment found',
          suggestion: 'Address the TODO or create a proper issue'
        });
      }
    });

    // Large functions
    const functionMatches = content.match(/function\s+\w+|const\s+\w+\s*=\s*\(/g);
    if (functionMatches) {
      functionMatches.forEach(match => {
        const functionStart = content.indexOf(match);
        const functionContent = this.extractFunctionContent(content, functionStart);
        if (functionContent.split('\n').length > 50) {
          this.addIssue({
            file,
            line: this.findLineNumber(content, new RegExp(this.escapeRegex(match))),
            type: 'quality',
            severity: 'medium',
            message: `Large function detected (${functionContent.split('\n').length} lines)`,
            suggestion: 'Consider breaking down into smaller functions'
          });
        }
      });
    }

    // Magic numbers
    const magicNumbers = content.match(/\b(?!0|1|2|10|100|1000)\d{3,}\b/g);
    if (magicNumbers) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /\b(?!0|1|2|10|100|1000)\d{3,}\b/),
        type: 'quality',
        severity: 'low',
        message: 'Magic numbers detected',
        suggestion: 'Use named constants for magic numbers'
      });
    }
  }

  checkTypeScriptPatterns(file, content, lines) {
    if (!file.endsWith('.ts') && !file.endsWith('.tsx')) return;

    // Any type usage
    if (content.includes(': any') || content.includes('<any>')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /:\s*any|<any>/),
        type: 'typescript',
        severity: 'medium',
        message: 'Use of "any" type',
        suggestion: 'Use specific types instead of "any"'
      });
    }

    // Missing return types
    const functionRegex = /function\s+\w+\s*\([^)]*\)\s*\{|const\s+\w+\s*=\s*\([^)]*\)\s*=>/g;
    let match;
    while ((match = functionRegex.exec(content)) !== null) {
      const beforeFunction = content.substring(0, match.index);
      const afterFunction = content.substring(match.index + match[0].length);
      
      if (!beforeFunction.includes(': ') && !afterFunction.startsWith(': ')) {
        this.addIssue({
          file,
          line: this.findLineNumber(content, match[0]),
          type: 'typescript',
          severity: 'low',
          message: 'Function without explicit return type',
          suggestion: 'Add explicit return type annotations'
        });
      }
    }

    // Non-null assertion without null check
    if (content.includes('!.') || content.includes('!;')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /!\.|!;/),
        type: 'typescript',
        severity: 'medium',
        message: 'Non-null assertion operator used',
        suggestion: 'Consider proper null checking instead of ! operator'
      });
    }
  }

  checkReactPatterns(file, content, lines) {
    if (!file.endsWith('.jsx') && !file.endsWith('.tsx')) return;

    // Missing key prop in lists
    if (content.includes('.map(') && !content.includes('key=')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /\.map\(/),
        type: 'react',
        severity: 'medium',
        message: 'Missing key prop in mapped elements',
        suggestion: 'Add unique key prop to mapped elements'
      });
    }

    // Direct state mutation
    if (content.includes('state.') && content.includes(' = ')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /state\.\w+\s*=/),
        type: 'react',
        severity: 'high',
        message: 'Direct state mutation detected',
        suggestion: 'Use setState or state setter functions'
      });
    }

    // Missing dependency array in useEffect
    if (content.includes('useEffect(') && !content.includes(', [')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /useEffect\(/),
        type: 'react',
        severity: 'medium',
        message: 'useEffect without dependency array',
        suggestion: 'Add dependency array to useEffect'
      });
    }

    // Large components
    if (content.includes('function ') || content.includes('const ')) {
      const componentMatch = content.match(/(?:function\s+\w+|const\s+\w+\s*=)/);
      if (componentMatch && lines.length > 200) {
        this.addIssue({
          file,
          line: 1,
          type: 'react',
          severity: 'medium',
          message: `Large component (${lines.length} lines)`,
          suggestion: 'Consider breaking down into smaller components'
        });
      }
    }
  }

  checkAsyncPatterns(file, content, lines) {
    // Async without await
    if (content.includes('async ') && !content.includes('await ')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /async\s/),
        type: 'async',
        severity: 'low',
        message: 'Async function without await',
        suggestion: 'Remove async keyword if not using await'
      });
    }

    // Missing error handling for async operations
    if (content.includes('await ') && !content.includes('try') && !content.includes('.catch(')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /await\s/),
        type: 'async',
        severity: 'medium',
        message: 'Async operation without error handling',
        suggestion: 'Add try-catch or .catch() for async operations'
      });
    }

    // Promise constructor anti-pattern
    if (content.includes('new Promise') && content.includes('async')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /new Promise/),
        type: 'async',
        severity: 'medium',
        message: 'Promise constructor with async function',
        suggestion: 'Use async/await instead of Promise constructor'
      });
    }
  }

  checkErrorHandling(file, content, lines) {
    // Empty catch blocks
    if (content.includes('catch') && content.includes('catch () {}')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /catch\s*\(\s*\)\s*\{\s*\}/),
        type: 'error_handling',
        severity: 'medium',
        message: 'Empty catch block',
        suggestion: 'Handle errors appropriately or log them'
      });
    }

    // Generic error catching
    if (content.includes('catch (e)') || content.includes('catch (error)')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /catch\s*\(\s*(e|error)\s*\)/),
        type: 'error_handling',
        severity: 'low',
        message: 'Generic error catching',
        suggestion: 'Consider catching specific error types when possible'
      });
    }
  }

  checkTestingPatterns(file, content, lines) {
    if (!file.includes('test') && !file.includes('spec')) return;

    // Missing test descriptions
    if (content.includes('it(') && content.includes('it(() => {')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /it\(\(\)\s*=>/),
        type: 'testing',
        severity: 'low',
        message: 'Test without description',
        suggestion: 'Add descriptive test names'
      });
    }

    // Large test files
    if (lines.length > 300) {
      this.addIssue({
        file,
        line: 1,
        type: 'testing',
        severity: 'medium',
        message: `Large test file (${lines.length} lines)`,
        suggestion: 'Consider splitting into multiple test files'
      });
    }

    // Missing assertions
    if (content.includes('it(') && !content.includes('expect(') && !content.includes('assert')) {
      this.addIssue({
        file,
        line: this.findLineNumber(content, /it\(/),
        type: 'testing',
        severity: 'medium',
        message: 'Test without assertions',
        suggestion: 'Add proper assertions to verify behavior'
      });
    }
  }

  addIssue({ file, line, type, severity, message, suggestion }) {
    this.issues.push({
      file,
      line,
      type,
      severity,
      message,
      suggestion
    });
  }

  findLineNumber(content, pattern) {
    const lines = content.split('\n');
    for (let i = 0; i < lines.length; i++) {
      try {
        if (lines[i].match(pattern)) {
          return i + 1;
        }
      } catch (e) {
        // Skip invalid regex patterns
        continue;
      }
    }
    return 1;
  }

  escapeRegex(string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  }

  extractFunctionContent(content, start) {
    let braceCount = 0;
    let i = start;
    
    // Find the opening brace
    while (i < content.length && content[i] !== '{') {
      i++;
    }
    
    if (i >= content.length) return '';
    
    const functionStart = i;
    braceCount = 1;
    i++;
    
    // Find the matching closing brace
    while (i < content.length && braceCount > 0) {
      if (content[i] === '{') {
        braceCount++;
      } else if (content[i] === '}') {
        braceCount--;
      }
      i++;
    }
    
    return content.substring(functionStart, i);
  }

  printReport(report) {
    const { summary, issues } = report;

    console.log('\n' + '='.repeat(60));
    console.log('🚀 JAVASCRIPT/TYPESCRIPT CODE REVIEW REPORT');
    console.log('='.repeat(60));
    console.log(`📁 Total files: ${summary.total_files}`);
    console.log(`🚀 JS/TS files analyzed: ${summary.js_files}`);
    console.log(`🚨 Total issues: ${summary.issues}`);
    console.log(`🔴 High severity: ${summary.high_severity}`);
    console.log(`🟡 Medium severity: ${summary.medium_severity}`);
    console.log(`🟢 Low severity: ${summary.low_severity}`);

    if (issues.length > 0) {
      console.log('\n📋 ISSUES FOUND:');
      console.log('-'.repeat(40));

      const groupedIssues = {};
      issues.forEach(issue => {
        if (!groupedIssues[issue.type]) {
          groupedIssues[issue.type] = [];
        }
        groupedIssues[issue.type].push(issue);
      });

      Object.entries(groupedIssues).forEach(([type, typeIssues]) => {
        console.log(`\n${this.getTypeIcon(type)} ${type.toUpperCase().replace(/_/g, ' ')}:`);
        typeIssues.forEach(issue => {
          const severityIcon = issue.severity === 'high' ? '🔴' : 
                              issue.severity === 'medium' ? '🟡' : '🟢';
          console.log(`  ${severityIcon} ${issue.file}:${issue.line} - ${issue.message}`);
          console.log(`    💡 ${issue.suggestion}`);
        });
      });
    } else {
      console.log('\n✅ No issues found! Great job!');
    }
  }

  getTypeIcon(type) {
    const icons = {
      'security': '🔒',
      'performance': '⚡',
      'quality': '🎯',
      'typescript': '🔷',
      'react': '⚛️',
      'async': '⏱️',
      'error_handling': '🛡️',
      'testing': '🧪'
    };
    return icons[type] || '📝';
  }
}

// Main execution
if (require.main === module) {
  const repoPath = process.argv[2] || '.';
  const targetBranch = process.argv[3] || 'main';
  
  const analyzer = new JSReviewAnalyzer(repoPath, targetBranch);
  analyzer.analyze().then(report => {
    analyzer.printReport(report);
    
    // Output JSON if requested
    if (process.argv.includes('--json')) {
      console.log('\n' + JSON.stringify(report, null, 2));
    }
  }).catch(error => {
    console.error('Error during analysis:', error);
    process.exit(1);
  });
}
