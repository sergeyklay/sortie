# Jira REST API — Adapter Research Notes

> Jira Cloud REST API v3, researched March 2026.
> Reference for implementing the Jira `TrackerAdapter`.

---

## Authentication

Jira supports several authentication methods depending on hosting environment and use case.

### Basic auth with API token (Cloud) — recommended for Sortie

The standard method for scripts and service integrations. Uses the user's email and an
API token generated from their Atlassian account profile.

- Generate a token at `https://id.atlassian.com/manage/api-tokens`.
- Header: `Authorization: Basic <base64(email:api_token)>`

This is the recommended method for Sortie because it runs as a background service, not
an interactive user-facing application.

### OAuth 2.0 (Cloud)

The recommended method for external apps accessing Jira on behalf of users. Uses the
authorization code grant type (3LO — Three-Legged OAuth). More secure as it restricts
scope and doesn't require sharing user credentials.

Not suitable for Sortie — OAuth 2.0 requires an interactive authorization flow and
callback URL, which a headless background service cannot provide.

### Personal Access Tokens (Data Center / Server)

PATs act as a secure alternative to Basic Auth passwords, behaving like bearer tokens.

- Header: `Authorization: Bearer <your_pat>`
- Primarily available in Jira Data Center and Server.
- Recently introduced to certain Cloud contexts.

Relevant only if Sortie adds Data Center support in the future.

### Config mapping

| Config field       | Value                                           |
| ------------------ | ----------------------------------------------- |
| `tracker.endpoint` | `https://<site>.atlassian.net` (no trailing /)  |
| `tracker.api_key`  | `email:api_token` — adapter splits on first `:` |
| `tracker.project`  | Jira project key, e.g. `SORT`                   |

Encoding `email:token` in a single field follows curl convention (`-u email:token`) and avoids
adding Jira-specific config keys to the core schema.

**CAPTCHA caveat:** After several failed logins Jira triggers CAPTCHA and returns
`X-Seraph-LoginReason: AUTHENTICATION_DENIED`. The adapter should detect this header and
return `tracker_auth_error`.

---

## Endpoints

Each `TrackerAdapter` operation maps to a Jira REST v3 endpoint.

### 1. `FetchCandidateIssues` → `GET /rest/api/3/search`

JQL:

```
project = "<KEY>" AND status IN ("<state1>", "<state2>", ...) ORDER BY priority ASC, created ASC
```

Query params: `jql`, `fields`, `maxResults`, `nextPageToken`

Request only needed fields:
`summary`, `status`, `priority`, `labels`, `assignee`, `issuetype`, `parent`,
`issuelinks`, `created`, `updated`, `description`

Does **not** request `comment` (separate call — comments use a dedicated endpoint).

Note: `POST /rest/api/3/search` also accepts JQL in the request body and avoids URI length
limits for very long queries. However, POST uses offset-based pagination and Atlassian
recommends the GET endpoint with cursor-based pagination. Sortie's JQL queries are short
enough for GET.

### 2. `FetchIssueByID` → `GET /rest/api/3/issue/{issueIdOrKey}`

Query param `fields` to select specific fields. Returns a single issue with full detail.

The `description` field uses **ADF** (Atlassian Document Format) — a JSON tree, not plain text.
Must be flattened (see ADF section below).

### 3. `FetchIssuesByStates` → `GET /rest/api/3/search`

JQL:

```
project = "<KEY>" AND status IN ("<state1>", ...) ORDER BY created ASC
```

Same endpoint as candidate fetch, different JQL. Used for startup terminal cleanup.
Paginate to fetch all matching issues.

### 4. `FetchIssueStatesByIDs` → `GET /rest/api/3/search`

JQL:

```
key IN ("PROJ-1", "PROJ-2", ...) ORDER BY key ASC
```

Request only `status` field to minimize payload. Used for active-run reconciliation.

Note: `id IN (...)` uses numeric internal IDs; `key IN (...)` uses project-prefixed keys.

### 5. `FetchIssueComments` → `GET /rest/api/3/issue/{issueIdOrKey}/comment`

Query params: `startAt`, `maxResults`, `orderBy`

Response: `{ startAt, maxResults, total, comments: [...] }`

Comment body uses ADF — must be flattened. Each comment has `id`, `author.displayName`,
`body` (ADF), `created`, `updated`.

### Transitions (reference only)

Sortie is a tracker **reader** — state transitions are handled by the coding agent,
not the orchestrator. These endpoints are documented for reference only.

- `GET /rest/api/3/issue/{issueIdOrKey}/transitions` — lists available transitions for an
  issue based on the current user's permissions and workflow rules.
- `POST /rest/api/3/issue/{issueIdOrKey}/transitions` — executes a transition, moving the
  issue to a new status. Request body: `{ "transition": { "id": "<transition_id>" } }`

If an optional `tracker_api` client-side tool extension is implemented, these endpoints
would be relevant.

---

## Field Mapping

`domain.Issue` field → Jira REST response path:

| `domain.Issue` field | Jira field                        | Notes                           |
| -------------------- | --------------------------------- | ------------------------------- |
| `ID`                 | `id` (string)                     | Numeric ID as string            |
| `Identifier`         | `key` (string)                    | e.g. `"PROJ-123"`               |
| `Title`              | `fields.summary`                  |                                 |
| `Description`        | `fields.description` (ADF)        | Flatten ADF → plain text        |
| `Priority`           | `fields.priority.id`              | Parse to int; non-numeric → nil |
| `State`              | `fields.status.name`              | Preserve original casing        |
| `BranchName`         | —                                 | Not natively available          |
| `URL`                | `{endpoint}/browse/{key}`         | Constructed                     |
| `Labels`             | `fields.labels` (string array)    | Lowercase each                  |
| `Assignee`           | `fields.assignee.displayName`     | Nil → empty string              |
| `IssueType`          | `fields.issuetype.name`           |                                 |
| `Parent`             | `fields.parent.id`, `.parent.key` | Nil when absent                 |
| `Comments`           | Separate endpoint                 | ADF → plain text                |
| `BlockedBy`          | `fields.issuelinks[]` (filtered)  | See blocker extraction below    |
| `CreatedAt`          | `fields.created` (ISO-8601)       |                                 |
| `UpdatedAt`          | `fields.updated` (ISO-8601)       |                                 |

### Blocker extraction from `issuelinks`

The "Blocks" link type has `type.inward = "Blocked by"` and `type.outward = "Blocks"`.

To find issues blocking the current issue, filter for links where:

- `type.name == "Blocks"` AND `outwardIssue` is present (meaning another issue blocks this one)
- Extract `outwardIssue.key` as the blocker identifier.

Note: The link type name "Blocks" may be renamed by Jira admins.

### ADF (Atlassian Document Format) flattening

Jira v3 returns `description` and comment `body` as ADF JSON:

```json
{
  "type": "doc",
  "version": 1,
  "content": [
    {
      "type": "paragraph",
      "content": [{ "type": "text", "text": "Hello world" }]
    }
  ]
}
```

The adapter must recursively walk the tree and extract all `text` node values, joining
paragraphs with newlines. Without this, `Description` and comment `Body` would be raw JSON.

**v2 API alternative:** The v2 API (`/rest/api/2/...`) returns `description` and comment
`body` as rendered HTML or plain text instead of ADF. If ADF flattening proves too complex,
the adapter could use v2 endpoints for these fields. However, v3 is the current API and
ADF flattening gives the adapter full control over text extraction.

---

## Pagination

### Search endpoint (`GET /rest/api/3/search`) — cursor-based

- First request: omit `nextPageToken`, set `maxResults` (recommend `50`).
- Subsequent requests: pass the `nextPageToken` from the previous response.
- Stop when `nextPageToken` is absent from the response.

Note: `POST /rest/api/3/search` uses offset-based (`startAt`/`total`) but is deprecated.
Use `GET` with cursor-based pagination.

### Comment endpoint — offset-based

- `startAt` (0-indexed), `maxResults` (default 50)
- Response includes `total`. Continue while `startAt + len(comments) < total`.

---

## Rate Limiting

Jira Cloud enforces three independent rate limiting systems:

### 1. Points-based quota (per-hour)

- Each call consumes points: base 1 + object costs (e.g., single issue GET = 2 points).
- Default quota: **65,000 points/hour** (resets at top of UTC hour).
- API-token traffic may be exempt from points-based limits (as of March 2026).

### 2. Burst rate limits (per-second, per-endpoint)

- Token bucket algorithm per endpoint per tenant.
- Defaults: GET 100 req/s, POST 100 req/s, PUT 50 req/s, DELETE 50 req/s.
- `GET /rest/api/3/issue/{id}`: 150 req/s burst bucket.
- `GET /rest/api/3/search`: 100 req/s.

### 3. Per-issue write limits

- 20 writes/2s, 100 writes/30s per issue.
- **Not relevant** — Sortie is read-only from the tracker.

### 429 handling

All limits return `429 Too Many Requests` with:

| Header                  | Value                                               |
| ----------------------- | --------------------------------------------------- |
| `Retry-After`           | Seconds to wait (integer)                           |
| `X-RateLimit-Remaining` | Remaining capacity                                  |
| `X-RateLimit-Reset`     | ISO-8601 reset timestamp                            |
| `RateLimit-Reason`      | `jira-quota-global-based`, `jira-burst-based`, etc. |

**Adapter guidance:** Respect `Retry-After` as minimum delay. Exponential backoff with jitter
(base 2s, max 30s, jitter ±30%). Map 429 → `tracker_api_error` with `Retry-After` preserved.

---

## Error Mapping

HTTP status → error category:

| HTTP Status | Condition                  | Error Category            |
| ----------- | -------------------------- | ------------------------- |
| 200         | Success                    | —                         |
| 400         | Bad JQL, invalid request   | `tracker_payload_error`   |
| 401         | Invalid/expired token      | `tracker_auth_error`      |
| 403         | Insufficient permissions   | `tracker_auth_error`      |
| 404         | Issue/resource not found   | `tracker_api_error`       |
| 429         | Rate limited               | `tracker_api_error`       |
| 5xx         | Server error               | `tracker_transport_error` |
| TCP/DNS     | Network failure            | `tracker_transport_error` |
| —           | JSON decode failure on 200 | `tracker_payload_error`   |
| —           | CAPTCHA (X-Seraph header)  | `tracker_auth_error`      |

---

## Config Notes

- **`tracker.api_key` format:** `email:api_token` — split on first `:`.
- **`tracker.endpoint`:** Site URL without trailing slash or path. Adapter appends `/rest/api/3/...`.
- **`tracker.project`:** Jira project key used in all JQL queries.
- **`tracker.active_states`:** Common defaults: `["Backlog", "Selected for Development", "In Progress"]`.
- **`tracker.terminal_states`:** Common defaults: `["Done", "Cancelled"]`.
- **JQL quoting:** Always quote string values in JQL to handle special characters in state names.
- **Network timeout:** 30,000 ms.
