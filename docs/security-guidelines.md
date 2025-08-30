# Security Guidelines for Claude AI Workflows

## 🔐 API Key Security

### Current Security Measures

#### ✅ **Secret Management**
- **Storage**: All API keys stored in GitHub Secrets (never hardcoded)
- **Access**: Secrets accessed via `${{ secrets.ANTHROPIC_API_KEY }}` syntax only
- **Masking**: GitHub automatically masks secrets in workflow logs
- **Validation**: Each workflow validates API key presence before execution
- **OIDC Authentication**: Claude Code Action uses `id-token: write` for secure GitHub authentication

#### ✅ **Logging Protection**
```yaml
# SECURE: API key in environment variable (masked in logs)
env:
  ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}

# SECURE: Validation without exposing key
if [ -z "$ANTHROPIC_API_KEY" ]; then
  echo "::error::API key not configured"
fi
```

#### ❌ **What to Avoid**
```yaml
# INSECURE: Never log or echo secrets
- run: echo "Key: ${{ secrets.ANTHROPIC_API_KEY }}"

# INSECURE: Never hardcode secrets
anthropic_api_key: "sk-ant-actual-key-here"

# INSECURE: Never use secrets in conditions that might log
if: ${{ secrets.ANTHROPIC_API_KEY != '' }}
```

## 🛡️ Workflow Security

### Permission Model (Principle of Least Privilege)

Each Claude workflow uses minimal permissions:

```yaml
permissions:
  contents: read          # Read repository files
  issues: write          # Create/update management reports  
  pull-requests: read    # Read PR information
  actions: read          # Read CI/CD results (where needed)
  id-token: write        # ✅ REQUIRED - Claude Code Action uses OIDC for GitHub authentication
```

### Security Validation Steps

All workflows include security validation:

```yaml
- name: Security Validation
  env:
    ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
  run: |
    # Validate without exposing
    if [ -z "$ANTHROPIC_API_KEY" ]; then
      echo "::error::API key not configured"
      exit 1
    fi
    
    if [[ ${#ANTHROPIC_API_KEY} -lt 10 ]]; then
      echo "::error::API key appears invalid"
      exit 1
    fi
    
    echo "✅ Security checks passed"
```

## 🔍 Monitoring & Auditing

### Automated Security Audits

- **Schedule**: Weekly security audits (Sundays 2 AM UTC)
- **Checks**: Hardcoded secrets, insecure logging, missing validations
- **Reports**: Generated in workflow artifacts

### Manual Security Checks

#### Monthly Review Checklist:
- [ ] Check GitHub Security tab for secret scanning alerts
- [ ] Review workflow logs for any anomalies  
- [ ] Verify all new workflows include security validation
- [ ] Confirm no hardcoded secrets in codebase
- [ ] Validate secret rotation schedule

#### Incident Response:
1. **If API key is compromised**:
   - Immediately rotate `ANTHROPIC_API_KEY` in GitHub Secrets
   - Review recent workflow executions
   - Check for unauthorized API usage
   - Update key in all environments

2. **If secret exposure is detected**:
   - Revoke compromised key immediately
   - Generate new API key from Anthropic Console
   - Update GitHub Secret
   - Review and enhance security measures

## 🚨 Security Alerts

### GitHub Secret Scanning
- Automatically scans for known secret patterns
- Alerts appear in Security tab → Secret scanning
- Immediate action required if alerts detected

### Monitoring Points
- Unusual API usage patterns
- Failed authentication attempts
- Workflow failures related to authentication
- Unauthorized access to sensitive endpoints

## 📋 Security Best Practices

### Development
1. **Never commit secrets** to the repository
2. **Use GitHub Secrets** for all sensitive data
3. **Validate inputs** in all workflows
4. **Follow least privilege** principle for permissions
5. **Review security regularly** through audits

### Deployment
1. **Rotate keys periodically** (quarterly recommended)
2. **Monitor usage patterns** for anomalies
3. **Keep workflows updated** with latest security practices
4. **Document security changes** in commit messages

### Team Guidelines
1. **Security awareness**: All team members understand secret handling
2. **Code review**: Security-focused reviews for workflow changes
3. **Incident reporting**: Clear escalation path for security issues
4. **Regular training**: Keep team updated on security best practices

## 🔧 Security Tools

### Repository Security Features
- **Secret scanning**: Enabled for automatic detection
- **Dependency scanning**: Monitors for vulnerable dependencies  
- **Code scanning**: Static analysis for security issues
- **Branch protection**: Prevents direct pushes to main branches

### Workflow Security
- **Environment protection**: Restrict workflow execution contexts
- **Required reviews**: Security-focused reviews for sensitive changes
- **Audit logs**: Complete audit trail of all workflow executions

## 📞 Security Contacts

- **Security Issues**: Report via GitHub Security tab
- **Team Lead**: Escalate critical security concerns
- **Anthropic Support**: For API key or service-related security issues

---

**Remember**: Security is everyone's responsibility. When in doubt, err on the side of caution and ask for a security review.