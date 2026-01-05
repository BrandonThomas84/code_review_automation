#!/usr/bin/env ruby

# Ruby-specific code review automation
# Analyzes Ruby code for common issues and best practices

require 'json'
require 'pathname'

class RubyReviewAnalyzer
  attr_reader :issues, :repo_path, :target_branch

  def initialize(repo_path = '.', target_branch = 'main')
    @repo_path = Pathname.new(repo_path)
    @target_branch = target_branch
    @issues = []
  end

  def analyze
    puts "🔍 Starting Ruby code analysis..."
    
    changed_files = get_changed_files
    ruby_files = changed_files.select { |f| f.end_with?('.rb', '.rake', 'Rakefile', 'Gemfile') }
    
    if ruby_files.empty?
      return {
        'summary' => { 'total_files' => 0, 'ruby_files' => 0, 'issues' => 0 },
        'issues' => []
      }
    end

    puts "📁 Found #{ruby_files.length} Ruby files to analyze"

    ruby_files.each { |file| analyze_file(file) }

    {
      'summary' => {
        'total_files' => changed_files.length,
        'ruby_files' => ruby_files.length,
        'issues' => @issues.length,
        'high_severity' => @issues.count { |i| i['severity'] == 'high' },
        'medium_severity' => @issues.count { |i| i['severity'] == 'medium' },
        'low_severity' => @issues.count { |i| i['severity'] == 'low' }
      },
      'issues' => @issues,
      'files_analyzed' => ruby_files
    }
  end

  def print_report(report)
    summary = report['summary']
    issues = report['issues']

    puts "\n#{'=' * 60}"
    puts "💎 RUBY CODE REVIEW REPORT"
    puts '=' * 60
    puts "📁 Total files: #{summary['total_files']}"
    puts "💎 Ruby files analyzed: #{summary['ruby_files']}"
    puts "🚨 Total issues: #{summary['issues']}"
    puts "🔴 High severity: #{summary['high_severity']}"
    puts "🟡 Medium severity: #{summary['medium_severity']}"
    puts "🟢 Low severity: #{summary['low_severity']}"

    if issues.any?
      puts "\n📋 ISSUES FOUND:"
      puts '-' * 40

      grouped_issues = issues.group_by { |issue| issue['type'] }

      grouped_issues.each do |type, type_issues|
        puts "\n#{get_type_icon(type)} #{type.upcase.gsub('_', ' ')}:"
        type_issues.each do |issue|
          severity_icon = case issue['severity']
                         when 'high' then '🔴'
                         when 'medium' then '🟡'
                         else '🟢'
                         end
          puts "  #{severity_icon} #{issue['file']}:#{issue['line']} - #{issue['message']}"
          puts "    💡 #{issue['suggestion']}"
        end
      end
    else
      puts "\n✅ No issues found! Great job!"
    end
  end

  private

  def get_changed_files
    result = `git diff --name-only #{@target_branch}..HEAD 2>/dev/null`
    if $?.success?
      result.strip.split("\n").reject(&:empty?)
    else
      puts "⚠️  Warning: Could not get git diff, analyzing all Ruby files"
      get_all_ruby_files
    end
  rescue
    puts "⚠️  Warning: Git not available, analyzing all Ruby files"
    get_all_ruby_files
  end

  def get_all_ruby_files
    files = []
    Dir.glob("#{@repo_path}/**/*.rb").each do |file|
      files << file.sub("#{@repo_path}/", '')
    end
    files
  end

  def analyze_file(file_path)
    full_path = @repo_path.join(file_path)
    return unless full_path.exist?

    content = full_path.read
    lines = content.split("\n")

    # Run various checks
    check_rails_patterns(file_path, content, lines)
    check_security_issues(file_path, content, lines)
    check_performance_issues(file_path, content, lines)
    check_code_style(file_path, content, lines)
    check_error_handling(file_path, content, lines)
    check_testing_patterns(file_path, content, lines)
    check_database_patterns(file_path, content, lines)
  end

  def check_rails_patterns(file, content, lines)
    # N+1 query patterns
    if content.match?(/\.each\s+do.*\.find/)
      add_issue(
        file: file,
        line: find_line_number(content, /\.each\s+do.*\.find/),
        type: 'rails_performance',
        severity: 'high',
        message: 'Potential N+1 query detected',
        suggestion: 'Use includes, joins, or preload to avoid N+1 queries'
      )
    end

    # Missing strong parameters
    if content.include?('params[') && !content.include?('permit')
      add_issue(
        file: file,
        line: find_line_number(content, /params\[/),
        type: 'rails_security',
        severity: 'high',
        message: 'Direct params access without strong parameters',
        suggestion: 'Use strong parameters with permit() method'
      )
    end

    # Fat controllers
    if file.include?('controller') && content.split("\n").length > 100
      add_issue(
        file: file,
        line: 1,
        type: 'rails_structure',
        severity: 'medium',
        message: "Large controller file (#{content.split("\n").length} lines)",
        suggestion: 'Consider extracting logic to services or concerns'
      )
    end

    # Missing validations in models
    if file.include?('model') && content.include?('class') && !content.include?('validates')
      add_issue(
        file: file,
        line: 1,
        type: 'rails_structure',
        severity: 'medium',
        message: 'Model without validations',
        suggestion: 'Add appropriate validations to ensure data integrity'
      )
    end

    # Callback hell
    callback_count = content.scan(/before_|after_|around_/).length
    if callback_count > 5
      add_issue(
        file: file,
        line: 1,
        type: 'rails_structure',
        severity: 'medium',
        message: "Too many callbacks (#{callback_count})",
        suggestion: 'Consider using service objects or observers'
      )
    end
  end

  def check_security_issues(file, content, lines)
    # SQL injection risks
    if content.match?(/where\s*\(\s*["'].*#\{/)
      add_issue(
        file: file,
        line: find_line_number(content, /where\s*\(\s*["'].*#\{/),
        type: 'security',
        severity: 'high',
        message: 'Potential SQL injection with string interpolation',
        suggestion: 'Use parameterized queries or ActiveRecord methods'
      )
    end

    # Mass assignment vulnerabilities
    if content.include?('.update_attributes(params[') || content.include?('.create(params[')
      add_issue(
        file: file,
        line: find_line_number(content, /params\[/),
        type: 'security',
        severity: 'high',
        message: 'Mass assignment vulnerability',
        suggestion: 'Use strong parameters to whitelist allowed attributes'
      )
    end

    # Hardcoded secrets
    secret_patterns = [
      /password\s*=\s*["'][^"']+["']/i,
      /api_key\s*=\s*["'][^"']+["']/i,
      /secret\s*=\s*["'][^"']+["']/i,
      /token\s*=\s*["'][^"']+["']/i
    ]

    secret_patterns.each do |pattern|
      if content.match?(pattern)
        add_issue(
          file: file,
          line: find_line_number(content, pattern),
          type: 'security',
          severity: 'high',
          message: 'Hardcoded secret detected',
          suggestion: 'Use environment variables or Rails credentials'
        )
      end
    end

    # Unsafe eval
    if content.include?('eval(') || content.include?('instance_eval(')
      add_issue(
        file: file,
        line: find_line_number(content, /eval\(/),
        type: 'security',
        severity: 'high',
        message: 'Use of eval() - security risk',
        suggestion: 'Avoid eval() or use safer alternatives like send()'
      )
    end
  end

  def check_performance_issues(file, content, lines)
    # Inefficient queries in loops
    if content.match?(/\.each.*do.*\.(find|where|create|update)/)
      add_issue(
        file: file,
        line: find_line_number(content, /\.each.*do.*\.(find|where|create|update)/),
        type: 'performance',
        severity: 'medium',
        message: 'Database query inside loop',
        suggestion: 'Use bulk operations or move queries outside loops'
      )
    end

    # Missing database indexes (migration files)
    if file.include?('migrate') && content.include?('add_column') && !content.include?('add_index')
      add_issue(
        file: file,
        line: find_line_number(content, 'add_column'),
        type: 'performance',
        severity: 'low',
        message: 'New column without index consideration',
        suggestion: 'Consider adding indexes for frequently queried columns'
      )
    end

    # Inefficient string concatenation
    if content.match?(/\+\s*=.*["']/)
      add_issue(
        file: file,
        line: find_line_number(content, /\+\s*=.*["']/),
        type: 'performance',
        severity: 'low',
        message: 'String concatenation with +=',
        suggestion: 'Use Array#join or String interpolation for better performance'
      )
    end
  end

  def check_code_style(file, content, lines)
    lines.each_with_index do |line, index|
      line_num = index + 1

      # Long lines
      if line.length > 120
        add_issue(
          file: file,
          line: line_num,
          type: 'style',
          severity: 'low',
          message: "Line too long (#{line.length} characters)",
          suggestion: 'Break long lines for better readability'
        )
      end

      # TODO/FIXME comments
      if line.match?(/TODO|FIXME|HACK/i)
        add_issue(
          file: file,
          line: line_num,
          type: 'style',
          severity: 'low',
          message: 'TODO/FIXME comment found',
          suggestion: 'Address the TODO or create a proper issue'
        )
      end

      # Trailing whitespace
      if line.match?(/\s+$/)
        add_issue(
          file: file,
          line: line_num,
          type: 'style',
          severity: 'low',
          message: 'Trailing whitespace',
          suggestion: 'Remove trailing whitespace'
        )
      end
    end

    # Class/method naming conventions
    if content.match?(/class\s+[a-z]/)
      add_issue(
        file: file,
        line: find_line_number(content, /class\s+[a-z]/),
        type: 'style',
        severity: 'low',
        message: 'Class name should be CamelCase',
        suggestion: 'Use CamelCase for class names'
      )
    end

    # Method naming conventions
    if content.match?(/def\s+[A-Z]/)
      add_issue(
        file: file,
        line: find_line_number(content, /def\s+[A-Z]/),
        type: 'style',
        severity: 'low',
        message: 'Method name should be snake_case',
        suggestion: 'Use snake_case for method names'
      )
    end
  end

  def check_error_handling(file, content, lines)
    # Rescue without specific exception
    if content.match?(/rescue\s*$/) || content.match?(/rescue\s*=>/)
      add_issue(
        file: file,
        line: find_line_number(content, /rescue/),
        type: 'error_handling',
        severity: 'medium',
        message: 'Generic rescue clause',
        suggestion: 'Rescue specific exceptions instead of using generic rescue'
      )
    end

    # Empty rescue blocks
    if content.match?(/rescue.*\n\s*end/)
      add_issue(
        file: file,
        line: find_line_number(content, /rescue/),
        type: 'error_handling',
        severity: 'medium',
        message: 'Empty rescue block',
        suggestion: 'Handle errors appropriately or log them'
      )
    end

    # Missing error handling for external calls
    if content.include?('Net::HTTP') && !content.include?('rescue')
      add_issue(
        file: file,
        line: find_line_number(content, 'Net::HTTP'),
        type: 'error_handling',
        severity: 'medium',
        message: 'External HTTP call without error handling',
        suggestion: 'Add proper error handling for network calls'
      )
    end
  end

  def check_testing_patterns(file, content, lines)
    return unless file.include?('spec') || file.include?('test')

    # Missing test descriptions
    if content.include?('it ') && content.match?(/it\s+["'][^"']*["']\s+do/)
      # This is actually good - has description
    elsif content.include?('it do')
      add_issue(
        file: file,
        line: find_line_number(content, 'it do'),
        type: 'testing',
        severity: 'low',
        message: 'Test without description',
        suggestion: 'Add descriptive test names'
      )
    end

    # Large test files
    if lines.length > 200
      add_issue(
        file: file,
        line: 1,
        type: 'testing',
        severity: 'medium',
        message: "Large test file (#{lines.length} lines)",
        suggestion: 'Consider splitting into multiple test files'
      )
    end

    # Missing assertions
    if content.include?('it ') && !content.match?(/expect|assert|should/)
      add_issue(
        file: file,
        line: find_line_number(content, 'it '),
        type: 'testing',
        severity: 'medium',
        message: 'Test without assertions',
        suggestion: 'Add proper assertions to verify behavior'
      )
    end
  end

  def check_database_patterns(file, content, lines)
    # Missing foreign key constraints
    if file.include?('migrate') && content.include?('references') && !content.include?('foreign_key')
      add_issue(
        file: file,
        line: find_line_number(content, 'references'),
        type: 'database',
        severity: 'medium',
        message: 'Reference without foreign key constraint',
        suggestion: 'Add foreign_key: true for referential integrity'
      )
    end

    # Missing null constraints
    if file.include?('migrate') && content.include?('add_column') && !content.include?('null:')
      add_issue(
        file: file,
        line: find_line_number(content, 'add_column'),
        type: 'database',
        severity: 'low',
        message: 'Column without null constraint specification',
        suggestion: 'Explicitly specify null: true or null: false'
      )
    end
  end

  def add_issue(file:, line:, type:, severity:, message:, suggestion:)
    @issues << {
      'file' => file,
      'line' => line,
      'type' => type,
      'severity' => severity,
      'message' => message,
      'suggestion' => suggestion
    }
  end

  def find_line_number(content, pattern)
    lines = content.split("\n")
    lines.each_with_index do |line, index|
      return index + 1 if line.match?(pattern)
    end
    1
  end

  def get_type_icon(type)
    case type
    when 'rails_performance', 'performance' then '⚡'
    when 'rails_security', 'security' then '🔒'
    when 'rails_structure' then '🏗️'
    when 'style' then '🎨'
    when 'error_handling' then '🛡️'
    when 'testing' then '🧪'
    when 'database' then '🗄️'
    else '📝'
    end
  end
end

# Main execution
if __FILE__ == $0
  repo_path = ARGV[0] || '.'
  target_branch = ARGV[1] || 'main'
  
  analyzer = RubyReviewAnalyzer.new(repo_path, target_branch)
  report = analyzer.analyze
  
  analyzer.print_report(report)
  
  # Output JSON if requested
  if ARGV.include?('--json')
    puts "\n#{JSON.pretty_generate(report)}"
  end
end
