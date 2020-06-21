package cmd

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type fileWatcher struct {
	watcher *fsnotify.Watcher
	events  chan struct{}
}

func newFileWatcher(filenames ...string) (*fileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// TODO: this does an amount of relative file string matching, which I'm
	// not will work in every scenario.
	filenameSet := map[string]bool{}
	for _, filename := range filenames {
		filenameSet[filename] = true
	}

	fw := &fileWatcher{
		watcher: watcher,
		events:  make(chan struct{}, 1),
	}

	go func() {
		// TODO: we can optimise this by writing to fw.events max once per e.g.
		// 100ms
		for {
			select {
			case event := <-watcher.Events:
				if !filenameSet[event.Name] {
					continue
				}
				fw.events <- struct{}{}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// Vim saves files by saving a copy and then renaming it to the filename.
	// This breaks fsnotify, and means we can't watch the filename directly.
	// Instead, we watch the parent directory, and filter for the file name
	// when we receive an event
	for _, filename := range filenames {
		if err := fw.watcher.Add(filepath.Dir(filename)); err != nil {
			return nil, err
		}
	}

	// hack: write an initial event to get tangle to compile when the command
	// if first run
	fw.events <- struct{}{}

	return fw, nil
}
