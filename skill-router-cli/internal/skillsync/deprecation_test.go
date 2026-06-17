package skillsync

import (
	"strings"
	"testing"
)

func TestDeprecationNoticeMentionsCLI(t *testing.T) {
	msg := DeprecationNotice()
	for _, want := range []string{"deprecated", "skill-router", "serve"} {
		if !strings.Contains(msg, want) {
			t.Errorf("notice missing %q: %s", want, msg)
		}
	}
}

func TestDeprecationNoticePointsToMigrationDoc(t *testing.T) {
	if !strings.Contains(DeprecationNotice(), "ADAPTER_DEPRECATION.md") {
		t.Errorf("notice should reference the migration doc: %s", DeprecationNotice())
	}
}
