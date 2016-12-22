package project

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func checkFileNotStale(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.ModTime().Before(time.Now().AddDate(0, -11, 0)) {
		log.Printf("WARNING: %s is over a month old", path)
	} else {
		// log.Printf("%s is fresh", path)
	}
}

func ensureFile(path, remedy string) {
	if !fileExists(path) {
		log.Fatalf("ERROR: %s does not exist: %s", path, remedy)
	}
}

func findAllDirs(path string, skipDirs []string) []string {
	dirs := []string{}
	skipSet := newStringSet(skipDirs...)

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if path == "." || !info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, ".") || skipSet[path] {
			return filepath.SkipDir
		}

		dirs = append(dirs, path)
		return nil
	})

	return dirs
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
