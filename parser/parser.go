package parser

// package parser is a non-generic Markdown parser which returns a list of the
// code blocks in a Markdown file.

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type CodeBlock struct {
	Language string
	// Name is the name assigned to the code block by the user. It can be empty
	Name string
	Code string
}

func Parse(filename string) ([]*CodeBlock, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fencedCodeBlocks, err := getFencedCodeBlocksFromMarkdown(source)
	if err != nil {
		return nil, err
	}

	var codeBlocks []*CodeBlock

	for _, block := range fencedCodeBlocks {
		var name string
		info := string(block.Info.Text(source))
		if infoParts := strings.Fields(info); len(infoParts) >= 2 {
			name = infoParts[1]
		}

		var code bytes.Buffer
		for j := 0; j < block.Lines().Len(); j++ {
			line := block.Lines().At(j)
			code.Write(line.Value(source))
		}

		codeBlock := &CodeBlock{
			Language: string(block.Language(source)),
			Name:     name,
			Code:     code.String(),
		}

		codeBlocks = append(codeBlocks, codeBlock)
	}

	return codeBlocks, nil
}

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
