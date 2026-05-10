# OpenUI README Map

Load only the reference needed for the task.

## Root

- `root/README.md` — OpenUI overview, quick start, package table, repository structure, Agent Skill install commands.
- `root/package.json` — workspace metadata, npm keywords, repository, license.
- `root/LICENSE` — MIT license.
- `root/CONTRIBUTING.md` and `root/SECURITY.md` — contribution and vulnerability reporting process.

## Packages

- `packages/lang-core/README.md` — framework-agnostic OpenUI Lang parser, streaming parser, prompt generation, runtime evaluator, tool specs, merge helpers.
- `packages/react-lang/README.md` — React component definitions with Zod, `createLibrary`, `library.prompt()`, `<Renderer />`, parser errors, hooks, form validation.
- `packages/react-headless/README.md` — chat state provider, thread/message hooks, streaming adapters, message format conversion.
- `packages/react-ui/README.md` — prebuilt chat layouts, built-in component libraries, theming, component imports, CSS import requirements.
- `packages/openui-cli/README.md` — `create` and `generate` CLI workflows, options, prompt/JSON Schema generation.
- `packages/svelte-lang/README.md` — Svelte 5 OpenUI Lang runtime.
- `packages/vue-lang/README.md` — Vue 3 OpenUI Lang runtime.
- `packages/browser-bundle/README.md` — browser bundle API and size notes.
- `packages/react-email/README.md` — OpenUI-driven React Email generation and rendering.

## Examples

- `examples/openui-chat/README.md` — scaffolded Next.js chat app and Docker usage.
- `examples/openui-dashboard/README.md` — dashboard example and OpenAI-compatible provider environment variables.
- `examples/vercel-ai-chat/README.md` — Vercel AI SDK chat integration.
- `examples/shadcn-chat/README.md` — shadcn component library integration.
- `examples/supabase-chat/README.md` — Supabase authentication, persistence, and realtime chat.
- `examples/svelte-chat/README.md` — Svelte chat setup and components.
- `examples/openui-react-native/README.md` — React Native app plus backend prompt generation.
- `examples/form-generator/README.md` — custom component library for generated forms.
- `examples/multi-agent-chat/README.md` — multi-agent orchestrator and assistant rendering pattern.

## Docs And Benchmarks

- `docs/README.md` — Fumadocs/Next.js docs site structure and local docs commands.
- `benchmarks/README.md` — OpenUI Lang token benchmark methodology, sample generation, offline benchmark report.
