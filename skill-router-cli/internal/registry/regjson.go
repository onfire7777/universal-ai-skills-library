// Package registry is the Go owner of the build step: it scans skills/ + the
// registry config and emits the CLI-first registry artifacts (manifest.json and
// docs/build_manifest.json), reproducing the legacy Node generator
// (scripts/registry/generate-registry.mjs) byte-for-byte so the Node tool can be
// retired behind the --check parity gate.
//
// regjson.go provides an insertion-ordered JSON value model plus a serializer
// that is byte-compatible with Node's `JSON.stringify(obj, null, 2) + "\n"`:
//   - 2-space indentation, trailing newline
//   - keys emitted in INSERTION order (Go maps sort keys; JSON.stringify does not)
//   - no HTML escaping of < > & (Go's encoding/json escapes them by default)
//   - control chars escaped as \uXXXX (lowercase), everything >= 0x20 emitted raw
package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// OM is an insertion-ordered JSON object.
type OM struct {
	keys []string
	vals map[string]any
}

// NewOM returns an empty ordered object.
func NewOM() *OM { return &OM{vals: map[string]any{}} }

// Set appends key k (preserving first-insertion order) and stores v.
func (o *OM) Set(k string, v any) *OM {
	if _, ok := o.vals[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.vals[k] = v
	return o
}

// Get returns the value for k and whether it is present.
func (o *OM) Get(k string) (any, bool) { v, ok := o.vals[k]; return v, ok }

// Keys returns the keys in insertion order.
func (o *OM) Keys() []string { return o.keys }

// Delete removes a key while preserving order of the rest.
func (o *OM) Delete(k string) {
	if _, ok := o.vals[k]; !ok {
		return
	}
	delete(o.vals, k)
	for i, kk := range o.keys {
		if kk == k {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
}

// Clone deep-copies the ordered value tree (objects and arrays are duplicated;
// scalars are shared, which is safe because they are immutable).
func (o *OM) Clone() *OM { return cloneVal(o).(*OM) }

func cloneVal(v any) any {
	switch t := v.(type) {
	case *OM:
		n := NewOM()
		for _, k := range t.keys {
			n.Set(k, cloneVal(t.vals[k]))
		}
		return n
	case []any:
		n := make([]any, len(t))
		for i := range t {
			n[i] = cloneVal(t[i])
		}
		return n
	default:
		return v
	}
}

// Parse decodes JSON bytes into the ordered model: objects -> *OM, arrays ->
// []any, plus string / json.Number-or-float64 / bool / nil. Object key order is
// preserved exactly as written, which is what lets embedded config sub-trees
// (routing, build provenance, ...) round-trip byte-for-byte.
func Parse(data []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return parseValue(dec)
}

func parseValue(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return tok, nil // string, json.Number, bool, or nil
	}
	switch delim {
	case '{':
		obj := NewOM()
		for dec.More() {
			keyTok, err := dec.Token()
			if err != nil {
				return nil, err
			}
			key, _ := keyTok.(string)
			val, err := parseValue(dec)
			if err != nil {
				return nil, err
			}
			obj.Set(key, val)
		}
		if _, err := dec.Token(); err != nil { // consume '}'
			return nil, err
		}
		return obj, nil
	case '[':
		arr := []any{}
		for dec.More() {
			val, err := parseValue(dec)
			if err != nil {
				return nil, err
			}
			arr = append(arr, val)
		}
		if _, err := dec.Token(); err != nil { // consume ']'
			return nil, err
		}
		return arr, nil
	}
	return nil, fmt.Errorf("regjson: unexpected delimiter %q", delim)
}

// Stringify renders v as Node's `JSON.stringify(v, null, 2) + "\n"`.
func Stringify(v any) string {
	var b strings.Builder
	writeValue(&b, v, 0)
	b.WriteByte('\n')
	return b.String()
}

func writeIndent(b *strings.Builder, level int) {
	for i := 0; i < level*2; i++ {
		b.WriteByte(' ')
	}
}

func writeValue(b *strings.Builder, v any, level int) {
	switch t := v.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if t {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case string:
		writeJSString(b, t)
	case json.Number:
		b.WriteString(string(t))
	case int:
		b.WriteString(strconv.Itoa(t))
	case int64:
		b.WriteString(strconv.FormatInt(t, 10))
	case float64:
		b.WriteString(formatNumber(t))
	case []string:
		arr := make([]any, len(t))
		for i := range t {
			arr[i] = t[i]
		}
		writeArray(b, arr, level)
	case []any:
		writeArray(b, t, level)
	case *OM:
		writeObject(b, t, level)
	default:
		panic(fmt.Sprintf("regjson: unsupported type %T", v))
	}
}

func writeArray(b *strings.Builder, arr []any, level int) {
	if len(arr) == 0 {
		b.WriteString("[]")
		return
	}
	b.WriteByte('[')
	for i, e := range arr {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
		writeIndent(b, level+1)
		writeValue(b, e, level+1)
	}
	b.WriteByte('\n')
	writeIndent(b, level)
	b.WriteByte(']')
}

func writeObject(b *strings.Builder, o *OM, level int) {
	if len(o.keys) == 0 {
		b.WriteString("{}")
		return
	}
	b.WriteByte('{')
	for i, k := range o.keys {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
		writeIndent(b, level+1)
		writeJSString(b, k)
		b.WriteString(": ")
		writeValue(b, o.vals[k], level+1)
	}
	b.WriteByte('\n')
	writeIndent(b, level)
	b.WriteByte('}')
}

// formatNumber renders a number the way JSON.stringify would. Our artifacts only
// ever carry integers, so integral values are emitted without a fractional part.
func formatNumber(f float64) string {
	if !math.IsInf(f, 0) && !math.IsNaN(f) && f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// writeJSString escapes exactly like JSON.stringify: it does NOT escape < > & or
// '/', emits raw UTF-8 for code points >= 0x20, and escapes control chars as
// lowercase \uXXXX (with the named short escapes for \b \f \n \r \t).
func writeJSString(b *strings.Builder, s string) {
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(b, `\u%04x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
}
