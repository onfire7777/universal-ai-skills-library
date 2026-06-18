package api

import (
	"encoding/json"
	"testing"
)

func TestPathSegmentEscapesIDs(t *testing.T) {
	got := pathSegment(`task/../projects?x=1#frag`)
	want := "task%2F..%2Fprojects%3Fx=1%23frag"
	if got != want {
		t.Fatalf("pathSegment() = %q, want %q", got, want)
	}
}

func TestJSONBodyEscapesUserStrings(t *testing.T) {
	body, err := jsonBody(map[string]any{"description": `x","priority":"urgent`})
	if err != nil {
		t.Fatalf("jsonBody returned error: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("jsonBody produced invalid JSON: %s", body)
	}
	if decoded["description"] != `x","priority":"urgent` {
		t.Fatalf("description = %q", decoded["description"])
	}
	if _, ok := decoded["priority"]; ok {
		t.Fatalf("injected priority field decoded from %s", body)
	}
}

func TestSplitEventsTrimsEmptyParts(t *testing.T) {
	got := splitEvents("task.completed, task.failed,,")
	want := []string{"task.completed", "task.failed"}
	if len(got) != len(want) {
		t.Fatalf("splitEvents() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("splitEvents()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
