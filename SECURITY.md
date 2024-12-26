# Security Guidelines

## Reporting Issues

If you discover a security issue in dns-updater, please create a GitHub issue with:
- A description of the problem
- Steps to reproduce
- Potential impact

## Basic Security Practices

### Credentials
- Don't commit API tokens or credentials directly to Git
- Store credentials in a config.yaml file (kept outside of version control)

### Docker Security
The container runs as a non-root user by default, which is sufficient for most private deployments.

### Configuration
- Keep your config.yaml file with restricted read permissions
- Use API tokens with necessary permissions for your domains
- Consider rotating credentials periodically (every few months is fine)

## Updates

Keep the application and its dependencies updated when convenient, particularly if you notice security warnings from GitHub.

## Development

Basic security practices when contributing:
1. Don't commit sensitive data
2. Keep dependencies reasonably up to date
3. Use HTTPS for API communications

## Note

This is a private project, and these guidelines are intentionally simplified. If you deploy this in a more sensitive environment, consider implementing stricter security measures.
