package registry

import (
	"sort"
)

// normalize.go ports the Node generator's drift comparison: --check compares by
// MEANING, not bytes, so it is invariant to the faithful-vs-optimize formatting
// choice and reports only real drift (skills/descriptions/aliases/script-sets).

func sortAnyStrings(arr []any) {
	sort.SliceStable(arr, func(i, j int) bool {
		a, _ := arr[i].(string)
		b, _ := arr[j].(string)
		return a < b
	})
}

// canonicalEntry mirrors canonicalEntry() in generate-registry.mjs.
func canonicalEntry(e *OM) *OM {
	scripts := []any{}
	if s, ok := e.Get("scripts"); ok {
		scripts = append(scripts, asArr(s)...)
	}
	sortAnyStrings(scripts)
	c := NewOM()
	c.Set("name", get(e, "name"))
	c.Set("directory", get(e, "directory"))
	c.Set("description", get(e, "description"))
	aliases := []any{}
	if a, ok := e.Get("aliases"); ok {
		if arr := asArr(a); len(arr) > 0 {
			aliases = arr
		}
	}
	c.Set("aliases", aliases)
	c.Set("has_scripts", len(scripts) > 0)
	c.Set("scripts", scripts)
	return c
}

func mapCanonical(obj *OM, key string) []any {
	src := asArr(get(obj, key))
	out := make([]any, 0, len(src))
	for _, e := range src {
		out = append(out, canonicalEntry(asOM(e)))
	}
	return out
}

// normalizeForCompare mirrors normalizeForCompare() in generate-registry.mjs.
func normalizeForCompare(key, text string) string {
	v, err := Parse([]byte(text))
	if err != nil {
		return text
	}
	obj, ok := v.(*OM)
	if !ok {
		return text
	}
	switch key {
	case "manifest":
		out := NewOM()
		out.Set("version", get(obj, "version"))
		out.Set("description", get(obj, "description"))
		out.Set("canonical_id_policy", get(obj, "canonical_id_policy"))
		out.Set("routing", get(obj, "routing"))
		out.Set("total_skills", get(obj, "total_skills"))
		out.Set("core_skills", mapCanonical(obj, "core_skills"))
		out.Set("library_skills", mapCanonical(obj, "library_skills"))
		return Stringify(out)
	case "build-manifest":
		obj.Delete("generated_at")
		for _, e := range asArr(get(obj, "skills")) {
			eo := asOM(e)
			if sc, ok := eo.Get("scripts"); ok {
				sorted := append([]any{}, asArr(sc)...)
				sortAnyStrings(sorted)
				eo.Set("scripts", sorted)
			}
		}
		return Stringify(obj)
	}
	return text
}
