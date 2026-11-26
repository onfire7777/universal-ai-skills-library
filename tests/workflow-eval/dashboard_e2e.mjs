// Run with SKILL_ROUTER_BIN and CHROME_BIN pointing to the binaries under test.
// Uses Chrome's private DevTools pipe and a temporary browser profile.
import assert from 'node:assert/strict';
import { execFileSync, spawn } from 'node:child_process';
import { createHash } from 'node:crypto';
import { mkdtemp, mkdir, readFile, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { StringDecoder } from 'node:string_decoder';
import { fileURLToPath, pathToFileURL } from 'node:url';

const repo = fileURLToPath(new URL('../../', import.meta.url));
const router = process.env.SKILL_ROUTER_BIN || 'skill-router';
const chrome = process.env.CHROME_BIN || 'google-chrome';
const prompt = 'Refresh the SaaS dashboard so it matches our existing brand and makes common tasks easier for keyboard users.';

function digest(bytes) {
  return createHash('sha256').update(bytes).digest('hex');
}

function workflow(method, request, env) {
  const stdout = execFileSync(router, ['skills', 'workflow', method], {
    cwd: repo, env, input: JSON.stringify(request), encoding: 'utf8', timeout: 240000,
    maxBuffer: 32 * 1024 * 1024,
  });
  return JSON.parse(stdout);
}

async function makePlan(env) {
  const specs = [
    ['brand', 'prerequisite', 'Read existing guidelines and extract the brand palette', 'brand',
      'inject brand context quick start', 'SKILL.md', 'Quick Start'],
    ['interface', 'deliverable', 'Apply the existing palette and native semantic controls to the dashboard ticket queue', 'ui-audit',
      'native HTML semantics clickable div button', 'rules/a11y-semantic-html-first.md', 'Prefer Native Semantics Before ARIA'],
    ['keyboard', 'verification', 'Check forward and reverse focus order, visible focus, and Enter activation', 'accessibility-testing',
      'manual accessibility testing keyboard focus visible tab order', 'SKILL.md', 'Manual Testing Checklist'],
  ];
  const requirements = [];
  for (const [id, kind, text, skill, query, file, title] of specs) {
    const result = workflow('search', { prompt, query, skill, limit: 100 }, env);
    const e = result.hits.find(({ element }) => element.file === file && element.title === title)?.element;
    assert.ok(e, 'Missing source section ' + skill + '/' + file + ' ' + title);
    const source = await readFile(path.join(repo, 'skills', skill, file));
    assert.equal(e.digest, digest(source));
    const { startLine, endLine } = e;
    const sourceSection = source.toString('utf8').split(/\r?\n/).slice(startLine - 1, endLine).join('\n');
    assert.ok(sourceSection.includes(title), 'Locator did not resolve to the selected source');
    requirements.push({
      id, kind, text, selection: 'explicit',
      dependsOn: requirements.length ? [requirements.at(-1).id] : [],
      elements: [{ skill, file, startLine, endLine, digest: e.digest }],
    });
  }
  const plan = workflow('plan', { prompt, requirements }, env);
  assert.deepEqual(plan.gaps ?? [], []);
  assert.deepEqual(plan.stages.map(s => s.requirement.id), specs.map(s => s[0]));
  for (const stage of plan.stages) {
    assert.deepEqual(stage.gaps ?? [], []);
    assert.equal(stage.contributions.length, 1);
    assert.equal(stage.contributions[0].requested, true);
    assert.equal(stage.contributions[0].elements.length, 1);
  }
  return plan;
}

async function loadStage(plan, stageID, env) {
  const stage = plan.stages.find(s => s.requirement.id === stageID);
  const contribution = stage.contributions[0];
  const output = workflow('load', { plan, stage: stageID }, env).contributions;
  assert.equal(output.length, 1);
  const loaded = output[0];
  assert.equal(loaded.id, contribution.id);
  const bundle = path.join(repo, 'skills', loaded.skill);
  assert.equal(loaded.body, await readFile(path.join(bundle, 'SKILL.md'), 'utf8'));
  const readFiles = [];
  for (const asset of loaded.assets ?? []) {
    const raw = await readFile(path.join(bundle, asset.file));
    assert.equal(digest(raw), asset.digest);
    if (asset.path) {
      assert.deepEqual(await readFile(asset.path), raw);
      readFiles.push(asset.file);
    } else {
      assert.equal(asset.body, raw.toString('utf8'));
    }
  }
  const element = contribution.elements[0];
  const record = (status, evidence) => ({
    status, evidence: [evidence], receipt: loaded.receipt, readFiles,
    contribution: contribution.id, elements: [element.file + '#L' + element.startLine + '-L' + element.endLine],
  });
  return { loaded, bundle, record };
}

function dashboard(palette) {
  const fixed = Boolean(palette);
  const style = fixed ? `
    body { background: ${palette.neutral[0]}; color: ${palette.neutral[1]}; font: 1rem/1.5 system-ui; }
    h1 { color: ${palette.primary[0]}; }
    button { background: ${palette.primary[0]}; color: ${palette.neutral[0]}; }
    :focus-visible { outline: 3px solid ${palette.primary[1]}; outline-offset: 2px; }` : '';
  const action = fixed ? '<button type="button" id="details">Open ticket H-104</button>'
    : '<div id="details">Open ticket H-104</div>';
  return `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Harbor queue</title>
    <style>${style}</style></head><body><h1>Harbor service desk</h1>
    <main><h2>Open tickets</h2><label for="filter">Filter tickets</label><input id="filter" type="search">
    <ul><li data-ticket="H-087">H-087 — Billing question</li><li data-ticket="H-104">H-104 — Account update ${action}</li></ul>
    <p id="result" role="status"></p></main><script>
    document.querySelector('#filter').addEventListener('input', event => {
      for (const row of document.querySelectorAll('[data-ticket]')) row.hidden = !row.dataset.ticket.includes(event.target.value);
    });
    document.querySelector('#details').addEventListener('click', () => {
      document.querySelector('#result').textContent = 'H-104: waiting for reply';
    });
    </script></body></html>`;
}

function openBrowser(profile) {
  const child = spawn(chrome, ['--headless=new', '--disable-gpu', '--no-sandbox', '--disable-background-networking',
    '--no-first-run', '--no-default-browser-check', '--remote-debugging-pipe', '--user-data-dir=' + profile, 'about:blank'],
  { stdio: ['ignore', 'ignore', 'pipe', 'pipe', 'pipe'] });
  let stderr = '', pending = '', serial = 0;
  const waiting = new Map();
  const decoder = new StringDecoder('utf8');
  child.stderr.on('data', chunk => { stderr = (stderr + chunk).slice(-2048); });
  function abort(error) {
    for (const request of waiting.values()) request.reject(error);
    waiting.clear();
  }
  child.on('error', abort);
  child.stdio[3].on('error', abort);
  const exited = new Promise(resolve => child.once('close', (code, signal) => {
    abort(new Error('Chrome exited (' + (signal ?? code) + '): ' + stderr));
    resolve();
  }));
  child.stdio[4].on('data', chunk => {
    pending += decoder.write(chunk);
    for (let stop; (stop = pending.indexOf('\0')) >= 0;) {
      const message = JSON.parse(pending.slice(0, stop));
      pending = pending.slice(stop + 1);
      const request = waiting.get(message.id);
      if (!request) continue; // CDP events do not carry request IDs.
      waiting.delete(message.id);
      if (message.error) request.reject(new Error(message.error.message));
      else request.resolve(message.result);
    }
  });
  function send(method, params = {}, sessionId) {
    const id = ++serial;
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        waiting.delete(id);
        reject(new Error('Chrome timed out on ' + method + ': ' + stderr));
      }, 15000);
      waiting.set(id, {
        resolve: value => { clearTimeout(timer); resolve(value); },
        reject: error => { clearTimeout(timer); reject(error); },
      });
      child.stdio[3].write(JSON.stringify({ id, method, params, sessionId }) + '\0');
    });
  }
  return { send, close: async () => { child.kill(); await exited; } };
}

async function measure(browser, file) {
  const { targetId } = await browser.send('Target.createTarget', { url: 'about:blank' });
  try {
    const { sessionId } = await browser.send('Target.attachToTarget', { targetId, flatten: true });
    const call = (method, params) => browser.send(method, params, sessionId);
    const evaluate = async expression => {
      const { result, exceptionDetails } = await call('Runtime.evaluate', { expression, returnByValue: true, awaitPromise: true });
      assert.ok(!exceptionDetails, exceptionDetails?.text);
      return result.value;
    };
    await call('Page.enable');
    const url = pathToFileURL(file).href;
    await call('Page.navigate', { url });
    const deadline = Date.now() + 15000;
    while (!(await evaluate('location.href === ' + JSON.stringify(url) + ' && document.readyState === "complete"'))) {
      assert.ok(Date.now() < deadline, 'Dashboard did not load');
      await new Promise(resolve => setTimeout(resolve, 20));
    }
    await call('Page.bringToFront');
    async function key(code, value, modifiers = 0) {
      // Enter carries a carriage return in Chrome's native keyDown event.
      // A rawKeyDown alone traverses focus but does not activate a button.
      const text = code === 'Enter' ? '\r' : '';
      await call('Input.dispatchKeyEvent', {
        type: text ? 'keyDown' : 'rawKeyDown', key: code, code, windowsVirtualKeyCode: value,
        modifiers, text, unmodifiedText: text,
      });
      await call('Input.dispatchKeyEvent', { type: 'keyUp', key: code, code, windowsVirtualKeyCode: value, modifiers });
    }
    await key('Tab', 9);
    const first = await evaluate('document.activeElement.id');
    await call('Input.insertText', { text: 'H-104' });
    const tickets = await evaluate('Array.from(document.querySelectorAll("[data-ticket]:not([hidden])"), row => row.dataset.ticket)');
    await key('Tab', 9);
    const second = await evaluate('document.activeElement.id');
    const focus = await evaluate(`(() => { const s = getComputedStyle(document.activeElement);
      return { style: s.outlineStyle, width: s.outlineWidth, color: s.outlineColor,
        focusVisible: document.activeElement.matches(':focus-visible') }; })()`);
    await key('Enter', 13);
    const activation = await evaluate('document.querySelector("#result").textContent');
    await key('Tab', 9, 8); // Chromium's shift modifier
    const reverse = await evaluate('document.activeElement.id');
    return { forward: [first, second], reverse, tickets, focus, activation };
  } finally {
    await browser.send('Target.closeTarget', { targetId });
  }
}

const root = await mkdtemp(path.join(tmpdir(), 'skill-router-dashboard-'));
let browser;
try {
  const env = { ...process.env, SKILL_ROUTER_REPO_DIR: repo, SKILL_ROUTER_SKILLS_SOURCE_DIR: path.join(repo, 'skills'),
    SKILL_ROUTER_MANIFEST: path.join(repo, 'manifest.json'), SKILL_ROUTER_TELEMETRY: '0',
    SKILL_ROUTER_CONFIG_DIR: process.env.SKILL_ROUTER_CONFIG_DIR || path.join(root, 'cache') };
  const task = path.join(root, 'service-desk');
  await mkdir(path.join(task, 'docs'), { recursive: true });
  const guidelines = `# Harbor service desk
### Primary Colors
| Role | Hex |
| --- | --- |
| Harbor teal | #124C48 |
| Focus amber | #F6AC19 |
### Neutral
| Role | Hex |
| --- | --- |
| Paper | #FAF8F2 |
| Text | #172C28 |
`;
  await writeFile(path.join(task, 'docs', 'brand-guidelines.md'), guidelines);
  const plan = await makePlan(env);
  const brand = await loadStage(plan, 'brand', env);
  const helper = 'scripts/inject-brand-context.cjs';
  assert.ok(brand.loaded.assets.some(a => a.file === helper), 'Brand section must deliver its helper');
  const context = JSON.parse(execFileSync(process.execPath, [path.join(brand.bundle, helper), '--json'], { cwd: task, encoding: 'utf8' }));
  assert.deepEqual(context.colors.primary, ['#124C48', '#F6AC19']);
  assert.deepEqual(context.colors.neutral, ['#FAF8F2', '#172C28']);
  await writeFile(path.join(task, 'brand-context.json'), JSON.stringify(context));
  const brandRecord = brand.record('applied', { kind: 'artifact', reference: 'brand-context.json',
    observation: 'The loaded helper extracted the four existing palette values from docs/brand-guidelines.md.' });

  const ui = await loadStage(plan, 'interface', env);
  const keyboard = await loadStage(plan, 'keyboard', env);
  browser = openBrowser(path.join(root, 'chrome-profile'));
  const page = path.join(task, 'dashboard.html');
  await writeFile(page, dashboard(null));
  const before = await measure(browser, page);
  assert.equal(before.forward[0], 'filter');
  assert.notEqual(before.forward[1], 'details', 'Broken clickable div should be absent from Tab order');
  assert.equal(before.activation, '', 'Keyboard unexpectedly activated the broken control');

  const html = dashboard(context.colors);
  assert.ok([...context.colors.primary, ...context.colors.neutral].every(color => html.includes(color)));
  await writeFile(page, html);
  const uiRecord = ui.record('applied', { kind: 'artifact', reference: 'dashboard.html',
    observation: 'The queue now uses the extracted palette, a labeled filter, and a native ticket-details button.' });
  const review = records => workflow('review', { plan, records }, env);
  const pending = review([brandRecord, uiRecord]);
  assert.equal(pending.complete, false);
  assert.equal(pending.stages.find(s => s.id === 'keyboard').complete, false);

  const after = await measure(browser, page);
  assert.deepEqual(after.forward, ['filter', 'details']);
  assert.deepEqual(after.tickets, ['H-104']);
  assert.equal(after.activation, 'H-104: waiting for reply');
  assert.equal(after.reverse, 'filter');
  assert.equal(after.focus.focusVisible, true);
  assert.equal(after.focus.style, 'solid');
  assert.ok(Number.parseFloat(after.focus.width) >= 2);
  assert.equal(after.focus.color, 'rgb(246, 172, 25)');
  await writeFile(path.join(task, 'keyboard-observations.json'), JSON.stringify({ before, after }, null, 2));
  const testRecord = keyboard.record('verified', { kind: 'test', outcome: 'passed', reference: 'keyboard-observations.json',
    observation: JSON.stringify(after) });
  const verified = review([brandRecord, uiRecord, testRecord]);
  assert.equal(verified.complete, true, JSON.stringify(verified.gaps));
  console.log(JSON.stringify({ before, after, completeBeforeObservation: pending.complete, completeAfterObservation: verified.complete }));
} finally {
  if (browser) await browser.close();
  await rm(root, { recursive: true, force: true });
}
