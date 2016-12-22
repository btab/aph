package aph

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func (p *Project) watch(notifications chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("ERROR: unable to build filesystem watcher: %s", err)
	}
	defer watcher.Close()

	for _, dir := range p.watchDirs.slice() {
		watcher.Add(dir)
	}

	for {
		select {
		case event := <-watcher.Events:
			notify := true

			if p.watchEventIsDirCreate(event) {
				p.watchDirs.add(event.Name)
				watcher.Add(event.Name)
			} else if p.watchEventIsDirRemove(event) {
				p.watchDirs.remove(event.Name)
				watcher.Remove(event.Name)
			} else if filepath.Ext(event.Name) != ".go" {
				notify = false
			}

			if notify {
				notifications <- true
			}

		case err := <-watcher.Errors:
			log.Fatalf("ERROR: while watching filesystem: %s", err)
		}
	}
}

func (p *Project) watchEventIsDirCreate(event fsnotify.Event) bool {
	if event.Op&fsnotify.Create == 0 {
		return false
	}

	info, err := os.Stat(event.Name)
	if err != nil {
		log.Fatalf("ERROR: while watching filesystem: %s", err)
	}

	return info.IsDir()
}

func (p *Project) watchEventIsDirRemove(event fsnotify.Event) bool {
	if event.Op&fsnotify.Remove == 0 {
		return false
	}

	return p.watchDirs.has(event.Name)
}
