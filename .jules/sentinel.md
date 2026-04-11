## 2026-04-11 - HTTP Server Hardening
**Vulnerability:** Lack of server timeouts and global state usage via `http.DefaultServeMux` made the application vulnerable to slowloris attacks and potentially unexpected routing behavior. Missing security headers allowed for MIME-sniffing.
**Learning:** Standard library defaults in Go's `net/http` package (like `ListenAndServe` with `nil` handler) use global state and lack timeouts, which are insufficient for production security.
**Prevention:** Always use a local `http.NewServeMux`, explicitly define `http.Server` with `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`, and use `http.MaxBytesReader` to limit request body sizes.
