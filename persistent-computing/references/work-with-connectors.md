# Working with Connectors in Persistent Environments

User-configured connectors are not available by default in WebDev, cloud computers, or local computers connected via the Manus desktop client. The authentication mechanisms for connectors are sandboxed and do not automatically transfer to these environments.

When the user's project requires connectors, there are two primary solutions:

## Option 1: Extract and Reuse Credentials from Connectors

If a connector is already configured in the user's Manus account, its underlying credentials can sometimes be extracted and reused in persistent environments (cloud computers, WebDev, or local computers).

**How it works:**
1. Inspect what credentials the connector exposes in the default sandbox
2. Determine whether those credentials can be safely transferred to the target environment
3. Store them in the target environment's secrets or configuration
4. Use them directly in the application code

**What can be extracted:**
- **Fixed API Keys** (e.g., OpenAI API Key) - Static credentials that don't expire. They can be directly extracted and reused.
- **Access Tokens** - Some connectors expose access tokens that may expire. These work for temporary or debugging use, but for sustained availability there are two paths: either the user creates their own OAuth App (so they can manage refresh tokens themselves), or fall back to Option 2.

**What cannot be extracted:**
- **Refresh Tokens** - Never exposed in the sandbox. If a connector relies on refresh tokens to maintain access, Option 1 is not viable.
- **OAuth Flows** - Some connectors handle OAuth authentication internally. If the underlying credentials are not exposed, Option 1 is not viable.

**Limitations:**
- Not all connectors expose their underlying credentials
- Some credentials are temporary and have no path to renewal without user-side OAuth setup
- The application is responsible for securely storing and managing extracted credentials
- Connector type matters: MCP connectors and Custom API connectors may expose credentials differently

**Process:**
1. In the default sandbox, check environment variables or connector outputs to see what credentials are available
2. Determine if the credentials are fixed (safe to transfer) or temporary (will eventually fail)
3. If fixed, extract them and store them in the target environment's secrets
4. If temporary, discuss with the user whether they can set up their own OAuth App; otherwise, fall back to Option 2

## Option 2: Use the Manus API

The Manus API can be called from the application to create a task for Manus to execute. The task includes using the user's configured connectors as one of its steps. Manus runs the task (which may involve multiple steps, logic, and AI reasoning) and returns the results.

**Advantages:**
- **Handles credential complexity** - Works with OAuth-based connectors and automatic token refresh without exposing credentials to the application
- **Intelligent processing** - The Manus Agent can apply logic before and after connector calls (data validation, transformation, aggregation, multi-step workflows)
- **Simplified configuration** - No need for the user to set up their own OAuth Client or manage individual API keys

**Limitations:**
- **Costs Credits** - Each task execution consumes Manus Credits. This is an ongoing cost the user must accept.
- Requires a Manus API Key, which grants access to all of the user's configured connectors
- Depends on Manus API availability

**When to use Option 2:**
- Connector uses OAuth with refresh tokens (cannot extract credentials)
- Connector doesn't expose underlying credentials
- The application needs intelligent pre/post-processing of connector data
- The user prefers not to manage credentials in the application

**How it works:**
1. The user creates a Manus API Key at [Manus API Integration settings](https://manus.im/app?show_settings=integrations&app_name=api)
2. The application calls the Manus API to create a task, specifying the required connector IDs
3. Manus executes the task (including connector access) and returns the results
4. The application receives the result, optionally validated against a JSON Schema via Structured Output

**Important security note:** A Manus API Key grants access to all of the user's configured connectors. The application must protect it with the same care as any sensitive credential.

For implementation details and API reference, see the `manus-api` skill.

## Decision Guide

| Scenario | Recommended Option |
|----------|-------------------|
| Connector exposes a fixed API Key that can be extracted | Option 1 |
| Connector uses OAuth with refresh tokens | Option 2 |
| Connector doesn't expose any credentials | Option 2 |
| The application needs pre/post-processing logic | Option 2 |
| Minimizing external dependencies is a priority | Option 1 (if possible) |

## Implementation Notes

- **For cloud computers:** Store extracted credentials in environment variables or configuration files. They can be documented in `agents.md` (without exposing the sensitive values themselves).
- **For WebDev:** Use environment variables configured through the WebDev project settings, or implement Option 2 (Manus API) for connector access.
- **For local computers:** Extracted credentials can be stored in the local environment or configuration files.

Refer to the `manus-api` skill for specific endpoints and authentication details when choosing Option 2.
