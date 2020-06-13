package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/jamesroutley/tangle/cmd"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var Usage = func() {
	fmt.Fprintln(os.Stderr, "Usage:  tangle <file.md>")
	flag.PrintDefaults()
}

var (
	outFile string
	watch   bool
)

func init() {
	flag.Usage = Usage
}

func main() {
	cmd.Execute()
}

// func main() {
// 	flag.StringVar(&outFile, "outfile", "", "The name of a file to write the output to")
// 	flag.BoolVar(&watch, "watch", false, "Watch the input file, and recompile when it changes")
// 	flag.Parse()

// 	filename := flag.Arg(0)
// 	if filename == "" {
// 		fmt.Fprintln(os.Stderr, "Error: no input file supplied")
// 		Usage()
// 		os.Exit(1)
// 	}

// 	if err := run(filename); err != nil {
// 		log.Fatal(err)
// 	}
// }

func run(filename string) error {

	if err := tangleAndWriteFile(filename); err != nil {
		return err
	}

	if watch {
		watchFile(filename)
	}

	return nil
}

func watchFile(filename string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(filename); err != nil {
		return err
	}

	log.Printf("Watching for changes to %s", filename)

	for {
		select {
		case event := <-watcher.Events:
			fmt.Println(event.Op)
			if event.Op != fsnotify.Write {
				continue
			}
			log.Printf("Change to %s. Rebuilding", event.Name)
			if err := tangleAndWriteFile(filename); err != nil {
				return err
			}

		case err := <-watcher.Errors:
			if err != nil {
				return err
			}
		}
	}
}

func tangleAndWriteFile(filename string) error {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tangled, err := Tangle(source)
	if err != nil {
		return err
	}

	writer := os.Stdout
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return err
		}
		defer f.Close()

		writer = f
	}

	fmt.Fprintf(writer, "%s", tangled)
	return nil
}

// Tangle pulls the code out of the Markdown fenced code blocks, then
// concatenates and returns them
func Tangle(source []byte) ([]byte, error) {
	parser := goldmark.DefaultParser()
	reader := text.NewReader(source)
	document := parser.Parse(reader)

	maxSectionNumber := 0
	sections := map[int]string{}

	walker := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		codeBlock, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		info := string(codeBlock.Info.Text(source))
		sectionNumber, err := getSectionNumber(info)
		if err != nil {
			return ast.WalkStop, err
		}
		if sectionNumber > maxSectionNumber {
			maxSectionNumber = sectionNumber
		}

		var lines bytes.Buffer
		for i := 0; i < codeBlock.Lines().Len(); i++ {
			line := codeBlock.Lines().At(i)
			_, err := lines.Write(line.Value(source))
			if err != nil {
				return ast.WalkStop, err
			}
		}
		sections[sectionNumber] = lines.String()

		return ast.WalkContinue, nil
	}

	if err := ast.Walk(document, walker); err != nil {
		return nil, err
	}

	var sectionNumbers []int
	for key := range sections {
		sectionNumbers = append(sectionNumbers, key)
	}
	sort.Ints(sectionNumbers)

	var tangled bytes.Buffer
	for _, num := range sectionNumbers {
		tangled.WriteString(sections[num])
		tangled.WriteRune('\n')
	}

	return bytes.TrimSuffix(tangled.Bytes(), []byte("\n")), nil
}

func getSectionNumber(info string) (int, error) {
	parts := strings.Fields(info)
	if len(parts) < 2 {
		return 0, fmt.Errorf("Expected a section number")
	}
	return strconv.Atoi(parts[1])
}
