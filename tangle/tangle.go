package tangle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

type CodeBlock struct {
	Language string
	Name     string
	Code     string
}

type Tangler struct {
	filters []Filter
}

func NewTangler(options ...TanglerOption) *Tangler {
	tangler := &Tangler{}
	for _, option := range options {
		option(tangler)
	}
	return tangler
}

func (t *Tangler) Tangle(sourceFiles ...string) ([]byte, error) {
	codeBlocks := map[string][]*CodeBlock{}

	// Pull code blocks from each of the source files
	for _, source := range sourceFiles {
		blocks, err := getCodeBlocksFromFile(source)
		if err != nil {
			return nil, err
		}

		codeBlocks[source] = blocks
	}

	// Filter out unnecessary code blocks
	for source, blocks := range codeBlocks {
		var filteredBlocks []*CodeBlock
		for _, block := range blocks {
			if !allFilters(block, t.filters...) {
				continue
			}
			filteredBlocks = append(filteredBlocks, block)
		}
		codeBlocks[source] = filteredBlocks
	}

	// Set up efficient data structures for output generation
	codeBlocksByName := map[string]*CodeBlock{}
	var nameOrder []string
	for _, source := range sourceFiles {
		blocks := codeBlocks[source]

		for _, block := range blocks {
			name := block.Name

			// Only write name if we haven't seen it before
			// TODO: this is subtle but crucial logic - I think this needs a
			// refactor to make this more explicit
			if _, ok := codeBlocksByName[name]; !ok {
				nameOrder = append(nameOrder, name)
			}

			// But always store name - that way later blocks can replace
			// earlier ones
			codeBlocksByName[name] = block
		}
	}

	// Generate output
	var output bytes.Buffer
	for _, name := range nameOrder {
		block := codeBlocksByName[name]

		output.WriteString(block.Code)
		output.WriteRune('\n')
	}

	return bytes.TrimSuffix(output.Bytes(), []byte("\n")), nil
}

func getCodeBlocksFromFile(filename string) ([]*CodeBlock, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fencedCodeBlocks, err := getFencedCodeBlocksFromMarkdown(source)
	if err != nil {
		return nil, err
	}

	var codeBlocks []*CodeBlock

	for i, block := range fencedCodeBlocks {
		var name string
		info := string(block.Info.Text(source))
		if infoParts := strings.Fields(info); len(infoParts) >= 2 {
			name = infoParts[1]
		}
		// Default name
		if name == "" {
			name = fmt.Sprintf("%s:%d", filename, i)
		}

		var code bytes.Buffer
		for i := 0; i < block.Lines().Len(); i++ {
			line := block.Lines().At(i)
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
