package tangle

type Filter func(*CodeBlock) bool

func allFilters(block *CodeBlock, filters ...Filter) bool {
	for _, filter := range filters {
		if !filter(block) {
			return false
		}
	}
	return true
}
