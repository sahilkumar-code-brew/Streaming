---
type: registry
owner: architect
last-updated: 2026-10-27
---
# Tech Stack (AUTHORITATIVE)

> Single source of truth for technology choices. Reference it: `registries/tech-stack.md`

| Layer | Choice | Notes |
|-------|--------|-------|
| Entitlement / control plane | Go (`net/http`, go-redis) | No heavy ORMs on playback path |
| Gateway | Envoy | |
| Hot store | Redis Cluster | Tokens, subscriptions, steering weights |
| Billing store | Postgres | Offline billing only, never on playback path |
| Telemetry | Kafka -> ClickHouse | QoE beacons |
| Ingest / packaging | SRT -> CMAF fMP4 | |
| Delivery | Varnish origin shield, multi-CDN | Akamai 40 / CloudFront 40 / Fastly 20 |
| Mobile | React Native (TypeScript strict), Zustand | Media3 (Android), AVPlayer (iOS) native modules |
| Admin web | Next.js App Router, Tailwind, TanStack Query | |
| Streaming profiles | LL-HLS + DASH (CMAF profile) from shared CMAF segments | Targets in `features/live-playback.md` |
| Encode / package | FFmpeg encode; packager with LL-HLS partial segments | Confirm packager LL-HLS support before choosing |
| DRM | Widevine (Android), FairPlay (iOS) via drm-license-proxy | |
| API envelope | Org CB envelope, camelCase, UUID v4 | D-008 |
| Runtime | Kubernetes (EKS), Helm | |
| Infra as code | Terraform | |
| CI | GitHub Actions | Runs `scripts/kb_lint.py` and code checks |
| Observability | OpenTelemetry -> Prometheus, Grafana, Tempo | |
| Go libraries | `log/slog`, envconfig, go-playground/validator, goose migrations, testify, testcontainers, mockery, golangci-lint | Org go-service profile |
| Mobile libraries | Expo (prebuild, dev client), Expo Router, TanStack Query, Zod, MMKV, expo-secure-store | D-011; org react-native profile |
| Mobile testing | Jest + React Native Testing Library, MSW, Detox | Org react-native profile |
| Admin libraries | shadcn/ui, Lucide, React Hook Form + Zod, Zod-validated env | Org react-web profile |
| Admin testing | Vitest + React Testing Library, MSW, Playwright, axe-core | Org react-web profile |
| Events | CloudEvents 1.0 on Kafka | D-012 |
