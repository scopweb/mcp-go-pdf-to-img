## Description
<!-- Describe your changes in detail -->

## Type of Change
<!-- Mark the relevant options with an "x" -->

- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update
- [ ] Security patch
- [ ] Dependency update
- [ ] Refactoring (no functional changes)

## Related Issue
<!-- Link the issue this PR closes -->

Closes #(issue number)

## Motivation and Context
<!-- Why is this change needed? What problem does it solve? -->

## How Has This Been Tested?
<!-- Describe the tests you ran and how to reproduce them -->

- [ ] Manual testing completed
- [ ] All existing tests pass
- [ ] Security tests pass
- [ ] Added new tests for changes

## Testing Instructions
<!-- Step-by-step instructions to verify the changes -->

```bash
# Example test commands
go build -o pdf2img ./cmd/pdf2img
go test ./test/security -v
./pdf2img -i example.pdf -o output
```

## Checklist
<!-- Ensure all items are complete before submitting -->

- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Security tests pass (19/19)
- [ ] Any dependent changes have been merged and published
- [ ] I have not introduced any hardcoded credentials or secrets
- [ ] I have no breaking changes to the public API

## Screenshots / Logs
<!-- If applicable, add screenshots or logs to demonstrate the change -->

## Performance Impact
<!-- Any performance implications of this change? -->

## Breaking Changes
<!-- List any breaking changes, if applicable -->

- None

## Additional Notes
<!-- Any additional context or notes -->

---

**Thank you for your contribution!** 🎉

Before marking as ready for review:
1. Ensure all tests pass locally
2. Review your own code first
3. Check for hardcoded secrets
4. Verify documentation is updated
5. Add appropriate labels
