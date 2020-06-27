package tangle

import (
	"regexp"

	"github.com/jamesroutley/tangle/parser"
)

type TanglerOption func(*Tangler)

func LanguageFilterOption(language string) TanglerOption {
	return func(tangler *Tangler) {
		filter := func(block *parser.CodeBlock) bool {
			return block.Language == language
		}
		tangler.filters = append(tangler.filters, filter)
	}
}

func RegexFilterOption(re *regexp.Regexp) TanglerOption {
	return func(tangler *Tangler) {
		filter := func(block *parser.CodeBlock) bool {
			return re.MatchString(block.Code)
		}
		tangler.filters = append(tangler.filters, filter)
	}
}

func CustomOrderOption(blockNames []string) TanglerOption {
	return func(tangler *Tangler) {
		tangler.order = blockNames
	}
}
