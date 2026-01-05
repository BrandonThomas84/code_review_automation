# Integration Guide for Other Repositories

This guide explains how to integrate the code-review tool into your other GitHub repositories.

## Quick Start (5 minutes)

### Step 1: Add GitHub Actions Workflow

Create `.github/workflows/code-review.yml` in your repository:

```yaml
name: Code Review

on:
  pull_request:
    branches:
      - main
      - develop

jobs:
  code-review:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'

      - name: Download code-review
        run: |
          RELEASE_URL=$(curl -s https://api.github.com/repos/BrandonThomas84/code-review-automation/releases/latest | grep browser_download_url | grep linux-amd64 | cut -d '"' -f 4)
          wget -q $RELEASE_URL -O code-review
          chmod +x code-review

      - name: Run code review
        run: ./code-review -t ${{ github.base_ref }} --json > report.json

      - name: Comment on PR
        if: always()
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const report = JSON.parse(fs.readFileSync('report.json', 'utf8'));
            const comment = `## Code Review Results\n\n**Issues:** ${report.summary.total_issues} (🔴 ${report.summary.high_severity} | 🟡 ${report.summary.medium_severity} | 🟢 ${report.summary.low_severity})`;
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });
```

### Step 2: Commit and Push

```bash
git add .github/workflows/code-review.yml
git commit -m "Add automated code review"
git push
```

### Step 3: Create a Pull Request

The workflow will automatically run on your next PR!

## Advanced Configuration

### Email Notifications

1. Add secrets to your repository:
   - `SMTP_HOST` - Your SMTP server (e.g., smtp.gmail.com)
   - `SMTP_USER` - Your email address
   - `SMTP_PASSWORD` - Your app password
   - `FROM_EMAIL` - Sender email address

2. Update workflow:
```yaml
- name: Run code review with email
  env:
    SMTP_HOST: ${{ secrets.SMTP_HOST }}
    SMTP_USER: ${{ secrets.SMTP_USER }}
    SMTP_PASSWORD: ${{ secrets.SMTP_PASSWORD }}
    FROM_EMAIL: ${{ secrets.FROM_EMAIL }}
  run: |
    ./code-review -t ${{ github.base_ref }} \
      --email your-email@example.com \
      --json > report.json
```

### Custom Branches

Modify the `on.pull_request.branches` section:

```yaml
on:
  pull_request:
    branches:
      - main
      - develop
      - staging
      - release/*
```

### Fail on High Severity Issues

Add this step to fail the workflow if high severity issues are found:

```yaml
- name: Check for high severity issues
  run: |
    HIGH=$(jq '.summary.high_severity' report.json)
    if [ "$HIGH" -gt 0 ]; then
      echo "❌ High severity issues found!"
      exit 1
    fi
```

### Full Codebase Scan

For comprehensive analysis of the entire codebase:

```yaml
- name: Run full codebase scan
  run: ./code-review -t ${{ github.base_ref }} --full-scan --json > report.json
```

## Local Testing

Before committing the workflow, test it locally:

```bash
# Download the binary
wget https://github.com/BrandonThomas84/code-review-automation/releases/download/v1.0.0/code-review-linux-amd64
chmod +x code-review-linux-amd64

# Run against your target branch
./code-review-linux-amd64 -t main

# Or with JSON output
./code-review-linux-amd64 -t main --json | jq .
```

## Troubleshooting

### Workflow not triggering
- Ensure the workflow file is in `.github/workflows/`
- Check that branch names match your repository
- Verify the workflow syntax: `git push` and check Actions tab

### Binary download fails
- Check the latest release: https://github.com/BrandonThomas84/code-review-automation/releases
- Verify your platform is supported (Linux, macOS, Windows)
- Check network connectivity in GitHub Actions

### Git diff shows no changes
- Ensure `fetch-depth: 0` is set in checkout step
- Verify the target branch exists in the repository
- Check that you're comparing against the correct branch

## Support

For issues or questions:
1. Check the main repository: https://github.com/BrandonThomas84/code-review-automation
2. Review the BUILD_AND_DEPLOY.md guide
3. Check GitHub Actions logs for detailed error messages

## Next Steps

1. ✅ Add the workflow to your repository
2. ✅ Create a test PR to verify it works
3. ✅ Configure email notifications (optional)
4. ✅ Customize severity thresholds (optional)
5. ✅ Share with your team!

