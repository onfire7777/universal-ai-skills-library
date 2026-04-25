# API Versioning Strategies

## Comparison

| Strategy | Pros | Cons | Example |
|----------|------|------|---------|
| URL Path | Simple, visible | URL pollution | `/v1/users` |
| Header | Clean URLs | Hidden version | `Accept: application/vnd.api.v1+json` |
| Query Param | Easy testing | Caching issues | `/users?version=1` |
| Content Negotiation | Standard-based | Client complexity | `Accept: application/json; version=1` |

**Recommendation:** URL Path versioning for simplicity and clarity.

## Deprecation Timeline

1. Announce deprecation (6 months before)
2. Add `Deprecation` header to responses
3. Add `Sunset` header with date
4. Monitor usage of deprecated version
5. Remove after sunset date

## Breaking vs Non-Breaking Changes

| Change | Breaking? |
|--------|-----------|
| Add optional field | No |
| Add new endpoint | No |
| New HTTP methods on existing endpoints | No |
| More permissive validation | No |
| Remove field | Yes |
| Rename field | Yes |
| Change field type | Yes |
| Add required field | Yes |
| Change URL structure | Yes |
| Stricter validation | Yes |
| Change authentication method | Yes |
| Change error response format | Yes |

## Version Migration Strategy

```markdown
## Version Migration Plan: v1 → v2

### Timeline
| Phase | Duration | Action |
|-------|----------|--------|
| Announcement | Week 1 | Notify consumers of v2 release |
| Parallel Operation | Weeks 2-12 | Both v1 and v2 available |
| Deprecation Notice | Week 8 | Add deprecation headers to v1 |
| v1 Sunset | Week 13 | v1 returns 410 Gone |

### Deprecation Headers
```http
Deprecation: true
Sunset: Sat, 01 Mar 2025 00:00:00 GMT
Link: </api/v2/users>; rel="successor-version"
```
```
