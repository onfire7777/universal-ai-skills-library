package cli

import "testing"

func TestEnsureLocalReadAllowedAllowsUnfilteredCollection(t *testing.T) {
	if err := ensureLocalReadAllowed("projects", true, "/projects", nil); err != nil {
		t.Fatalf("expected broad collection read to be allowed: %v", err)
	}
}

func TestEnsureLocalReadAllowedRejectsFilteredCollection(t *testing.T) {
	err := ensureLocalReadAllowed("projects", true, "/projects", map[string]string{"status": "active"})
	if err == nil {
		t.Fatal("expected filtered local read to be rejected")
	}
}

func TestEnsureLocalReadAllowedRejectsScopedNestedCollection(t *testing.T) {
	err := ensureLocalReadAllowed("tasks", true, "/projects/team-a/tasks", nil)
	if err == nil {
		t.Fatal("expected scoped nested local read to be rejected")
	}
}

func TestEnsureLocalReadAllowedRejectsSingleItemRead(t *testing.T) {
	err := ensureLocalReadAllowed("projects", false, "/projects/team-a", nil)
	if err == nil {
		t.Fatal("expected single-item local read to be rejected")
	}
}
