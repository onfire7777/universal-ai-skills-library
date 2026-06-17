package skillsync

// DeprecationNotice is the user-facing message printed when physical-copy
// adapter propagation runs. Physical copies of the wrapper skill into agent
// roots are deprecated in favor of agents calling the skill-router CLI directly
// or via the `serve` MCP server. This is guidance only; copy behavior is
// unchanged.
func DeprecationNotice() string {
	return "NOTE: physical-copy adapter propagation is deprecated. Agents should call " +
		"`skill-router route|search_skills|load_skill|compose` directly, or connect the " +
		"`skill-router serve` MCP server, instead of receiving copied skills. " +
		"See docs/ADAPTER_DEPRECATION.md."
}
