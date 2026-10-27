---
type: architecture
owner: backend-lead
surfaces: [backend, mobile, admin]
last-verified: 2026-10-27
verified-by: docs
---
# API Contracts

> Ownership and per-surface status: `../registries/contracts.md`. Every `## C-NNN` section here must match a row there (`scripts/check_kb_sync.sh`). All paths live under `/api/v1/` (D-009). Envelope, error codes, casing, IDs and timestamps follow the org CB conventions (D-008).

## Shared types
| Type | Definition |
|------|-----------|
| `StreamId` | string slug, `^[a-z0-9-]{3,64}$` |
| `UserId` | string, UUID v4 |
| `DeviceId` | string, UUID v4, generated once per install |
| `RequestId` | string, UUID v4 |
| `Timestamp` | string, ISO 8601 UTC, e.g. `2026-10-27T12:00:00.000Z` |
| `CdnPathway` | `"akamai" \| "cloudfront" \| "fastly"`. Unknown value on a client: ignore that pathway, keep the rest |
| `StreamState` | `"scheduled" \| "live" \| "interrupted" \| "ended"`. Transitions: scheduled->live, live->interrupted, interrupted->live, live->ended, interrupted->ended. Only stream-ops-service writes it |
| `Role` | `"viewer" \| "ops" \| "steering-admin"` |

JSON fields are camelCase. Required headers on every mobile and admin request: `X-Device-Id: DeviceId`, `X-App-Version: semver`, `X-Request-Id: RequestId`. Fields not marked optional are required and non-null.

## Response envelope (D-008)
Success: `{ "success": true, "data": {}, "meta": { "timestamp": "Timestamp", "requestId": "RequestId" } }`. Below, each contract shows only `data`.
Exceptions: C-002 returns the HLS steering document unwrapped, C-004 is a static CDN file, C-003 returns `202` with no body.

Error:
```json
{
  "success": false,
  "error": {
    "code": "ENTITLEMENT_NOT_SUBSCRIBED",
    "message": "string, human readable, never parsed",
    "details": [{ "field": "string", "message": "string", "code": "string" }],
    "retryable": false,
    "retryAfterMs": "integer >= 0, present only when retryable"
  },
  "meta": { "timestamp": "Timestamp", "requestId": "RequestId" }
}
```
`details` is optional and only used with `VALIDATION_ERROR`.

| HTTP | code | Source | retryable | Client action |
|------|------|--------|-----------|---------------|
| 400 | `VALIDATION_ERROR` | org | no | Show field errors |
| 400 | `INVALID_REQUEST` | org | no | Report bug |
| 401 | `UNAUTHORIZED` | org | no | Re-login |
| 401 | `SESSION_EXPIRED` | domain | no | Refresh via C-008, then re-login |
| 401 | `OTP_INVALID` | domain | no | Re-enter OTP |
| 403 | `FORBIDDEN` | org | no | Admin: hide action |
| 403 | `ENTITLEMENT_NOT_SUBSCRIBED` | domain | no | Paywall |
| 403 | `ENTITLEMENT_GEO_BLOCKED` | domain | no | Region message |
| 403 | `ENTITLEMENT_DEVICE_LIMIT` | domain | no | Reserved (R-009) |
| 404 | `NOT_FOUND` | org | no | Back to catalog |
| 409 | `CONFLICT` | org | no | Refetch |
| 409 | `STREAM_NOT_LIVE` | domain | no | Show schedule |
| 409 | `STREAM_STATE_TRANSITION_INVALID` | domain | no | Refetch state |
| 412 | `VERSION_CONFLICT` | domain | no | Refetch, re-apply |
| 426 | `CLIENT_UPGRADE_REQUIRED` | domain | no | Force-update screen |
| 429 | `RATE_LIMITED` | org | yes | Wait `retryAfterMs` |
| 500 | `INTERNAL_ERROR` | org | yes | Full-jitter backoff (D-010) |
| 503 | `SERVICE_UNAVAILABLE` | org | yes | Full-jitter backoff, cap 30s (D-006) |

Clients branch on `code` only. Unknown code: retryable if status >= 500, otherwise not.

## C-001 Entitlement
`POST /api/v1/streams/{stream_id}/entitlement`. Headers: `Authorization: Bearer <session_jwt>`. Body: none. Client timeout 2s; server budget 50ms. Rate limit 10/min per user.
```json
{
  "playbackToken": "string, compact JWS EdDSA, <= 1024 bytes",
  "tokenExpiresAt": "Timestamp",
  "refreshAfterS": 240,
  "manifestUrl": "string, https URL",
  "steeringManifestUrl": "string, https URL",
  "telemetrySampleRate": "number 0.0..1.0"
}
```
Token claims: `sub` UserId, `sid` StreamId, `did` sha256(DeviceId), `kid`, `iat`, `exp` = iat + 300. Edge accepts up to 30s past `exp` (clock skew); edge holds current and previous `kid` public keys (R-004).
Errors: `SESSION_EXPIRED`, `ENTITLEMENT_NOT_SUBSCRIBED`, `ENTITLEMENT_GEO_BLOCKED`, `NOT_FOUND`, `STREAM_NOT_LIVE` (state not live or interrupted), `RATE_LIMITED`, `SERVICE_UNAVAILABLE` (Redis miss, D-006).
Refresh: client re-calls C-001 at `refreshAfterS + rand(0, 30)` s; keeps playing on the old token until `exp + 30s`. Token reaches segment URLs via origin manifest rewrite (D-007).

## C-002 Steering
`GET /api/v1/steering.json?_HLS_pathway=<CdnPathway>&_HLS_throughput=<int bps>`. HLS content steering (draft-pantos-hls-rfc8216bis). `Cache-Control: max-age=` equal to TTL.
```json
{
  "VERSION": 1,
  "TTL": 30,
  "RELOAD-URI": "string, optional, https URL",
  "PATHWAY-PRIORITY": ["CdnPathway", "..."]
}
```
`PATHWAY-PRIORITY` non-empty, no duplicates. On failure clients keep the last list and retry at TTL with jitter.

## C-003 Telemetry
`POST /api/v1/telemetry/qoe`. Policy (interval, batch, sampling): `../features/live-playback.md`. Max body 16 KB.
```json
{
  "deviceId": "DeviceId",
  "streamId": "StreamId",
  "appVersion": "semver",
  "events": [
    { "ts": "Timestamp", "type": "stallStart | stallEnd | bitrateChange | cdnSwitch | error",
      "cdn": "CdnPathway", "bufferMs": "integer >= 0", "bitrateKbps": "integer >= 0, optional",
      "errorCode": "string, optional" }
  ]
}
```
`events` length 1..20. Responses: `202` empty, `429` with `retryAfterMs`. Telemetry is lossy: drop a batch after 2 failed retries.

## C-004 Live scores
`GET /api/v1/scores/{stream_id}.json`, served from CDN. Origin: `Cache-Control: public, max-age=1, stale-while-revalidate=1`.
```json
{
  "streamId": "StreamId",
  "seq": "integer, monotonic per stream",
  "updatedAt": "Timestamp",
  "matchState": "pre | live | break | ended",
  "teams": [{ "id": "string", "name": "string", "score": "string, e.g. 245/6", "overs": "string, optional" }],
  "summary": "string, <= 140 chars"
}
```
Client polls every 2s + rand(0, 500) ms; drops responses with `seq` <= last seen; backs off with full jitter after errors; stops when `matchState` is `ended`.

## C-005 Steering weights write
`PUT /api/v1/admin/steering/weights`. Headers: `Authorization` (role `steering-admin`), `If-Match: <version>`, `Idempotency-Key: UUIDv4`.
```json
{ "weights": { "akamai": 40, "cloudfront": 40, "fastly": 20 }, "confirmLargeShift": false }
```
Each weight integer 0..100, every known `CdnPathway` present, sum = 100. A shift > 50 points on any pathway needs `confirmLargeShift: true`. Response `200` data `{ "version": integer, "weights": {}, "updatedBy": "UserId", "updatedAt": "Timestamp" }`. Errors: `VALIDATION_ERROR`, `FORBIDDEN`, `VERSION_CONFLICT`. Every write appends an audit row.

## C-006 Stream state
`PUT /api/v1/admin/streams/{stream_id}/state`. Headers: `Authorization` (role `ops`), `If-Match: <version>`.
```json
{ "state": "StreamState", "reason": "string, 1..200 chars" }
```
Response `200` data `{ "streamId": "StreamId", "state": "StreamState", "version": integer, "updatedAt": "Timestamp" }`. Errors: `NOT_FOUND`, `STREAM_STATE_TRANSITION_INVALID`, `VERSION_CONFLICT`, `FORBIDDEN`.

## C-007 NOC rebuffer query
Not HTTP. admin-web server reads ClickHouse `client_qoe_telemetry` as a read-only user: `max_execution_time` 5s, `max_result_rows` 10000, window <= 1h, grouped by `cdn` per minute. Never from the browser.

## C-008 Session token
`POST /api/v1/auth/token`.
```json
{ "grantType": "otp", "phone": "string, E.164", "otp": "string, ^[0-9]{6}$" }
{ "grantType": "refreshToken", "refreshToken": "string, opaque, 64 chars" }
```
Response `200` data `{ "sessionJwt": "JWS EdDSA", "expiresAt": "Timestamp", "refreshToken": "string", "roles": ["Role"] }`. `sessionJwt` lives 15 min; refresh token 30 days, single use, rotated on every call; reuse of a rotated token revokes the whole family; every grant and reuse is audit-logged. Client refreshes at 80% of lifetime with jitter, or on `SESSION_EXPIRED`. Errors: `OTP_INVALID`, `SESSION_EXPIRED`, `RATE_LIMITED`.

## C-009 OTP send
`POST /api/v1/auth/otp`. Body `{ "phone": "string, E.164" }`. Response `202` empty. Rate limit 3 per 10 min per phone and per device. Errors: `VALIDATION_ERROR`, `RATE_LIMITED`.
