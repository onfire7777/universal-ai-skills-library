package skillservice

import (
	"strings"
	"testing"
)

func TestExternalSkillRootsExcludeRetiredMarketplaceAndLegacyRoots(t *testing.T) {
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", "")

	for _, root := range externalSkillRoots() {
		if root.ID == "claude-market" || strings.Contains(root.Path, "plugins/marketplaces") {
			t.Fatalf("retired Claude marketplace root is still auto-discovered: %#v", root)
		}
		if root.ID == "legacy-compat" || strings.Contains(root.Path, ".manus") {
			t.Fatalf("legacy Manus root is still auto-discovered: %#v", root)
		}
	}
}

func TestExternalSkillRootsKeepExplicitOptInRoots(t *testing.T) {
	extraRoot := t.TempDir()
	t.Setenv("SKILL_ROUTER_EXTERNAL_SKILL_ROOTS", extraRoot)

	found := false
	for _, root := range externalSkillRoots() {
		if root.ID == "extra-1" && root.Path == extraRoot {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("explicit external skill root was not preserved")
	}
}
