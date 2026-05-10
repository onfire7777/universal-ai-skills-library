$ErrorActionPreference = "Stop"

$name = "lightpanda-ai-browser"
$image = "lightpanda/browser:nightly"

docker pull $image | Out-Null

$existing = docker ps -a --filter "name=^/$name$" --format "{{.Names}}"
if ($existing -eq $name) {
  $running = docker ps --filter "name=^/$name$" --format "{{.Names}}"
  if ($running -ne $name) {
    docker start $name | Out-Null
  }
} else {
  docker run -d `
    --name $name `
    -p 127.0.0.1:9222:9222 `
    -e LIGHTPANDA_DISABLE_TELEMETRY=true `
    --restart unless-stopped `
    $image /bin/lightpanda serve --host 0.0.0.0 --advertise-host 127.0.0.1 --port 9222 --log-level info --obey-robots | Out-Null
}

Invoke-RestMethod -Uri "http://127.0.0.1:9222/json/version" -TimeoutSec 10 | ConvertTo-Json -Depth 4
