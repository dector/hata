# Hata

Minimalist home automation hub written in Go.

**Status:** Currently in development and not ready for usage.

## Links

- **Sources:** https://github.com/dector/hata
- **YT videos (rare):** [here](https://www.youtube.com/playlist?list=PL8DPygpcWgCJ_xVC5t-NxpUqQNv2QRRI2)

## Development

### Server dev mode (`ror dev`) and Tailscale port 4500

`server/ror.kdl` provides two server tasks:

- `ror dev` — live reload + local proxy (`127.0.0.1:4500 -> 127.0.0.1:4501`)
- `ror dev:server` — live reload without proxy (direct app on `127.0.0.1:4501`)

Why this setup exists:

- Air's built-in proxy binds to wildcard host (`0.0.0.0`) and cannot be forced to bind localhost-only.
- On this machine, port `4500` is exposed by Tailscale, so wildcard bind causes conflicts.
- Workaround: use Air only for rebuild/restart and run a tiny custom reverse proxy (`cmd/devproxy`) bound to `127.0.0.1:4500`.

## Roadmap

- [ ] Role-based access
- [ ] Device scanning
- [ ] REST API
- [ ] Mobile app
- [ ] Historical data tracking
- [ ] Configurable dashboards
- [ ] Automation scenarios
- [ ] Guest mode
