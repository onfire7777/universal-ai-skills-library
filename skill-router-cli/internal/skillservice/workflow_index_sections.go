package skillservice

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

func workflowFrontmatterEnd(lines []string) int {
	if len(lines) == 0 || strings.TrimPrefix(strings.TrimSpace(lines[0]), "\uFEFF") != "---" {
		return 0
	}
	for i, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			return i + 2
		}
	}
	return 0
}

// Read just the fields relevant to planning. Imported frontmatter may contain
// unrelated keys with duplicate names or values in host-specific formats.
func readWorkflowFrontmatter(data []byte, b *WorkflowBundle) {
	lines := strings.Split(string(data), "\n")
	end := workflowFrontmatterEnd(lines)
	if end == 0 {
		return
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:end-1], "\n")), &doc); err != nil || len(doc.Content) == 0 {
		return
	}
	n := doc.Content[0]
	if n.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		v := n.Content[i+1]
		switch n.Content[i].Value {
		case "description":
			if v.Kind == yaml.ScalarNode && v.Tag == "!!str" {
				b.SourceDescription = strings.Join(strings.Fields(v.Value), " ")
			}
		case "version":
			if v.Kind == yaml.ScalarNode {
				b.Version = v.Value
			}
		case "capabilities":
			b.Capabilities = yamlStrings(v)
		case "conflicts":
			b.Conflicts = yamlStrings(v)
		case "supersedes":
			b.Supersedes = yamlStrings(v)
		case "requires":
			if v.Kind != yaml.MappingNode {
				b.Requires = yamlStrings(v)
				continue
			}
			for j := 0; j+1 < len(v.Content); j += 2 {
				key := v.Content[j].Value
				for _, item := range yamlStrings(v.Content[j+1]) {
					b.Requires = append(b.Requires, key+":"+item)
				}
			}
		}
	}
}

func yamlStrings(n *yaml.Node) []string {
	switch n.Kind {
	case yaml.ScalarNode:
		if n.Tag == "!!str" && strings.TrimSpace(n.Value) != "" {
			return []string{n.Value}
		}
	case yaml.SequenceNode:
		var out []string
		for _, c := range n.Content {
			out = append(out, yamlStrings(c)...)
		}
		return out
	}
	return nil
}

var (
	workflowHeading     = regexp.MustCompile(`^ {0,3}#{1,6}[[:space:]]+(.*)$`)
	workflowSetext      = regexp.MustCompile(`^ {0,3}(=+|-+)[[:space:]]*$`)
	workflowSymbol      = regexp.MustCompile(`^[[:space:]]*(?:(?:export|default|async|public|private|protected|static)[[:space:]]+)*(?:(?:def|class|function|interface|type)[[:space:]]+([\pL_][\pL\pN_$]*)|func[[:space:]]+(?:\([^)]*\)[[:space:]]*)?([\pL_][\pL\pN_]*)|([\pL_][\pL\pN_-]*)[[:space:]]*\(\)[[:space:]]*\{)`)
	workflowLinks       = regexp.MustCompile(`\]\(\s*(?:<([^>]+)>|([^\s)]+))`)
	workflowPaths       = regexp.MustCompile(`[\pL\pN_./-]+\.[\pL\pN]{1,12}`)
	workflowHTMLHeading = regexp.MustCompile(`(?i)^\s*<h([1-6])(?:\s[^>]*|)>(.*?)</h[1-6]>`)
	workflowHTMLTags    = regexp.MustCompile(`<[^>]*>`)
	workflowXMLSymbol   = regexp.MustCompile(`(?i)^\s*<(?:(?:xs|xsd):)(?:complexType|simpleType|element|group)\s+name=["']([^"']+)["']`)
)

type workflowSection struct {
	first         int
	title, origin string
	context       []string
	summary       string
}

func extractWorkflowFile(b *WorkflowBundle, f indexedBundleFile, known map[string]bool) ([]WorkflowElement, []string, []string) {
	role := workflowFileRole(f.path)
	base := WorkflowElement{Skill: b.Skill, File: f.path, Digest: f.digest, Role: role, AppliesWhen: b.Description, Requires: b.Requires}
	base.Origins = WorkflowOrigins{Role: "inferred from path", AppliesWhen: "bundle description"}
	if len(base.Requires) > 0 {
		base.Origins.Requires = "bundle frontmatter"
	}
	if !utf8.Valid(f.content) || bytes.IndexByte(f.content, 0) >= 0 {
		base.Binary = true
		base.SectionDigest = f.digest
		base.Title = path.Base(f.path)
		base.Purpose = base.Title
		base.Origins.Title = "filename"
		base.Origins.Purpose = "inferred from filename"
		return []WorkflowElement{base}, []string{workflowVocabulary(f.path)}, []string{""}
	}
	lines := strings.Split(string(f.content), "\n")
	if len(lines) > 1 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	sections := splitWorkflowSections(f.path, lines)
	el := make([]WorkflowElement, 0, len(sections))
	terms := make([]string, 0, len(sections))
	bodies := make([]string, 0, len(sections))
	for i, sec := range sections {
		stop := len(lines)
		if i+1 < len(sections) {
			stop = sections[i+1].first
		}
		e := base
		e.Title = sec.title
		e.Context = sec.context
		e.Origins.Title = sec.origin
		e.StartLine = sec.first + 1
		e.EndLine = stop
		chunk := lines[sec.first:stop]
		body := strings.Join(chunk, "\n")
		digest := sha256.Sum256([]byte(strings.TrimSpace(strings.ReplaceAll(body, "\r\n", "\n"))))
		e.SectionDigest = hex.EncodeToString(digest[:])
		e.Purpose, e.Origins.Purpose = workflowPurpose(chunk, sec.title, sec.origin)
		if sec.summary != "" {
			e.Purpose = sec.summary
			e.Origins.Purpose = "JSON record description"
		}
		roleHint := strings.ToLower(sec.title)
		if sec.origin == "heading" || sec.origin == "setext heading" {
			switch {
			case strings.Contains(roleHint, "when to use") || strings.Contains(roleHint, "when to apply") || strings.Contains(roleHint, "applicability"):
				if s := workflowDetails(chunk); len(s) > 0 {
					e.AppliesWhen = strings.Join(s, " ")
					e.Origins.AppliesWhen = "section text"
				}
			case roleHint == "inputs" || roleHint == "prerequisites" || roleHint == "requirements":
				if s := workflowDetails(chunk); len(s) > 0 {
					e.Requires = append(slices.Clone(b.Requires), s...)
					e.Origins.Requires = "section text and bundle declarations"
				}
			case roleHint == "output" || roleHint == "outputs" || roleHint == "deliverables" || roleHint == "expected output":
				e.Produces = workflowDetails(chunk)
				if len(e.Produces) > 0 {
					e.Origins.Produces = "section text"
				}
			}
		}
		if role == "instruction" || role == "reference" || role == "workflow" {
			switch {
			case hasWorkflowWord(roleHint, "verify", "verification", "checks", "checklist", "testing", "audit", "validation"):
				e.Role = "check"
				e.Origins.Role = "inferred from heading"
			case hasWorkflowWord(roleHint, "workflow", "steps", "procedure", "method"):
				e.Role = "workflow"
				e.Origins.Role = "inferred from heading"
			case hasWorkflowWord(roleHint, "example", "examples"):
				e.Role = "example"
				e.Origins.Role = "inferred from heading"
			}
		}
		if sec.origin == "heading" || sec.origin == "setext heading" {
			workflowLabeledDetails(chunk, &e)
		}
		e.References = workflowReferences(f.path, body, known)
		el = append(el, e)
		terms = append(terms, workflowVocabulary(body))
		bodies = append(bodies, body)
	}
	return el, terms, bodies
}

func workflowFileRole(file string) string {
	if file == "SKILL.md" {
		return "instruction"
	}
	p := strings.ToLower(file)
	base := path.Base(p)
	if base == ".openskills.json" || base == "license" || strings.HasPrefix(base, "license.") || strings.HasPrefix(base, "package-lock.") {
		return "metadata"
	}
	for _, part := range strings.Split(p, "/") {
		switch part {
		case "templates", "template":
			return "template"
		case "examples", "example", "samples":
			return "example"
		case "scripts", "bin":
			return "script"
		case "workflows", "workflow":
			return "workflow"
		case "assets", "canvas-fonts", "fonts":
			return "asset"
		case "references", "reference", "rules", "resources", "docs":
			return "reference"
		}
	}
	if workflowCode(p) {
		return "script"
	}
	if path.Ext(p) == ".md" || path.Ext(p) == ".txt" {
		return "reference"
	}
	return "asset"
}

func workflowCode(p string) bool {
	switch strings.ToLower(path.Ext(p)) {
	case ".py", ".go", ".js", ".mjs", ".cjs", ".ts", ".jsx", ".tsx", ".sh", ".bash", ".zsh", ".ps1", ".rb", ".rs":
		return true
	}
	return false
}

func splitWorkflowSections(file string, lines []string) []workflowSection {
	first := 0
	ext := strings.ToLower(path.Ext(file))
	if ext == ".json" && len(lines) > 1 && workflowFileRole(file) != "metadata" {
		if sections := workflowJSONSections(lines); len(sections) > 0 {
			if sections[0].first > 0 && strings.Trim(strings.Join(lines[:sections[0].first], " "), " \t\r\n{[") != "" {
				sections = append([]workflowSection{{first: 0, title: path.Base(file), origin: "filename"}}, sections...)
			}
			return sections
		}
	}
	markdown := ext == ".md" || ext == ".markdown" || ext == ".txt"
	if markdown {
		first = workflowFrontmatterEnd(lines)
	}
	// Metadata-only files still have a source location.
	if first >= len(lines) {
		return []workflowSection{{first: 0, title: path.Base(file), origin: "filename"}}
	}
	sections := []workflowSection{}
	var parents [6]string
	add := func(n int, title, origin string, level int) {
		if len(sections) > 0 && sections[len(sections)-1].first == n {
			return
		}
		sec := workflowSection{first: n, title: title, origin: origin}
		if level > 0 {
			for _, p := range parents[:level-1] {
				if p != "" {
					sec.context = append(sec.context, p)
				}
			}
			clear(parents[level-1:])
			parents[level-1] = title
		}
		sections = append(sections, sec)
	}
	fence := byte(0)
	fenceSize := 0
	for i := first; i < len(lines); i++ {
		line := strings.TrimSuffix(lines[i], "\r")
		trim := strings.TrimSpace(line)
		if markdown {
			if len(trim) >= 3 && (trim[0] == 0x60 || trim[0] == '~') && len(line)-len(strings.TrimLeft(line, " ")) <= 3 {
				n := 0
				for n < len(trim) && trim[n] == trim[0] {
					n++
				}
				if n >= 3 && fence == 0 {
					fence = trim[0]
					fenceSize = n
					continue
				}
				if fence == trim[0] && n >= fenceSize && strings.TrimSpace(trim[n:]) == "" {
					fence = 0
					continue
				}
			}
			if fence != 0 {
				continue
			}
			if match := workflowHeading.FindStringSubmatch(line); len(match) > 1 {
				title := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(match[1]), "#"))
				if title != "" {
					add(i, title, "heading", len(strings.TrimLeft(line, " "))-len(strings.TrimLeft(strings.TrimLeft(line, " "), "#")))
				}
			} else if i > first && workflowSetext.MatchString(line) {
				prev := strings.TrimSpace(lines[i-1])
				if prev != "" && !strings.HasPrefix(prev, "#") && !strings.ContainsAny(prev, "|`") && !strings.HasPrefix(prev, "-") {
					level := 2
					if strings.HasPrefix(trim, "=") {
						level = 1
					}
					add(i-1, prev, "setext heading", level)
				}
			}
		} else if ext == ".html" || ext == ".htm" {
			if m := workflowHTMLHeading.FindStringSubmatch(line); len(m) > 2 {
				title := strings.TrimSpace(html.UnescapeString(workflowHTMLTags.ReplaceAllString(m[2], "")))
				if title != "" {
					add(i, title, "HTML heading", int(m[1][0]-'0'))
				}
			}
		} else if ext == ".xsd" {
			if m := workflowXMLSymbol.FindStringSubmatch(line); len(m) > 1 {
				add(i, m[1], "schema declaration", 0)
			}
		} else if workflowCode(file) {
			match := workflowSymbol.FindStringSubmatch(line)
			if len(match) > 1 {
				for _, m := range match[1:] {
					if m != "" {
						add(i, m, "source symbol", 0)
						break
					}
				}
			}
		}
	}
	if len(sections) == 0 {
		return []workflowSection{{first: first, title: path.Base(file), origin: "filename"}}
	}
	if sections[0].first > first && strings.TrimSpace(strings.Join(lines[first:sections[0].first], " ")) != "" {
		return append([]workflowSection{{first: first, title: path.Base(file), origin: "filename"}}, sections...)
	}
	return sections
}

func hasWorkflowWord(s string, wants ...string) bool {
	for _, word := range strings.FieldsFunc(s, workflowSeparator) {
		if slices.Contains(wants, word) {
			return true
		}
	}
	return false
}

func workflowPurpose(lines []string, title, origin string) (string, string) {
	begin := 0
	if origin == "heading" || origin == "source symbol" {
		begin = 1
	} else if origin == "setext heading" {
		begin = 2
	}
	for _, line := range lines[min(begin, len(lines)):] {
		t := strings.TrimSpace(line)
		if t == "" || t == "---" {
			continue
		}
		// A table, code block, or symbol body is better identified by its
		// authored title than by a truncated first line of source code.
		if origin == "source symbol" || strings.HasPrefix(t, "|") || strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			break
		}
		return workflowExcerpt(t, 240), "section text"
	}
	return workflowExcerpt(title, 240), "inferred from " + origin
}

func workflowDetails(lines []string) []string {
	var out []string
	for _, line := range lines[1:] {
		t := strings.TrimSpace(line)
		if t != "" && t != "---" && !workflowSetext.MatchString(t) {
			out = append(out, strings.TrimSpace(strings.TrimPrefix(t, "- ")))
		}
	}
	return out
}

func workflowExcerpt(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return fmt.Sprintf("%s…", string([]rune(s)[:n]))
}

func workflowReferences(file, body string, known map[string]bool) []string {
	set := map[string]bool{}
	add := func(target string) {
		if strings.Contains(target, "://") || strings.HasPrefix(target, "/") {
			return
		}
		target, _, _ = strings.Cut(target, "#")
		target, _, _ = strings.Cut(target, "?")
		if decoded, err := url.PathUnescape(target); err == nil {
			target = decoded
		}
		for _, p := range []string{path.Join(path.Dir(file), target), path.Clean(target)} {
			if p != file && known[p] {
				set[p] = true
				break
			}
		}
	}
	for _, link := range workflowLinks.FindAllStringSubmatch(body, -1) {
		add(link[1] + link[2])
	}
	for _, word := range workflowPaths.FindAllString(body, -1) {
		add(word)
	}
	if len(set) == 0 {
		return nil
	}
	r := make([]string, 0, len(set))
	for p := range set {
		r = append(r, p)
	}
	slices.Sort(r)
	return r
}

// A procedure can declare its own preconditions and outputs inline. Only
// promote explicit labels outside code fences to metadata.
func workflowLabeledDetails(lines []string, e *WorkflowElement) {
	var fence string
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			if fence == "" {
				fence = t[:3]
			} else if strings.HasPrefix(t, fence) {
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}
		key, value, ok := strings.Cut(t, ":")
		if !ok {
			continue
		}
		value = strings.TrimSpace(strings.TrimLeft(value, "* "))
		if value == "" {
			continue
		}
		switch strings.ToLower(strings.Trim(key, " -*")) {
		case "applies when", "when to use", "use when":
			e.AppliesWhen = value
			e.Origins.AppliesWhen = "section text"
		case "input", "inputs", "requires", "prerequisite", "prerequisites":
			e.Requires = append(slices.Clone(e.Requires), value)
			e.Origins.Requires = "section text and bundle declarations"
		case "output", "outputs", "produces", "expected result", "deliverable", "deliverables":
			e.Produces = append(e.Produces, value)
			e.Origins.Produces = "section text"
		}
	}
}

// Record boundaries and object keys preserve the scope of structured reference
// data. Compact JSON records sharing a physical line remain one retrievable
// unit, because the public locator is a line range.
func workflowJSONSections(lines []string) []workflowSection {
	text := strings.Join(lines, "\n")
	decoder := json.NewDecoder(strings.NewReader(text))
	begin, err := decoder.Token()
	if err != nil {
		return nil
	}
	delim, ok := begin.(json.Delim)
	if !ok || (delim != '[' && delim != '{') {
		return nil
	}
	var sections []workflowSection
	line, offset := 0, 0
	for n := 1; decoder.More(); n++ {
		from := int(decoder.InputOffset())
		for from < len(text) && strings.ContainsRune(" \t\r\n,", rune(text[from])) {
			from++
		}
		title := ""
		summary := ""
		if delim == '{' {
			key, e := decoder.Token()
			if e != nil {
				return nil
			}
			title, ok = key.(string)
			if !ok {
				return nil
			}
		}
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return nil
		}
		if len(raw) == 0 {
			continue
		}
		if delim == '[' && raw[0] != '{' && raw[0] != '[' {
			continue
		}
		var record struct {
			Name        string
			Title       string
			ID          string
			Description string
		}
		if raw[0] == '{' {
			_ = json.Unmarshal(raw, &record)
			summary = workflowExcerpt(record.Description, 240)
		}
		if delim == '[' {
			switch {
			case record.Title != "":
				title = record.Title
			case record.Name != "":
				title = record.Name
			case record.ID != "":
				title = record.ID
			default:
				title = fmt.Sprintf("Record %d", n)
			}
		}
		line += strings.Count(text[offset:from], "\n")
		offset = from
		if len(sections) == 0 || sections[len(sections)-1].first != line {
			sections = append(sections, workflowSection{first: line, title: workflowExcerpt(title, 160), origin: "JSON record", summary: summary})
		}
	}
	if _, err := decoder.Token(); err != nil {
		return nil
	}
	return sections
}
