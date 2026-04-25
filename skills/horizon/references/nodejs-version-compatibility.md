# Node.js Version Compatibility

Feature availability by Node.js version for backend modernization.

## LTS Timeline

| Version | Status | Active Support | Maintenance | EOL |
|---------|--------|----------------|-------------|-----|
| 18.x | Maintenance LTS | 2022-10 to 2023-10 | 2023-10 to 2025-04 | 2025-04 |
| 20.x | Active LTS | 2023-10 to 2024-10 | 2024-10 to 2026-04 | 2026-04 |
| 22.x | Current | 2024-10 (LTS) | 2025-10 to 2027-04 | 2027-04 |

## Feature Matrix

| Feature | Node 18 | Node 20 | Node 22 | Replaces |
|---------|---------|---------|---------|----------|
| Native `fetch` | ✅ | ✅ | ✅ | node-fetch, axios |
| Native test runner | ✅ | ✅ | ✅ | jest, mocha |
| `--watch` mode | ✅ | ✅ | ✅ | nodemon |
| `crypto.randomUUID` | ✅ | ✅ | ✅ | uuid |
| `structuredClone` | ✅ | ✅ | ✅ | lodash.cloneDeep |
| `.env` file loading | ❌ | ✅ | ✅ | dotenv |
| Native WebSocket | ❌ | ❌ | ✅ | ws |
| Permission model | ❌ | ✅ (exp) | ✅ | - |
| Single executable | ❌ | ✅ (exp) | ✅ | pkg |
| ESM by default | ✅ | ✅ | ✅ | - |
| Top-level await | ✅ | ✅ | ✅ | - |

## Upgrade Path Recommendations

### Node.js Upgrade Checklist

**From 16.x to 18.x:**
- [ ] Replace node-fetch with native fetch
- [ ] Update OpenSSL-dependent code (v3 changes)
- [ ] Review V8 engine changes
- [ ] Test npm workspaces compatibility

**From 18.x to 20.x:**
- [ ] Remove dotenv (use --env-file)
- [ ] Update to new test runner if desired
- [ ] Enable permission model for security
- [ ] Review experimental features used

**From 20.x to 22.x:**
- [ ] Replace ws with native WebSocket
- [ ] Consider single executable apps
- [ ] Review TypeScript 5.x compatibility
- [ ] Test with updated V8 engine

## package.json Engine Specification

```json
{
  "engines": {
    "node": ">=20.0.0",
    "npm": ">=10.0.0"
  }
}
```
