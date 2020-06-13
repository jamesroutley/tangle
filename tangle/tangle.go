package tangle

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Tangler struct{}

func NewTangler() *Tangler {
	return &Tangler{}
}

func getCodeBlocksFromMarkdown(source []byte) ([]*ast.FencedCodeBlock, error) {
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

// Tangle pulls the code out of the Markdown fenced code blocks, then
// concatenates and returns them
func Tangle(source []byte) ([]byte, error) {
	codeBlocks, err := getCodeBlocksFromMarkdown(source)
	if err != nil {
		return nil, err
	}

	var output [][]byte
	for _, codeBlock := range codeBlocks {
		lines, err := getCodeBlockLines(codeBlock, source)
		if err != nil {
			return nil, err
		}
		output = append(output, lines)
	}

	return bytes.Join(output, []byte("\n")), nil
}

func getCodeBlockLines(codeBlock *ast.FencedCodeBlock, source []byte) ([]byte, error) {
	var lines bytes.Buffer
	for i := 0; i < codeBlock.Lines().Len(); i++ {
		line := codeBlock.Lines().At(i)
		_, err := lines.Write(line.Value(source))
		if err != nil {
			return nil, err
		}
	}
	return lines.Bytes(), nil
}

func getSectionNumber(info string) (int, error) {
	parts := strings.Fields(info)
	if len(parts) < 2 {
		return 0, fmt.Errorf("Expected a section number")
	}
	return strconv.Atoi(parts[1])
}
