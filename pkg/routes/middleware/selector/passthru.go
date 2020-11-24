package selector

// NewPassthruHandler creates a selector that always succeeds (good for other fallbacks)
func NewPassthruHandler() SelectorGroup {
	return NewSelectorGroup(SelectorAlways)
}
