package tangle

import "github.com/jamesroutley/tangle/parser"

type Filter func(*parser.CodeBlock) bool

func allFilters(block *parser.CodeBlock, filters ...Filter) bool {
	for _, filter := range filters {
		if !filter(block) {
			return false
		}
	}
	return true
}
