$ErrorActionPreference = "Stop"

$name = "lightpanda-ai-browser"
$running = docker ps --filter "name=^/$name$" --format "{{.Names}}"
if ($running -eq $name) {
  docker stop $name | Out-Null
}
