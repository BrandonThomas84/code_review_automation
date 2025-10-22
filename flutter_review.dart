#!/usr/bin/env dart

/// Flutter-specific code review automation
/// Analyzes Dart/Flutter code for common issues and best practices

import 'dart:io';
import 'dart:convert';

class FlutterReviewAnalyzer {
  final String repoPath;
  final String targetBranch;
  final List<Map<String, dynamic>> _issues = [];

  FlutterReviewAnalyzer(this.repoPath, [this.targetBranch = 'main']);

  /// Main analysis entry point
  Future<Map<String, dynamic>> analyze() async {
    print('🔍 Starting Flutter code analysis...');

    final changedFiles = await _getChangedFiles();
    final dartFiles = changedFiles.where((f) => f.endsWith('.dart')).toList();

    if (dartFiles.isEmpty) {
      return {
        'summary': {'total_files': 0, 'dart_files': 0, 'issues': 0},
        'issues': [],
      };
    }

    print('📁 Found ${dartFiles.length} Dart files to analyze');

    for (final file in dartFiles) {
      await _analyzeFile(file);
    }

    return {
      'summary': {
        'total_files': changedFiles.length,
        'dart_files': dartFiles.length,
        'issues': _issues.length,
        'high_severity': _issues.where((i) => i['severity'] == 'high').length,
        'medium_severity': _issues
            .where((i) => i['severity'] == 'medium')
            .length,
        'low_severity': _issues.where((i) => i['severity'] == 'low').length,
      },
      'issues': _issues,
      'files_analyzed': dartFiles,
    };
  }

  /// Get list of changed files from git
  Future<List<String>> _getChangedFiles() async {
    try {
      final result = await Process.run('git', [
        'diff',
        '--name-only',
        '$targetBranch..HEAD',
      ], workingDirectory: repoPath);

      if (result.exitCode != 0) {
        print('⚠️  Warning: Could not get git diff, analyzing all Dart files');
        return await _getAllDartFiles();
      }

      return result.stdout
          .toString()
          .trim()
          .split('\n')
          .where((line) => line.isNotEmpty)
          .toList();
    } catch (e) {
      print('⚠️  Warning: Git not available, analyzing all Dart files');
      return await _getAllDartFiles();
    }
  }

  /// Get all Dart files in the project
  Future<List<String>> _getAllDartFiles() async {
    final files = <String>[];
    final libDir = Directory('$repoPath/lib');

    if (await libDir.exists()) {
      await for (final entity in libDir.list(recursive: true)) {
        if (entity is File && entity.path.endsWith('.dart')) {
          files.add(entity.path.replaceFirst('$repoPath/', ''));
        }
      }
    }

    return files;
  }

  /// Analyze a single Dart file
  Future<void> _analyzeFile(String filePath) async {
    final file = File('$repoPath/$filePath');

    if (!await file.exists()) {
      return;
    }

    final content = await file.readAsString();
    final lines = content.split('\n');

    // Run various checks
    _checkWidgetStructure(filePath, content, lines);
    _checkStateManagement(filePath, content, lines);
    _checkPerformance(filePath, content, lines);
    _checkCodeStyle(filePath, content, lines);
    _checkNullSafety(filePath, content, lines);
    _checkAsyncPatterns(filePath, content, lines);
  }

  /// Check widget structure and composition
  void _checkWidgetStructure(String file, String content, List<String> lines) {
    // Large build methods
    final buildMethodRegex = RegExp(r'Widget\s+build\s*\([^)]*\)\s*\{');
    final matches = buildMethodRegex.allMatches(content);

    for (final match in matches) {
      final startLine = content.substring(0, match.start).split('\n').length;
      final methodContent = _extractMethodContent(content, match.start);

      if (methodContent.split('\n').length > 50) {
        _addIssue(
          file: file,
          line: startLine,
          type: 'widget_structure',
          severity: 'medium',
          message:
              'Large build method detected (${methodContent.split('\n').length} lines). Consider extracting widgets.',
          suggestion:
              'Break down into smaller widget methods or separate widget classes',
        );
      }
    }

    // Nested Container widgets
    if (content.contains(RegExp(r'Container\s*\([^)]*child:\s*Container'))) {
      _addIssue(
        file: file,
        line: _findLineNumber(content, 'Container'),
        type: 'widget_structure',
        severity: 'low',
        message: 'Nested Container widgets detected',
        suggestion:
            'Consider using Padding, Margin, or other specific widgets instead',
      );
    }

    // Empty widgets
    final emptyWidgetPatterns = [
      r'Column\s*\([^)]*children:\s*\[\s*\]',
      r'Row\s*\([^)]*children:\s*\[\s*\]',
      r'ListView\s*\([^)]*children:\s*\[\s*\]',
    ];

    for (final pattern in emptyWidgetPatterns) {
      if (content.contains(RegExp(pattern))) {
        _addIssue(
          file: file,
          line: _findLineNumber(content, pattern),
          type: 'widget_structure',
          severity: 'low',
          message: 'Empty widget detected: ${pattern.split(r'\s*\(')[0]}',
          suggestion: 'Remove empty widgets or add placeholder content',
        );
      }
    }
  }

  /// Check state management patterns
  void _checkStateManagement(String file, String content, List<String> lines) {
    // Empty setState calls
    if (content.contains(RegExp(r'setState\s*\(\s*\(\s*\)\s*\{\s*\}\s*\)'))) {
      _addIssue(
        file: file,
        line: _findLineNumber(content, 'setState'),
        type: 'state_management',
        severity: 'medium',
        message: 'Empty setState call detected',
        suggestion: 'Remove empty setState or add state changes',
      );
    }

    // setState in build method
    if (content.contains('build(') && content.contains('setState(')) {
      final buildStart = content.indexOf('build(');
      final setStateInBuild = content.indexOf('setState(', buildStart);
      if (setStateInBuild != -1) {
        _addIssue(
          file: file,
          line: _findLineNumber(content, 'setState'),
          type: 'state_management',
          severity: 'high',
          message: 'setState called in build method',
          suggestion: 'Move setState calls outside of build method',
        );
      }
    }

    // Missing dispose for controllers
    if (content.contains('Controller') && !content.contains('dispose()')) {
      _addIssue(
        file: file,
        line: 1,
        type: 'state_management',
        severity: 'medium',
        message: 'Controller detected but no dispose method found',
        suggestion: 'Implement dispose() method to clean up controllers',
      );
    }
  }

  /// Check performance-related issues
  void _checkPerformance(String file, String content, List<String> lines) {
    // Expensive operations in build
    final expensivePatterns = [
      r'DateTime\.now\(\)',
      r'Random\(\)',
      r'\.toList\(\)',
      r'\.where\(',
      r'\.map\(',
    ];

    for (final pattern in expensivePatterns) {
      if (content.contains('build(') && content.contains(RegExp(pattern))) {
        final buildStart = content.indexOf('build(');
        final operationInBuild = content.indexOf(RegExp(pattern), buildStart);
        if (operationInBuild != -1) {
          _addIssue(
            file: file,
            line: _findLineNumber(content, pattern),
            type: 'performance',
            severity: 'medium',
            message: 'Expensive operation in build method: $pattern',
            suggestion:
                'Move expensive operations outside build or use memoization',
          );
        }
      }
    }

    // Missing const constructors
    final widgetPatterns = [
      r'Text\s*\(',
      r'Icon\s*\(',
      r'Padding\s*\(',
      r'Container\s*\(',
    ];

    for (final pattern in widgetPatterns) {
      final matches = RegExp(pattern).allMatches(content);
      for (final match in matches) {
        final beforeMatch = content.substring(0, match.start);
        if (!beforeMatch.endsWith('const ')) {
          _addIssue(
            file: file,
            line: _findLineNumber(content, pattern),
            type: 'performance',
            severity: 'low',
            message: 'Missing const keyword for widget: ${match.group(0)}',
            suggestion: 'Add const keyword for better performance',
          );
        }
      }
    }
  }

  /// Check code style and conventions
  void _checkCodeStyle(String file, String content, List<String> lines) {
    // Check for proper naming conventions
    final classRegex = RegExp(r'class\s+([a-z][a-zA-Z0-9_]*)\s+');
    final matches = classRegex.allMatches(content);

    for (final match in matches) {
      final className = match.group(1)!;
      if (className[0].toLowerCase() == className[0]) {
        _addIssue(
          file: file,
          line: _findLineNumber(content, 'class $className'),
          type: 'code_style',
          severity: 'low',
          message: 'Class name should start with uppercase: $className',
          suggestion:
              'Rename to ${className[0].toUpperCase()}${className.substring(1)}',
        );
      }
    }

    // Check for TODO/FIXME comments
    for (int i = 0; i < lines.length; i++) {
      final line = lines[i];
      if (line.contains(
        RegExp(r'//\s*(TODO|FIXME|HACK)', caseSensitive: false),
      )) {
        _addIssue(
          file: file,
          line: i + 1,
          type: 'code_style',
          severity: 'low',
          message: 'TODO/FIXME comment found',
          suggestion: 'Address the TODO or create a proper issue',
        );
      }
    }
  }

  /// Check null safety patterns
  void _checkNullSafety(String file, String content, List<String> lines) {
    // Force unwrapping without null checks
    final forceUnwrapRegex = RegExp(r'[a-zA-Z_][a-zA-Z0-9_]*!(?!\s*[=<>])');
    final matches = forceUnwrapRegex.allMatches(content);

    for (final match in matches) {
      _addIssue(
        file: file,
        line: _findLineNumber(content, match.group(0)!),
        type: 'null_safety',
        severity: 'medium',
        message: 'Force unwrapping detected: ${match.group(0)}',
        suggestion: 'Consider using null-aware operators or proper null checks',
      );
    }

    // Missing null checks for context operations
    if (content.contains('.of(context)') &&
        !content.contains('context.mounted')) {
      _addIssue(
        file: file,
        line: _findLineNumber(content, '.of(context)'),
        type: 'null_safety',
        severity: 'medium',
        message: 'Context operation without mounted check',
        suggestion: 'Check context.mounted before using context operations',
      );
    }
  }

  /// Check async patterns
  void _checkAsyncPatterns(String file, String content, List<String> lines) {
    // Async without await
    if (content.contains('async') && !content.contains('await')) {
      _addIssue(
        file: file,
        line: _findLineNumber(content, 'async'),
        type: 'async_patterns',
        severity: 'low',
        message: 'Async function without await',
        suggestion: 'Remove async keyword if not using await',
      );
    }

    // Missing error handling for async operations
    if (content.contains('await') &&
        !content.contains('try') &&
        !content.contains('catchError')) {
      _addIssue(
        file: file,
        line: _findLineNumber(content, 'await'),
        type: 'async_patterns',
        severity: 'medium',
        message: 'Async operation without error handling',
        suggestion: 'Add try-catch or .catchError for async operations',
      );
    }
  }

  /// Helper method to add an issue
  void _addIssue({
    required String file,
    required int line,
    required String type,
    required String severity,
    required String message,
    required String suggestion,
  }) {
    _issues.add({
      'file': file,
      'line': line,
      'type': type,
      'severity': severity,
      'message': message,
      'suggestion': suggestion,
    });
  }

  /// Helper to find line number of a pattern
  int _findLineNumber(String content, String pattern) {
    final index = content.indexOf(RegExp(pattern));
    if (index == -1) return 1;
    return content.substring(0, index).split('\n').length;
  }

  /// Extract method content from a starting position
  String _extractMethodContent(String content, int start) {
    int braceCount = 0;
    int i = start;

    // Find the opening brace
    while (i < content.length && content[i] != '{') {
      i++;
    }

    if (i >= content.length) return '';

    final methodStart = i;
    braceCount = 1;
    i++;

    // Find the matching closing brace
    while (i < content.length && braceCount > 0) {
      if (content[i] == '{') {
        braceCount++;
      } else if (content[i] == '}') {
        braceCount--;
      }
      i++;
    }

    return content.substring(methodStart, i);
  }

  /// Print analysis report
  void printReport(Map<String, dynamic> report) {
    final summary = report['summary'] as Map<String, dynamic>;
    final issues = (report['issues'] as List).cast<Map<String, dynamic>>();

    print('\n${'=' * 60}');
    print('🎯 FLUTTER CODE REVIEW REPORT');
    print('=' * 60);
    print('📁 Total files: ${summary['total_files']}');
    print('🎯 Dart files analyzed: ${summary['dart_files']}');
    print('🚨 Total issues: ${summary['issues']}');
    print('🔴 High severity: ${summary['high_severity']}');
    print('🟡 Medium severity: ${summary['medium_severity']}');
    print('🟢 Low severity: ${summary['low_severity']}');

    if (issues.isNotEmpty) {
      print('\n📋 ISSUES FOUND:');
      print('-' * 40);

      final groupedIssues = <String, List<Map<String, dynamic>>>{};
      for (final issue in issues) {
        final type = issue['type'] as String;
        groupedIssues.putIfAbsent(type, () => []).add(issue);
      }

      for (final entry in groupedIssues.entries) {
        print(
          '\n${_getTypeIcon(entry.key)} ${entry.key.toUpperCase().replaceAll('_', ' ')}:',
        );
        for (final issue in entry.value) {
          final severity = issue['severity'] as String;
          final icon = severity == 'high'
              ? '🔴'
              : severity == 'medium'
              ? '🟡'
              : '🟢';
          print(
            '  $icon ${issue['file']}:${issue['line']} - ${issue['message']}',
          );
          print('    💡 ${issue['suggestion']}');
        }
      }
    } else {
      print('\n✅ No issues found! Great job!');
    }
  }

  String _getTypeIcon(String type) {
    switch (type) {
      case 'widget_structure':
        return '🏗️';
      case 'state_management':
        return '🔄';
      case 'performance':
        return '⚡';
      case 'code_style':
        return '🎨';
      case 'null_safety':
        return '🛡️';
      case 'async_patterns':
        return '⏱️';
      default:
        return '📝';
    }
  }
}

void main(List<String> args) async {
  final repoPath = args.isNotEmpty ? args[0] : '.';

  final analyzer = FlutterReviewAnalyzer(repoPath);
  final report = await analyzer.analyze();

  analyzer.printReport(report);

  // Output JSON if requested
  if (args.contains('--json')) {
    print('\n${jsonEncode(report)}');
  }
}
