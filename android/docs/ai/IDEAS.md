# Hata Android Roadmap Ideas

Prioritized with simpler and higher-impact work first.

- Priority: `P0` (now), `P1` (next), `P2` (later)
- Effort: `S` (small), `M` (medium), `L` (large)

1. **[P0][S] Idempotent init/collector guards** (fix)
   - Ensure repeated `Init` actions cannot create duplicate Flow collectors in ViewModels.

2. **[P0][S] Personalization cleanup** (fix)
   - Remove hardcoded values like "Welcome home, Dan" and drive labels from session/home data.

3. **[P0][S] Manual refresh + last sync** (feature/fix)
   - Add pull-to-refresh and display last successful sync time in Home.

4. **[P0][S] Session timeout UX** (fix)
   - Show clear re-auth prompts and preserve unsaved UI intent when session expires.

5. **[P0][S] Accessibility pass** (quality)
   - Add semantic labels, larger tap targets, contrast checks, and TalkBack-friendly wording.

6. **[P0][M] Encrypted session/token storage** (fix/security)
   - Store auth/session data in encrypted storage and handle token expiration explicitly.

7. **[P0][M] Structured error model** (fix/architecture)
   - Replace generic `Exception` usage with typed errors for API/network/auth/device failures.

8. **[P0][M] Crash reporting and logging strategy** (quality/ops)
   - Integrate crash reporting and sanitize logs; avoid sensitive values in debug/network logs.

9. **[P0][M] Test coverage expansion** (quality)
   - Add unit tests for ViewModels/repositories and Compose UI tests for login/home/profile flows.

10. **[P1][M] Real notification events** (feature/fix)
    - Replace debug sample notifications with domain events (device offline, automation fired, login/session issues).

11. **[P1][M] Search and filter in Home** (feature)
    - Add search by device name and filters by room, state, and integration.

12. **[P1][M] Room-based grouping** (feature)
    - Group the device grid by room with collapsible sections.

13. **[P1][M] Capability-aware controls** (feature/fix)
    - Detect device capabilities and hide unsupported controls to avoid failed actions.

14. **[P1][M] Network diagnostics screen** (feature)
    - Show server reachability, local LAN/Wi-Fi checks, and quick troubleshooting actions.

15. **[P1][M] Localization support** (feature)
    - Move user-facing strings to resources and prepare for multi-language support.

16. **[P1][M] Analytics and event telemetry** (feature/ops)
    - Track key flows (login success/failure, device toggle latency, sync failures) for product health.

17. **[P1][L] WiZ device discovery onboarding** (feature)
    - Add an in-app LAN scan flow using `WizDiscovery` so users can discover and pair devices without manual server-side setup.

18. **[P1][L] Device details screen** (feature)
    - Add a dedicated screen with richer controls (brightness, color temp, RGB where supported).

19. **[P1][L] Offline sync queue for toggles** (feature/fix)
    - Queue device toggles while offline and retry when connectivity returns instead of only optimistic rollback.

20. **[P1][L] Multi-home switching** (feature)
    - Turn the current top-bar home label into a real selector with persisted active-home state.

21. **[P2][M] Import/export app settings** (feature)
    - Let users back up and restore non-sensitive preferences and layout settings.

22. **[P2][L] Background sync with WorkManager** (feature)
    - Periodically refresh device states and notifications in the background.

23. **[P2][L] Push notifications (FCM)** (feature)
    - Add Firebase Cloud Messaging for real-time alerts when app is in background.

24. **[P2][L] Tablet and foldable adaptive layouts** (feature)
    - Improve Home/Profile/Notifications for larger screens with responsive panes.

25. **[P2][L] Scenes and automations** (feature)
    - Add user-defined scenes (for example, Movie/Night) and schedule triggers.
