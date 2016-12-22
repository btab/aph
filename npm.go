package project

import (
	"os"
)

func npmInstall() {
	if !run("npm install", "npm", []string{"install"}) {
		os.Exit(1)
	}
}

func npmWatch() {
	start("npm run watch", "npm", []string{"run", "watch"})
}
