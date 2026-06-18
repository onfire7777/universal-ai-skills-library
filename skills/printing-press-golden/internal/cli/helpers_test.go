package cli

import "testing"

func TestReplacePathParamEscapesPathSegment(t *testing.T) {
	got := replacePathParam("/projects/{projectId}/tasks", "projectId", "team/a?debug=true#frag")
	want := "/projects/team%2Fa%3Fdebug=true%23frag/tasks"
	if got != want {
		t.Fatalf("replacePathParam() = %q, want %q", got, want)
	}
}
