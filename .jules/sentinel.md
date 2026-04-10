## 2026-04-10 - [Security Hardening]
**Vulnerability:** Multiple defense-in-depth gaps including lack of server timeouts, unrestricted HTTP methods, no request body size limits, and unvalidated version input.
**Learning:** Default Go HTTP server and mux configurations are permissive and vulnerable to DoS (Slowloris, large payloads) and potentially other attacks (MIME-sniffing).
**Prevention:** Always use a local 'http.NewServeMux', configure 'http.Server' with explicit timeouts, use 'http.MaxBytesReader' to limit payload size, and strictly validate all file-based or environment-based inputs.
