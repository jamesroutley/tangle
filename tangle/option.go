package tangle

import "regexp"

type TanglerOption func(*Tangler)

func LanguageFilterOption(language string) TanglerOption {
	return func(tangler *Tangler) {
		filter := func(block *CodeBlock) bool {
			return block.Language == language
		}
		tangler.filters = append(tangler.filters, filter)
	}
}

func RegexFilterOption(re *regexp.Regexp) TanglerOption {
	return func(tangler *Tangler) {
		filter := func(block *CodeBlock) bool {
			return re.MatchString(block.Code)
		}
		tangler.filters = append(tangler.filters, filter)
	}
}
