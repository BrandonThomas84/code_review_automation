# Deployment Checklist

Use this checklist to deploy the code review automation to your repositories.

## Phase 1: Build & Test (Local)

- [ ] Clone/navigate to code-review-automation repository
- [ ] Verify Go 1.25+ is installed: `go version`
- [ ] Build the binary: `make build`
- [ ] Verify binary works: `./bin/code-review --help`
- [ ] Test with JSON output: `./bin/code-review -t main --json`
- [ ] Build for all platforms: `make build-all`
- [ ] Verify all binaries exist: `ls -lh bin/`

## Phase 2: Create Release

- [ ] Update version in `internal/cmd/version.go` (if needed)
- [ ] Commit changes: `git add . && git commit -m "Release v1.0.0"`
- [ ] Create git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] Verify GitHub Actions runs
- [ ] Check GitHub Releases page for binaries
- [ ] Download and test a binary from release

## Phase 3: Prepare Target Repository

For each repository where you want to add code reviews:

- [ ] Clone/navigate to target repository
- [ ] Create workflows directory: `mkdir -p .github/workflows`
- [ ] Copy workflow template: `cp templates/github-actions-workflow.yml .github/workflows/code-review.yml`
- [ ] Review workflow file for customization needs
- [ ] Customize branch names if needed
- [ ] Commit workflow: `git add .github/workflows/code-review.yml`
- [ ] Commit message: `git commit -m "Add automated code review"`
- [ ] Push to repository: `git push`

## Phase 4: Test Workflow

For each target repository:

- [ ] Create a test branch: `git checkout -b test/code-review`
- [ ] Make a small code change
- [ ] Commit and push: `git push origin test/code-review`
- [ ] Create a pull request
- [ ] Wait for GitHub Actions to run
- [ ] Verify PR comment appears with results
- [ ] Check workflow logs for any errors
- [ ] Delete test branch: `git branch -D test/code-review`

## Phase 5: Email Configuration (Optional)

If you want email notifications:

- [ ] Obtain SMTP credentials:
  - [ ] SMTP_HOST (e.g., smtp.gmail.com)
  - [ ] SMTP_USER (your email)
  - [ ] SMTP_PASSWORD (app-specific password)
  - [ ] FROM_EMAIL (sender email)

- [ ] For each target repository:
  - [ ] Go to Settings → Secrets and variables → Actions
  - [ ] Add secret: `SMTP_HOST`
  - [ ] Add secret: `SMTP_USER`
  - [ ] Add secret: `SMTP_PASSWORD`
  - [ ] Add secret: `FROM_EMAIL`

- [ ] Update workflow to use secrets:
  - [ ] Add environment variables section
  - [ ] Add email flag to code-review command
  - [ ] Test with a new PR

## Phase 6: Customize Rules (Optional)

- [ ] Review `internal/review/analyzer.go`
- [ ] Add custom security patterns
- [ ] Add custom quality checks
- [ ] Rebuild: `make build`
- [ ] Test changes locally
- [ ] Create new release with updates

## Phase 7: Team Rollout

- [ ] Document the workflow for your team
- [ ] Share INTEGRATION_GUIDE.md with team
- [ ] Share QUICK_REFERENCE.md with team
- [ ] Conduct team training/demo
- [ ] Add to team documentation
- [ ] Monitor initial PRs for issues
- [ ] Gather feedback from team

## Phase 8: Monitoring & Maintenance

- [ ] Monitor GitHub Actions logs
- [ ] Track email delivery (if enabled)
- [ ] Collect team feedback
- [ ] Identify false positives
- [ ] Refine rules based on feedback
- [ ] Update documentation as needed
- [ ] Plan regular reviews of rules

## Troubleshooting Checklist

If something goes wrong:

- [ ] Check GitHub Actions logs
- [ ] Verify binary is executable: `chmod +x code-review`
- [ ] Verify git is installed: `git --version`
- [ ] Verify target branch exists: `git branch -a`
- [ ] Check SMTP credentials (if using email)
- [ ] Review workflow syntax
- [ ] Test binary locally first
- [ ] Check network connectivity

## Success Criteria

- [ ] Binary builds successfully
- [ ] Binary works locally
- [ ] GitHub Actions workflow runs
- [ ] PR comments appear with results
- [ ] Email sends (if configured)
- [ ] Team can use the tool
- [ ] No critical errors in logs

## Documentation Checklist

- [ ] Team has access to GETTING_STARTED.md
- [ ] Team has access to QUICK_REFERENCE.md
- [ ] Team has access to INTEGRATION_GUIDE.md
- [ ] Team knows how to read reports
- [ ] Team knows how to configure email
- [ ] Team knows how to customize rules

## Post-Deployment

- [ ] Monitor for 1 week
- [ ] Collect feedback from team
- [ ] Fix any issues found
- [ ] Document lessons learned
- [ ] Plan improvements
- [ ] Schedule regular reviews
- [ ] Update documentation

## Rollback Plan

If you need to disable code reviews:

- [ ] Delete `.github/workflows/code-review.yml` from repository
- [ ] Commit and push
- [ ] Workflow will no longer run on new PRs
- [ ] Existing PR comments will remain

## Support Resources

- **Quick Start**: GETTING_STARTED.md
- **Commands**: QUICK_REFERENCE.md
- **Building**: BUILD_AND_DEPLOY.md
- **Integration**: INTEGRATION_GUIDE.md
- **Technical**: GO_IMPLEMENTATION_SUMMARY.md
- **GitHub**: https://github.com/BrandonThomas84/code-review-automation

## Notes

Use this space to track your deployment:

```
Repository: ___________________
Date Started: ___________________
Date Completed: ___________________
Issues Encountered: ___________________
Team Feedback: ___________________
Next Steps: ___________________
```

---

**Status**: Ready for deployment ✅

