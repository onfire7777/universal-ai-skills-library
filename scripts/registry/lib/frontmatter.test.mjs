import test from "node:test";
import assert from "node:assert/strict";
import { parseFrontMatter, collapseWhitespace, readScalarField } from "./frontmatter.mjs";

test("parses simple inline name + description", () => {
  const md = "---\nname: alpha\ndescription: Do a thing.\n---\n# Alpha\n";
  assert.deepEqual(parseFrontMatter(md), { name: "alpha", description: "Do a thing." });
});

test("strips surrounding quotes", () => {
  const md = `---\nname: "beta"\ndescription: 'Quoted desc'\n---\n`;
  assert.deepEqual(parseFrontMatter(md), { name: "beta", description: "Quoted desc" });
});

test("folds block scalar (|) into a single whitespace-collapsed line", () => {
  const md = "---\nname: cso\ndescription: |\n  Chief Security Officer mode.\n  OWASP and STRIDE security audits.\n---\n";
  assert.equal(
    parseFrontMatter(md).description,
    "Chief Security Officer mode. OWASP and STRIDE security audits."
  );
});

test("folds folded scalar (>) and stops at the next top-level key", () => {
  const md = "---\ndescription: >\n  line one\n  line two\nname: gamma\n---\n";
  const fm = parseFrontMatter(md);
  assert.equal(fm.description, "line one line two");
  assert.equal(fm.name, "gamma");
});

test("ignores indented (nested) keys", () => {
  const md = "---\nname: delta\nmeta:\n  description: nested should be ignored\ndescription: real one\n---\n";
  assert.equal(parseFrontMatter(md).description, "real one");
});

test("returns empty object when there is no front matter", () => {
  assert.deepEqual(parseFrontMatter("# no front matter\n"), {});
});

test("tolerates CRLF line endings and a BOM", () => {
  const md = "﻿---\r\nname: eps\r\ndescription: windows\r\n---\r\n";
  assert.deepEqual(parseFrontMatter(md), { name: "eps", description: "windows" });
});

test("collapseWhitespace + readScalarField helpers", () => {
  assert.equal(collapseWhitespace("a   b\n c "), "a b c");
  assert.equal(readScalarField("name: zeta\ndescription: hi", "name"), "zeta");
  assert.equal(readScalarField("name: zeta", "missing"), undefined);
});
