package tangle

// func watchFile(filename string) error {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		return err
// 	}
// 	defer watcher.Close()

// 	if err := watcher.Add(filename); err != nil {
// 		return err
// 	}

// 	log.Printf("Watching for changes to %s", filename)

// 	for {
// 		select {
// 		case event := <-watcher.Events:
// 			fmt.Println(event.Op)
// 			if event.Op != fsnotify.Write {
// 				continue
// 			}
// 			log.Printf("Change to %s. Rebuilding", event.Name)
// 			if err := tangleAndWriteFile(filename); err != nil {
// 				return err
// 			}

// 		case err := <-watcher.Errors:
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// }
