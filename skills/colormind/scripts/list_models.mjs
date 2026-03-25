#!/usr/bin/env node

const url = "http://colormind.io/list/";

const resp = await fetch(url, { method: "GET" });
if (!resp.ok) {
  const text = await resp.text().catch(() => "");
  console.error(`Colormind list failed (${resp.status}): ${text}`);
  process.exit(1);
}

const data = await resp.json();
console.log(JSON.stringify(data, null, 2));
