# First-Time Setup (Before Image Generation)

## 1. First Trigger: Readiness First

Before the first real generation attempt, the agent should run a readiness pass silently:

```bash
npm run ensure-ready -- --project . --workflow <workflow>
```

This should happen before model choice. The purpose is:

- detect missing local runtime dependencies
- auto-bootstrap local script dependencies when possible
- confirm whether `IMAGE_GEN_API_KEY` is already available

If local dependencies are still missing after the check, the agent should ask the user for approval to install them on the user's behalf.

If `IMAGE_GEN_API_KEY` is missing, the agent should tell the user that image generation needs an API key and offer to configure it now. After approval, the agent should persist it to `.image-skills/image-generation/.env` or run `npm run setup -- --project . --workflow <workflow> --persist-api-key` when the key is already available in env. Do not ask the user to debug local setup manually before this readiness pass.

Treat the API key as a secret:

- prefer writing it locally on the user's behalf instead of making the user edit files
- do not ask the user to paste the key into normal chat if a local write path is available
- if the user does share the key, never echo it back or include it in follow-up messages

## 2. API Key

Get a key from the current gateway console at [WeryAI API Keys](https://weryai.com/api/keys), then configure only:

- environment variable **`IMAGE_GEN_API_KEY`** (recommended), or
- project/user `.image-skills/image-generation/.env` with:
  `IMAGE_GEN_API_KEY=<your Bearer token>`

Preferred agent path:

```bash
IMAGE_GEN_API_KEY=<secret> npm run setup -- --project . --workflow <workflow> --persist-api-key
```

## 3. Default Model Initialization

When no default model is configured yet, the **agent** should initialize the workspace with **Nano Banana 2** and clearly tell the user they can switch later at any time:

1. Start from the bundled starter registry (silent, no user involvement). If the agent needs a fresher list, refresh it from docs:

```bash
npm run discover-image-models -- --out .image-skills/image-generation/MODELS.json
```

2. Write the default config immediately:

```bash
mkdir -p .image-skills/image-generation
cat > .image-skills/image-generation/EXTEND.md << 'EOF'
---
version: 1
default_model: "GEMINI_3_1_FLASH_IMAGE"
default_quality: 2k
default_aspect_ratio: "1:1"
batch:
  max_workers: 4
---
EOF
```

3. Tell the user what happened:

> All set — I'll use **Nano Banana 2** as the default model for this workspace.
> If you want, I can switch models anytime later for different quality, style, or text-layout needs.

4. If the user wants a different model right now, run recommendation for the user's workflow:

```bash
npm run recommend-model -- --workflow <workflow> --role <role> --json
```

```bash
npm run recommend-model -- --workflow <workflow> --role <role> --json
```

5. Present alternatives only when needed:

> You're currently set to **Nano Banana 2** by default. If you'd like a better fit for this task, here are other strong options:
> 1. **GPT Image 1.5** — balanced and reliable; good for general images and text-heavy layouts
> 2. **Recraft V4 Pro** — stronger for infographic-style visuals and cleaner layout control
> 3. **WeryAI Image 2.0** — high-quality general image model, better for cover-style visuals
>
> Pick one if you want to switch now — otherwise I can keep Nano Banana 2. You can switch anytime.

6. If the user switches, update `EXTEND.md` and confirm:

> Updated — you're now defaulting to **GPT Image 1.5** instead of **Nano Banana 2**. Let me know anytime if you want to switch again or adjust defaults.

The agent **must not** ask the user to manually edit `EXTEND.md` or figure out the file format. The entire flow should feel like a guided conversation, not a config tutorial.

If a requested model is missing or looks stale, the agent should first refresh the local registry from WeryAI docs before asking the user for any platform-side details.

If the user later asks to switch models, the agent should update `EXTEND.md` in place.

## 4. Verify

```bash
npx -y bun scripts/main.ts --prompt "test" --image ./out/test.png --ar 1:1 --dry-run
```

Use `--dry-run` first to verify the request body is correct without consuming API quota. If the dry-run looks good, run without `--dry-run` to confirm end-to-end.

## 5. If the Request Is Still Vague

If the user's visual direction is underspecified, generate a brief before prompt writing:

```bash
npm run build-visual-brief -- --workflow cover --topic "Habit systems"
```

Use its question menu to guide the conversation. See `SKILL.md` § "If The Request Is Underspecified" when this file is bundled with the skill.
