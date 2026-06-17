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
