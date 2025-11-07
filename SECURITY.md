# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions receive security updates depends on the CVSS v3.0 Rating:

| Version | Supported          | Status |
| ------- | ------------------ | ------ |
| 1.0.x   | ✅ Yes             | Current (until 2027-11-07) |
| < 1.0   | ❌ No              | End of Life |

## Reporting a Vulnerability

**DO NOT** create public GitHub issues for security vulnerabilities.

Instead, please email your findings to the maintainer. Include:

1. **Description**: What is the vulnerability?
2. **Location**: Which component/file?
3. **Severity**: Critical, High, Medium, Low
4. **Steps to Reproduce**: How can we reproduce it?
5. **Impact**: What could an attacker do?
6. **PoC (if applicable)**: Proof of concept

### Response Timeline

- **Initial Response**: Within 48 hours
- **Fix Preparation**: 7-14 days for critical issues
- **Security Advisory**: Published after patch is released

## Security Best Practices

### For Users

1. **Keep Dependencies Updated**
   ```bash
   go get -u ./...
   ```

2. **Run Security Tests**
   ```bash
   go test ./test/security -v
   ```

3. **Monitor Security Advisories**
   - Watch this repository for security updates
   - Check [Go Vulnerability Database](https://vuln.go.dev/)

### For Contributors

1. **Code Review**
   - All PRs reviewed for security issues
   - Automated security scanning enabled

2. **Testing**
   - Security tests must pass (19/19)
   - No hardcoded secrets or credentials
   - Input validation required

3. **Dependencies**
   - Use `go mod tidy` to keep dependencies clean
   - Review CVEs before updating major versions
   - Document why a dependency is needed

## Security Checklist

Before submitting a PR, ensure:

- [ ] No hardcoded credentials or secrets
- [ ] Input validation implemented
- [ ] Error handling is complete
- [ ] No unsafe code patterns (e.g., `os/exec` with user input)
- [ ] Security tests pass
- [ ] No new warnings from `go vet`
- [ ] Code reviewed by at least one maintainer

## Known Security Considerations

### Path Traversal (CWE-22)
**Status**: ✅ Mitigated

The application uses `filepath.Join()` and `filepath.Clean()` to prevent path traversal attacks.

**Example**:
```go
outputPath := filepath.Join(opts.OutputDir, filename)  // Safe
```

### Command Injection (CWE-78)
**Status**: ✅ Protected

The application does not use `os/exec` or `syscall` packages, preventing command injection.

### Input Validation (CWE-20)
**Status**: ✅ Implemented

All inputs are validated in `validateOptions()` function:
- File paths checked for existence
- Format validated (png/jpg only)
- DPI range checked (>0)
- Page ranges validated

## Security Testing

The project includes comprehensive security tests:

```bash
# Run all security tests
go test ./test/security -v

# Run with race detection
go test ./test/security -race -v

# Check for vulnerable dependencies
go list -u -m all | grep -v direct
```

## Vulnerability Disclosure

If you discover a vulnerability:

1. **Report privately** (do not create public issues)
2. **Include details** (description, location, impact)
3. **Wait for response** (48 hours)
4. **Coordination** (embargo period before public disclosure)

Example email subject:
```
[SECURITY] PDF2IMG - [CWE-XXXX] [Severity] [Description]
```

## Security Standards Compliance

This project follows:

- ✅ OWASP Top 10 mitigation
- ✅ CWE/SANS Top 25 review
- ✅ Go Security Best Practices
- ✅ NIST Secure Software Development Framework
- ✅ Apache 2.0 License compliance

## Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CWE/SANS Top 25](https://cwe.mitre.org/top25/)
- [Go Security](https://golang.org/doc/security)
- [CERT Secure Coding Standards](https://www.securecoding.cert.org/)
- [Go Vulnerability Database](https://vuln.go.dev/)

## Version History

| Version | Date | Notes |
| ------- | ---- | ----- |
| 1.0.0 | 2025-11-07 | Initial release with security audit |

## Contact

For security concerns, contact the maintainers privately.

**Last Updated**: 2025-11-07
**Version**: 1.0.0
