---
type: architecture
owner: mobile-lead
surfaces: [mobile]
last-verified: 2026-10-26
verified-by: docs
---
# Mobile App (React Native)

- Media3/ExoPlayer on Android, AVPlayer on iOS, in native modules (LL-003).
- NativeModule exposes buffer telemetry (ms) and stall state to JS at low frequency.
- Jittered polling hook for live scores.
