package skillservice

import (
	"slices"
	"strings"
	"testing"
)

func TestNaturalClausesPreserveObjectsAndSplitDistinctActions(t *testing.T) {
	input := "Archive and rename the PDF files, export their metadata. Verify the totals in v2.4.1. Keep the name and date together."
	want := []string{"Archive and rename the PDF files", "export their metadata", "Verify the totals in v2.4.1", "Keep the name and date together"}
	if got := workflowClauses(input); !slices.Equal(got, want) {
		t.Fatalf("%q", got)
	}
}
func TestNaturalFollowupCarriesScopedContext(t *testing.T) {
	rs := naturalRequirements("Create a PDF booklet. Inspect every page for clipping and create a calendar invitation.")
	if len(rs) != 3 || !slices.Equal(rs[1].DependsOn, []string{rs[0].ID}) || !strings.Contains(rs[1].Text, "PDF") || strings.Contains(rs[2].Text, "PDF") || len(rs[2].DependsOn) != 0 {
		t.Fatal(rs)
	}
}

func TestNaturalImplicitFollowupResolvesPlatform(t *testing.T) {
	rs := naturalRequirements("Draw a Mermaid timeline for the project. Show the dates as annotations. Create a PDF invoice.")
	if len(rs) != 3 || !strings.Contains(rs[1].Text, "Mermaid") || !slices.Equal(rs[1].DependsOn, []string{rs[0].ID}) ||
		strings.Contains(rs[2].Text, "Mermaid") || len(rs[2].DependsOn) > 0 {
		t.Fatal(rs)
	}
}

func TestStructuredStageKeepsItsOriginalPlatformInRetrieval(t *testing.T) {
	cs := naturalRequirements("Create a Kotlin inventory service for the warehouse. Create an Excel chart of quarterly spending.")
	r := WorkflowRequirement{ID: "service", Text: "Build the warehouse inventory service"}
	enriched := retrievalRequirement(cs, r)
	if !strings.Contains(enriched.Text, "Kotlin") || strings.Contains(enriched.Text, "Excel") {
		t.Fatal(enriched)
	}
	if r.Text != "Build the warehouse inventory service" {
		t.Fatal(r)
	}
	r.Text = "Draft a press release"
	if enriched = retrievalRequirement(cs, r); enriched.Text != r.Text {
		t.Fatal(enriched)
	}
}
