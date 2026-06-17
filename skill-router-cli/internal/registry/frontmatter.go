package registry

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// frontmatter.go is a byte-faithful Go port of scripts/registry/lib/frontmatter.mjs.
// It MUST reproduce that reader's normalization exactly, because ~1,804 of the
// 1,812 skill descriptions in manifest.json come straight from SKILL.md front
// matter — any divergence breaks byte-parity with the Node generator.

// jsWhitespace matches the exact set JavaScript's \s matches (used by the Node
// reader's collapseWhitespace / trim), so collapsing behaves identically.
var jsWhitespace = regexp.MustCompile(`[\x09\x0a\x0b\x0c\x0d \x{00a0}\x{1680}\x{2000}-\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000}\x{feff}]+`)

// blockScalarIndicator matches `|`/`>` block-scalar headers (with optional
// chomping/indent markers and trailing spaces): mirrors /^[|>][+-]?\d*\s*$/.
var blockScalarIndicator = regexp.MustCompile(`^[|>][+-]?[0-9]*[ \t]*$`)

// bom is the UTF-8 byte-order mark (U+FEFF), stripped from the head of SKILL.md.
const bom = string(rune(0xFEFF))

// isJSSpace reports whether r is in JavaScript's \s set (used for trimming).
func isJSSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ': // \x09..\x0d and SPACE
		return true
	case 0x00a0, 0x1680,
		0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006,
		0x2007, 0x2008, 0x2009, 0x200a,
		0x2028, 0x2029, 0x202f, 0x205f, 0x3000, 0xfeff:
		return true
	}
	return false
}

// collapseWhitespace folds runs of whitespace to a single space and trims the
// ends — equivalent to `value.replace(/\s+/g, " ").trim()`.
func collapseWhitespace(value string) string {
	value = jsWhitespace.ReplaceAllString(value, " ")
	return strings.TrimFunc(value, isJSSpace)
}

// stripQuotes removes a single matching pair of surrounding ' or " quotes,
// mirroring /^(['"])([\s\S]*)\1$/.
func stripQuotes(value string) string {
	if len(value) < 2 {
		return value
	}
	first, firstSize := utf8.DecodeRuneInString(value)
	last, lastSize := utf8.DecodeLastRuneInString(value)
	if (first == '\'' || first == '"') && first == last {
		return value[firstSize : len(value)-lastSize]
	}
	return value
}

// startsWithJSSpace reports whether s begins with a JS-whitespace rune.
func startsWithJSSpace(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	return r != utf8.RuneError && isJSSpace(r)
}

// extractFrontMatterBlock returns the inner body of a leading `---\n ... \n---`
// block (and true), or ("", false) when absent. Tolerates a UTF-8 BOM and CRLF.
func extractFrontMatterBlock(content string) (string, bool) {
	content = strings.TrimPrefix(content, bom)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	const open = "---\n"
	if !strings.HasPrefix(content, open) {
		return "", false
	}
	rest := content[len(open):]
	// Find the first "\n---" whose closing fence is terminated by \n or EOF
	// (the non-greedy /([\s\S]*?)\n---(?:\n|$)/).
	search := 0
	for {
		idx := strings.Index(rest[search:], "\n---")
		if idx < 0 {
			return "", false
		}
		idx += search
		after := idx + len("\n---")
		if after == len(rest) || rest[after] == '\n' {
			return rest[:idx], true
		}
		search = idx + 1
	}
}

// readScalarField reads a single top-level scalar field from a front-matter
// body, supporting inline, quoted, and `|`/`>` block-scalar values. Returns the
// whitespace-collapsed string, or "" if the key is absent (which the caller
// treats the same as an empty value, matching the Node generator).
func readScalarField(body, key string) string {
	prefix := key + ":"
	lines := strings.Split(body, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		// Top-level keys start at column 0; indented lines belong to a parent.
		if startsWithJSSpace(line) {
			continue
		}
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		rest := line[len(prefix):]
		// `replace(/^\s/, "")`: drop exactly one leading whitespace rune.
		if startsWithJSSpace(rest) {
			_, size := utf8.DecodeRuneInString(rest)
			rest = rest[size:]
		}
		if blockScalarIndicator.MatchString(rest) {
			collected := make([]string, 0, 8)
			for j := i + 1; j < len(lines); j++ {
				l := lines[j]
				if strings.TrimFunc(l, isJSSpace) == "" {
					collected = append(collected, "")
					continue
				}
				if !startsWithJSSpace(l) { // a non-indented line ends the block
					break
				}
				collected = append(collected, strings.TrimLeftFunc(l, isJSSpace))
			}
			return collapseWhitespace(strings.Join(collected, " "))
		}
		return collapseWhitespace(stripQuotes(strings.TrimFunc(rest, isJSSpace)))
	}
	return ""
}

// parseFrontMatterDescription returns the collapsed `description` field of a
// SKILL.md file, or "" when there is no front matter / no description.
func parseFrontMatterDescription(content string) string {
	body, ok := extractFrontMatterBlock(content)
	if !ok {
		return ""
	}
	return readScalarField(body, "description")
}
