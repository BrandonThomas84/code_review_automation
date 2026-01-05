# .autoreview-ignore Guide

The `.autoreview-ignore` file allows you to exclude files and directories from code review analysis.

## Creating the File

Create a `.autoreview-ignore` file in the root of your repository:

```bash
touch .autoreview-ignore
```

## File Format

Each line in `.autoreview-ignore` specifies a pattern to ignore:

```text
# Comments start with #
# Empty lines are ignored

# Exact file matches
config.json
secrets.env

# Wildcard patterns
*.min.js
*.test.py
build/*

# Directory patterns (end with /)
node_modules/
dist/
.venv/
vendor/

# Glob patterns
src/**/*.generated.go
tests/fixtures/*
```

## Pattern Types

### 1. Exact File Matches

```text
config.json
src/utils.py
```

Ignores files with exact paths.

### 2. Wildcard Patterns

```text
*.min.js
*.test.py
*.generated.go
```

Uses standard glob patterns with `*` and `?`.

### 3. Directory Patterns

```text
node_modules/
dist/
.venv/
```

End with `/` to ignore entire directories and their contents.

### 4. Glob Patterns

```text
build/*
tests/fixtures/*
src/**/*.generated.go
```

Supports `*` for any characters and `**` for nested directories.

## Examples

### Python Project

```text
# Virtual environments
.venv/
venv/
env/

# Generated files
*.pyc
__pycache__/
*.egg-info/

# Test coverage
.coverage
htmlcov/

# IDE
.vscode/
.idea/

# Dependencies
site-packages/
```

### JavaScript/Node Project

```text
# Dependencies
node_modules/
package-lock.json

# Build output
dist/
build/
.next/

# Generated files
*.min.js
*.min.css

# IDE
.vscode/
.idea/

# Environment
.env
.env.local
```

### Go Project

```text
# Build output
bin/
dist/

# Generated code
*.pb.go
*_gen.go

# Vendor
vendor/

# IDE
.vscode/
.idea/
```

### General

```text
# Compiled files
*.o
*.a
*.so
*.exe

# Archives
*.zip
*.tar.gz

# OS files
.DS_Store
Thumbs.db

# Temporary files
*.tmp
*.swp
*~
```

## Comments

Lines starting with `#` are treated as comments:

```text
# This is a comment
# Ignore all test files
*.test.js

# Ignore build directories
dist/
build/
```

## Whitespace

- Leading and trailing whitespace is automatically trimmed
- Empty lines are ignored
- Indentation is not significant

## How It Works

1. When the analyzer starts, it reads `.autoreview-ignore` from the repository root
2. Each non-comment, non-empty line is treated as a pattern
3. Files matching any pattern are excluded from analysis
4. Both `git diff` analysis and full codebase scans respect these patterns

## Testing Your Patterns

To verify your patterns work:

```bash
# Run with verbose output to see which files are being analyzed
go run ./cmd/code-review -t main -v

# Or with the built binary
./bin/code-review -t main -v
```

The verbose output will show which files are being analyzed.

## Best Practices

1. **Keep it organized** - Group related patterns with comments
2. **Be specific** - Use exact paths when possible to avoid unintended exclusions
3. **Document** - Add comments explaining why files are ignored
4. **Review regularly** - Update patterns as your project evolves
5. **Version control** - Commit `.autoreview-ignore` to your repository

## Common Mistakes

❌ **Wrong**: `node_modules` (without trailing slash)
✅ **Right**: `node_modules/`

❌ **Wrong**: `*.js` (too broad, ignores all JS)
✅ **Right**: `*.min.js` (only minified JS)

❌ **Wrong**: `/src/` (leading slash not needed)
✅ **Right**: `src/`

## Troubleshooting

### Files are still being analyzed

- Check that `.autoreview-ignore` is in the repository root
- Verify the pattern syntax is correct
- Run with `-v` flag to see which files are being analyzed

### Too many files are being ignored

- Review your patterns for overly broad matches
- Use more specific patterns
- Add comments to document why each pattern exists

### Pattern not working

- Ensure no leading/trailing spaces in patterns
- Check that directory patterns end with `/`
- Verify the file path matches the pattern exactly

## Examples in Action

### Example 1: Ignore generated files

```text
# .autoreview-ignore
*.generated.go
*.pb.go
*_gen.py
```

### Example 2: Ignore vendor and build directories

```text
# .autoreview-ignore
vendor/
node_modules/
dist/
build/
```

### Example 3: Ignore test files and fixtures

```text
# .autoreview-ignore
*.test.js
*.test.py
tests/fixtures/
test/mocks/
```

## Integration with CI/CD

The `.autoreview-ignore` file is automatically respected when running in GitHub Actions:

```yaml
- name: Run code review
  run: ./code-review -t main --json > report.json
```

No additional configuration needed!
