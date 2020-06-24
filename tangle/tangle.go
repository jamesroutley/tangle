package tangle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

type CodeBlock struct {
	Language string
	// Name is the name assigned to the code block by the user. It can be empty
	Name string
	Code string
}

type Tangler struct {
	filters []Filter
	order   []string
}

func NewTangler(options ...TanglerOption) *Tangler {
	tangler := &Tangler{}
	for _, option := range options {
		option(tangler)
	}
	return tangler
}

func (t *Tangler) Tangle(sourceFiles ...string) ([]byte, error) {
	blocks, err := extractBlocksFromFiles(sourceFiles)
	if err != nil {
		return nil, err
	}

	filteredBlocks := filterBlocks(blocks, t.filters)

	var orderedBlocks []*CodeBlock
	if len(t.order) > 0 {
		var err error
		orderedBlocks, err = explicitlyOrderBlocks(filteredBlocks, t.order)
		if err != nil {
			return nil, err
		}
	} else {
		orderedBlocks = orderBlocks(filteredBlocks)
	}

	// Generate output
	var output bytes.Buffer
	for _, block := range orderedBlocks {
		output.WriteString(block.Code)
		output.WriteRune('\n')
	}

	return bytes.TrimSuffix(output.Bytes(), []byte("\n")), nil
}

func extractBlocksFromFiles(sourceFiles []string) (allBlocks []*CodeBlock, err error) {
	for _, source := range sourceFiles {
		blocks, err := extractBlocksFromFile(source)
		if err != nil {
			return nil, err
		}
		allBlocks = append(allBlocks, blocks...)
	}
	return allBlocks, nil
}

func extractBlocksFromFile(filename string) ([]*CodeBlock, error) {
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

func filterBlocks(blocks []*CodeBlock, filters []Filter) (filtered []*CodeBlock) {
	for _, block := range blocks {
		if !allFilters(block, filters...) {
			continue
		}
		filtered = append(filtered, block)
	}
	return filtered
}

func orderBlocks(blocks []*CodeBlock) (orderedBlocks []*CodeBlock) {
	namedBlockIndex := map[string]int{}
	for _, block := range blocks {
		if block.Name != "" {
			index, ok := namedBlockIndex[block.Name]
			if ok {
				// We've previously seen a block with this name, replace it
				orderedBlocks[index] = block
				continue
			}
			// else write it to the list and store the index
			orderedBlocks = append(orderedBlocks, block)
			namedBlockIndex[block.Name] = len(orderedBlocks) - 1
			continue
		}
		// Unnamed block, write it to the list
		orderedBlocks = append(orderedBlocks, block)
	}
	return orderedBlocks
}

func explicitlyOrderBlocks(
	blocks []*CodeBlock, order []string,
) (orderedBlocks []*CodeBlock, err error) {
	blocksByName := map[string]*CodeBlock{}
	for _, block := range blocks {
		if block.Name == "" {
			continue
		}
		blocksByName[block.Name] = block
	}

	for _, name := range order {
		block, ok := blocksByName[name]
		if !ok {
			return nil, fmt.Errorf("the order specifies the name %s, but there's not block with that order", name)
		}
		orderedBlocks = append(orderedBlocks, block)
	}
	return orderedBlocks, nil
}
