package skillservice

// WorkflowSkillEligible checks the canonical source before a host's automatic
// prompt hook recommends loading a single named workflow.
func WorkflowSkillEligible(prompt, name string) bool {
	return workflowLoadAllowed(prompt, name)
}
