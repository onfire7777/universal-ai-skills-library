@echo off
set LIGHTPANDA_DISABLE_TELEMETRY=true
docker run --rm -i -e LIGHTPANDA_DISABLE_TELEMETRY=true lightpanda/browser:nightly /bin/lightpanda mcp --obey-robots
