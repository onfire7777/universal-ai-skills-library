# Manus Project Skills gRPC-web API Reference

## Base URL

```
https://api.manus.im
```

## Authentication

All requests require a JWT session token from the Manus browser session.

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
Connect-Protocol-Version: 1
```

### Extracting the Token

Run in browser console at `manus.im` (while logged in):
```javascript
document.cookie.split(';').find(c => c.trim().startsWith('session_id=')).split('=').slice(1).join('=')
```

Or use `browser_console_exec` tool to extract it programmatically.

## Endpoints

### ListProjectSkills

List all skills in a project.

```
POST /skill.v1.ProjectSkillService/ListProjectSkills
```

Request:
```json
{"project_uid": "KHUqYo7c2G8ZvsByd4Dvwx"}
```

Response:
```json
{
  "projectSkills": [
    {
      "skill": {
        "id": "abc123",
        "name": "skill-name",
        "description": "..."
      }
    }
  ]
}
```

### UploadProjectSkill

Upload a skill zip to a project. The zip must contain a valid `SKILL.md` at root.

```
POST /skill.v1.ProjectSkillService/UploadProjectSkill
```

Request:
```json
{
  "project_uid": "KHUqYo7c2G8ZvsByd4Dvwx",
  "content": "<base64-encoded-zip-bytes>"
}
```

Response: `200 OK` on success.

### DeleteProjectSkill

Delete a skill from a project by skill ID.

```
POST /skill.v1.ProjectSkillService/DeleteProjectSkill
```

Request:
```json
{
  "project_uid": "KHUqYo7c2G8ZvsByd4Dvwx",
  "skill_id": "abc123"
}
```

Response: `200 OK` on success.

## Limits

- **500 skills per project** (hard limit, returns error on exceed)
- **Rate limiting**: 429 responses on rapid requests (~150ms delay recommended)
- **Zip format**: Must contain `SKILL.md` at root with valid YAML frontmatter

## Error Codes

| Status | Meaning |
|---|---|
| 200 | Success |
| 400 | Bad request (missing/invalid fields) |
| 401 | Unauthorized (bad/expired token) |
| 429 | Rate limited (retry with exponential backoff) |
| 500 | Server error (often `SKILL_MD_INVALID_FORMAT`) |

## Project UIDs

Get project UIDs by navigating to project settings in the Manus UI:
```
https://manus.im/app/project/<PROJECT_UID>
```

Or extract from sidebar links via browser console:
```javascript
document.querySelectorAll('a[href*="/app/project/"]')
```
