package tangle

import (
	"bytes"
	"fmt"

	"github.com/jamesroutley/tangle/parser"
)

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

	var orderedBlocks []*parser.CodeBlock
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

func extractBlocksFromFiles(sourceFiles []string) (allBlocks []*parser.CodeBlock, err error) {
	for _, source := range sourceFiles {
		blocks, err := parser.Parse(source)
		if err != nil {
			return nil, err
		}
		allBlocks = append(allBlocks, blocks...)
	}
	return allBlocks, nil
}

func filterBlocks(blocks []*parser.CodeBlock, filters []Filter) (filtered []*parser.CodeBlock) {
	for _, block := range blocks {
		if !allFilters(block, filters...) {
			continue
		}
		filtered = append(filtered, block)
	}
	return filtered
}

func orderBlocks(blocks []*parser.CodeBlock) (orderedBlocks []*parser.CodeBlock) {
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
	blocks []*parser.CodeBlock, order []string,
) (orderedBlocks []*parser.CodeBlock, err error) {
	blocksByName := map[string]*parser.CodeBlock{}
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
