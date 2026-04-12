# Sentinel Security Journal

## 2026-04-12 - Initial Security Hardening
**Vulnerability:** Several security gaps were identified:
1. HTTP server using default configuration without timeouts, making it susceptible to Slowloris and other DoS attacks.
2. Global `http.ServeMux` being used, which is a shared resource and can lead to accidental route exposure or shadowing.
3. `/converter` endpoint not limiting request body size, allowing potential memory exhaustion DoS.
4. Missing standard security headers (`Content-Type`, `X-Content-Type-Options`).
5. Lack of validation for version strings from external files.

**Learning:** Lightweight APIs often overlook basic hardening because of their simplicity. Even simple conversion services should implement defense-in-depth.

**Prevention:** Use custom `http.Server` with timeouts, local `ServeMux`, and implement input validation/resource limiting at the handler level.
