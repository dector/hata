# Hata Server Roadmap Ideas

Prioritized with simpler and higher-impact work first.

- Priority: `P0` (now), `P1` (next), `P2` (later)
- Effort: `S` (small), `M` (medium), `L` (large)

1. **[P0][S] Add auth middleware for protected routes** (fix/architecture)
   - Move repeated bearer-token checks from handlers into shared chi middleware and inject authenticated user context.

2. **[P0][S] Add logout endpoint** (feature/security)
   - Implement `POST /api/latest/auth/logout` to invalidate the current session token (`invalid_since`) cleanly.

3. **[P0][S] Strengthen request validation** (fix)
   - Trim/validate required fields and reject malformed payloads consistently across login and management APIs.

4. **[P0][S] API versioning consistency pass** (fix)
   - Align route comments and endpoints (`latest` vs `next`) and keep version constants in one place.

5. **[P0][S] Add context timeouts for DB-backed handlers** (quality)
   - Wrap handler DB calls with bounded context deadlines to prevent requests hanging indefinitely.

6. **[P0][M] Replace fmt prints with structured logging** (quality/ops)
   - Use structured logs with request metadata and avoid leaking sensitive values in auth failure paths.

7. **[P0][M] Session lifecycle hardening** (security)
   - Add rotation and bounded concurrent sessions per user with explicit revocation controls.

8. **[P0][M] Uniform error contract and mapping** (architecture)
   - Centralize domain-to-HTTP error mapping so internal DB errors never leak and client codes stay stable.

9. **[P0][M] Expand API tests for houses + auth helpers** (quality)
   - Add coverage for house listing edge cases, expired/invalid session states, and malformed authorization headers.

10. **[P1][M] Pagination support for list endpoints** (feature)
    - Add `limit`/`offset` parameters and deterministic ordering for `/house`, `/device`, and house-device listings.

11. **[P1][M] Device state update endpoint** (feature)
    - Implement authenticated `PATCH /api/latest/device/{id}` to toggle/update device state with authorization checks.

12. **[P1][M] Role-based access rules cleanup** (fix/security)
    - Enforce per-role capabilities (owner/admin/habitant/guest) for read/write operations instead of read-only checks.

13. **[P1][M] Add request ID and trace propagation** (ops)
    - Generate/request pass-through IDs in middleware and include them in logs and error responses.

14. **[P1][M] Basic rate limiting on auth routes** (security)
    - Protect login endpoint against brute force attempts with IP and username windows.

15. **[P1][M] Health/readiness split** (ops)
    - Keep lightweight liveness endpoint and add readiness endpoint validating DB connectivity and migration status.

16. **[P1][M] Config via env with validation** (quality)
    - Move hardcoded server/db settings into env-config structs with defaults and startup validation.

17. **[P1][L] OpenAPI spec + generated clients** (feature/quality)
    - Publish machine-readable API contract and use it to keep Android/server payloads synchronized.

18. **[P1][L] Migration workflow improvements** (ops)
    - Add explicit migration commands, rollback strategy docs, and CI checks for schema drift.

19. **[P1][L] Integration plugin boundary for device providers** (architecture)
    - Introduce adapter interfaces for providers (for example WiZ) so discovery/control logic is modular and testable.

20. **[P2][M] Audit trail for security-sensitive actions** (security/ops)
    - Record user/session actions (login, logout, role assignment, device updates) with immutable timestamps.

21. **[P2][M] Idempotency keys for write endpoints** (quality)
    - Support idempotent retries for state-changing calls to avoid duplicate effects during network retries.

22. **[P2][L] Background sync workers for device status** (feature)
    - Add periodic jobs to reconcile persisted state with provider state and mark stale devices.

23. **[P2][L] Real-time updates channel** (feature)
    - Add SSE/WebSocket stream for live device/session/house events consumed by mobile clients.

24. **[P2][L] Multi-tenant hardening** (security/architecture)
    - Add explicit tenant scoping strategy to prevent cross-house data leaks as the data model grows.

25. **[P2][L] Backup/restore tooling for SQLite** (ops)
    - Provide safe snapshot and restore commands with lock handling and integrity verification.
