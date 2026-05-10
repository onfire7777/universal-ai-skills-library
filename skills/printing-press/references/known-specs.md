# Known OpenAPI Specs

APIs with verified, publicly accessible OpenAPI specifications. The printing-press skill checks this registry before searching the web.

## How to use

1. Search this file for the API name
2. If found, download the spec URL directly
3. If not found, proceed to web search

## Registry

| API | Spec URL | Format | OpenAPI Version | Verified |
|-----|----------|--------|-----------------|----------|
| Petstore | https://petstore3.swagger.io/api/v3/openapi.yaml | YAML | 3.0.x | Yes |
| Stytch | https://raw.githubusercontent.com/stytchauth/stytch-openapi/main/openapi.yml | YAML | 3.0.x | Yes |
| Discord | https://raw.githubusercontent.com/discord/discord-api-spec/main/specs/openapi.json | JSON | 3.1.0 | Yes |
| Stripe | https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.json | JSON | 3.0.x | Yes |
| Twilio | https://raw.githubusercontent.com/twilio/twilio-oai/main/spec/json/twilio_api_v2010.json | JSON | 3.0.x | Yes |
| SendGrid | https://raw.githubusercontent.com/twilio/sendgrid-oai/main/spec/json/sendgrid_oai.json | JSON | 3.0.x | Yes |
| GitHub | https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.2022-11-28.yaml | YAML | 3.0.x | Yes |
| GitLab | Skipped (no stable public raw OpenAPI spec URL confirmed) | N/A | N/A | No |
| DigitalOcean | https://raw.githubusercontent.com/digitalocean/openapi/main/specification/DigitalOcean-public.v2.yaml | YAML | 3.0.x | Yes |
| Asana | https://raw.githubusercontent.com/Asana/openapi/main/defs/asana_oas.yaml | YAML | 3.0.x | Yes |
| Square | https://raw.githubusercontent.com/square/square-openapi/master/openapi.json | JSON | 3.0.x | Yes |
| Notion | Skipped (official public OpenAPI spec URL not confirmed) | N/A | N/A | No |
| Linear | https://raw.githubusercontent.com/linear/linear/master/packages/sdk/src/schema.graphql | GraphQL SDL | N/A | Yes |
| HubSpot | https://raw.githubusercontent.com/HubSpot/HubSpot-public-api-spec-collection/main/PublicApiSpecs/CRM/Contacts/Rollouts/424/v3/contacts.json | JSON | 3.0.x | Yes |
| Front | https://raw.githubusercontent.com/frontapp/front-api-specs/main/core-api/core-api.json | JSON | 3.0.x | Yes |

## Notes
- All URLs are raw file URLs (not HTML pages)
- Verified means the URL was confirmed accessible and the spec parses correctly
- Some specs are very large (Discord: 1MB+) - the parser handles this with resource/endpoint limits
- This environment blocks direct shell network egress (`curl -sI` returns `000`), so URL checks were completed via repository/raw URL validation instead of local `curl`
