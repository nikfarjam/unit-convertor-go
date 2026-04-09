## 2026-04-09 - API Hardening and Denial of Service Prevention
**Vulnerability:** The application lacked several basic security protections:
- **Missing Request Body Size Limit:** The /converter endpoint was vulnerable to memory exhaustion attacks because it did not limit the size of the incoming request body.
- **Missing Server Timeouts:** The HTTP server used default configurations without Read, Write, or Idle timeouts, making it susceptible to Slowloris-style DoS attacks.
- **Missing Security Headers:** Responses were missing standard security headers like `X-Content-Type-Options: nosniff`, and did not explicitly set `Content-Type: application/json` for JSON responses.
- **Open HTTP Methods:** The `/converter` endpoint was open to all HTTP methods, while it only intended to support `POST`.

**Learning:** While simple Go web servers are easy to start, they are not secure by default. Standard library defaults for `http.ListenAndServe` do not include timeouts, and `json.NewDecoder(r.Body).Decode()` will attempt to read any size body unless limited.

**Prevention:**
1. Always use `http.MaxBytesReader` to limit request body size.
2. Always configure a custom `http.Server` with reasonable timeouts.
3. Explicitly set `Content-Type` and `X-Content-Type-Options` headers in all handlers.
4. Use Go 1.22+ method-based routing to restrict endpoints to necessary HTTP methods.
