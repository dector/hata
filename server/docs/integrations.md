# Device integrations

Devices store integration details in two fields:

- `integration_id` — integration identifier used by the server to choose a controller.
- `integration_data` — optional JSON object with integration-specific payload.

## Supported integrations

### WiZ

Supported `integration_id` values:

- `wiz`

Payload:

```json
{
  "ip": "192.168.1.41"
}
```

Fields:

- `ip` string, required — LAN IP address of the WiZ light.

Example device row:

```json
{
  "id": "d-office-top",
  "name": "Top Light",
  "integration": {
    "id": "wiz-1",
    "data": {
      "ip": "192.168.1.41"
    }
  },
  "state": "on"
}
```

Notes:

- Device control currently supports only `on` and `off` states.
- Extra payload fields are ignored by the current WiZ controller.

## Unsupported integrations

Any other `integration_id` can be stored, but physical control will fail with `unsupported-integration` until a controller is implemented for it.
