param()

$ErrorActionPreference = 'Stop'
$HomeDir = $env:USERPROFILE
$Root = Join-Path $HomeDir '.universal-ai-stack'
$HermesConfig = Join-Path $HomeDir '.hermes\config.yaml'
$Python = Join-Path $HomeDir '.hermes\hermes-agent\venv\Scripts\python.exe'

if (!(Test-Path -LiteralPath $HermesConfig)) {
  throw "Hermes config not found: $HermesConfig"
}
if (!(Test-Path -LiteralPath $Python)) {
  $Python = 'python'
}

$script = @'
from pathlib import Path
from datetime import datetime
import json
import yaml

path = Path.home() / ".hermes" / "config.yaml"
data = yaml.safe_load(path.read_text(encoding="utf-8")) or {}
backup = path.with_name(path.name + ".bak-universal-router-" + datetime.now().strftime("%Y%m%d-%H%M%S"))
backup.write_text(path.read_text(encoding="utf-8"), encoding="utf-8")

router_base = "http://127.0.0.1:18100/v1"
canonical_skills_root = str(Path.home() / "universal-ai-skills-library" / "skills")
router_provider = {
    "name": "Universal AI Stack Router",
    "base_url": router_base,
    "api_mode": "chat_completions",
    "key_env": "UNIVERSAL_AI_STACK_API_KEY",
    "default_model": "auto-coding",
    "models": [
        "auto-coding",
        "primary-api",
        "local-coding",
        "qwen3-coder-30b-a3b-q4",
        "kimi-k2.6-thinking",
        "qwen3-coder-next-q5",
        "qwen2.5-coder-32b-q4",
    ],
}

model = data.setdefault("model", {})
model.update({
    "default": "gpt-5.5",
    "provider": "openai-codex",
    "base_url": "https://chatgpt.com/backend-api/codex",
    "api_mode": "codex_responses",
    "context_length": 1000000,
})

providers = data.setdefault("providers", {})
if not isinstance(providers, dict):
    providers = {}
    data["providers"] = providers
for redundant in ("local-qwen3-coder", "local-qwen25-coder"):
    providers.pop(redundant, None)
providers["universal-router"] = router_provider

skills = data.setdefault("skills", {})
if not isinstance(skills, dict):
    skills = {}
    data["skills"] = skills
external_dirs = skills.setdefault("external_dirs", [])
if not isinstance(external_dirs, list):
    external_dirs = []
    skills["external_dirs"] = external_dirs
external_dirs = [str(p) for p in external_dirs if str(p) != canonical_skills_root]
external_dirs.insert(0, canonical_skills_root)
skills["external_dirs"] = external_dirs

data["fallback_providers"] = [{
    "provider": "universal-router",
    "model": "auto-coding",
    "base_url": router_base,
    "key_env": "UNIVERSAL_AI_STACK_API_KEY",
    "api_mode": "chat_completions",
}]

agent = data.setdefault("agent", {})
agent["reasoning_effort"] = "xhigh"
agent["service_tier"] = "fast"
agent["gateway_notify_interval"] = 180
agent["max_turns"] = 30
agent["max_iterations"] = 30

compression = data.setdefault("compression", {})
compression.update({
    "enabled": True,
    "threshold": 0.60,
    "target_ratio": 0.35,
    "protect_last_n": 12,
    "hygiene_hard_message_limit": 250,
})

cron = data.setdefault("cron", {})
cron["tick_interval_seconds"] = 600
cron["max_parallel_jobs"] = 1

aux = data.setdefault("auxiliary", {})

def route_aux(task, model_name="auto-coding", timeout=None):
    block = aux.setdefault(task, {})
    if not isinstance(block, dict):
        block = {}
        aux[task] = block
    block["provider"] = "universal-router"
    block["model"] = model_name
    block["base_url"] = router_base
    block["api_mode"] = "chat_completions"
    block.setdefault("extra_body", {})
    if timeout is not None:
        block["timeout"] = timeout
    block.pop("api_key", None)

route_aux("compression", "kimi-k2.6-thinking", 300)
aux["compression"]["context_length"] = 262144
for task, timeout in {
    "web_extract": 180,
    "session_search": 60,
    "skills_hub": 60,
    "approval": 60,
    "mcp": 60,
    "title_generation": 60,
    "triage_specifier": 120,
    "curator": 600,
}.items():
    route_aux(task, "auto-coding", timeout)

vision = aux.setdefault("vision", {})
if isinstance(vision, dict):
    vision["provider"] = "auto"

discord = data.setdefault("discord", {})
discord["require_mention"] = False
discord["auto_thread"] = True
discord["reactions"] = True
platforms = data.setdefault("platforms", {})
platforms.setdefault("discord", {})["enabled"] = True
api_server = platforms.setdefault("api_server", {})
api_server["enabled"] = True
api_server.setdefault("extra", {})["port"] = 8642
api_server.setdefault("extra", {})["host"] = "127.0.0.1"

path.write_text(yaml.safe_dump(data, sort_keys=False, allow_unicode=True), encoding="utf-8")

report = {
    "hermesConfig": str(path),
    "backup": str(backup),
    "primaryProvider": data.get("model", {}).get("provider"),
    "primaryModel": data.get("model", {}).get("default"),
    "fallbackProvider": data.get("fallback_providers", [{}])[0].get("provider"),
    "fallbackModel": data.get("fallback_providers", [{}])[0].get("model"),
    "compressionProvider": data.get("auxiliary", {}).get("compression", {}).get("provider"),
    "compressionModel": data.get("auxiliary", {}).get("compression", {}).get("model"),
    "discordRequireMention": data.get("discord", {}).get("require_mention"),
    "canonicalSkillsRoot": canonical_skills_root,
    "canonicalSkillsRootLinked": canonical_skills_root in data.get("skills", {}).get("external_dirs", []),
}
print(json.dumps(report, indent=2))
'@

$script | & $Python -
