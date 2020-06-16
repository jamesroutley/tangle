package tangle

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func getFencedCodeBlocksFromMarkdown(source []byte) ([]*ast.FencedCodeBlock, error) {
	parser := goldmark.DefaultParser()
	reader := text.NewReader(source)
	document := parser.Parse(reader)

	var codeBlocks []*ast.FencedCodeBlock

	walker := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		codeBlock, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		codeBlocks = append(codeBlocks, codeBlock)
		return ast.WalkContinue, nil
	}

	if err := ast.Walk(document, walker); err != nil {
		return nil, err
	}

	return codeBlocks, nil
}
