# Changelog

## v0.1.31

- Removed the process-global, implicitly registered S3 SigV4 authenticator and
  the `github.com/xxxsen/s3verify` dependency.
- Added `auth.NewBasic()` and `webapi.WithAuthenticators(...)` for explicit,
  ordered, engine-local authentication.
- `middleware.TryAuthMiddleware` remains source-compatible when called with
  only a user resolver; that form now enables Basic authentication only.
- Consumers that referenced `S3V4AuthName` or depended on implicit SigV4 must
  provide their own explicit authenticator or verify S3 requests in their S3
  protocol layer before upgrading.
