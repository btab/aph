package aph

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Project struct {
	CheckLatestGoVersion bool
	Executables          map[string][]string
	Services             map[string][]string
	SkipDirs             []string

	serviceCmdsByFile map[string]*exec.Cmd
	tmpDir            string
	watchDirs         *concurrentStringSet
}

func (p *Project) Run() {
	p.init()

	if p.CheckLatestGoVersion {
		checkLatestGoVersionNumber()
	}

	ensureFile("glide.lock", "run 'glide install'")
	checkFileNotStale("glide.lock")

	if fileExists("package.json") {
		checkFileNotStale("package.json")
		npmInstall()
		npmWatch()
	}

	watchNotifications := make(chan bool, 1)
	watchNotifications <- true
	go p.watch(watchNotifications)

	debouncedCallbackLoop(watchNotifications, 100*time.Millisecond, func() {
		log.Print("starting vet-lint-test loop")

		dirs := p.watchDirs.slice()

		if run("go vet", "go", append([]string{"vet"}, dirs...)) &&
			runForEach("golint", "golint", []string{"-set_exit_status"}, dirs) &&
			run("go test", "go", append([]string{"test"}, dirs...)) {
			p.restartServices()
			p.runExecutables()
		}
	})
}

func (p *Project) init() {
	p.serviceCmdsByFile = map[string]*exec.Cmd{}

	tmpDir, err := ioutil.TempDir("", "go-project-harness")
	if err != nil {
		log.Fatalf("ERROR: while establishing a tempdir: %s", err)
	}
	p.tmpDir = tmpDir

	p.watchDirs = newConcurrentStringSet()
	for _, dir := range findAllDirs(".", p.SkipDirs) {
		p.watchDirs.add("." + string(os.PathSeparator) + dir)
	}
}

func (p *Project) restartServices() {
	for goFile, args := range p.Services {
		if cmd := p.serviceCmdsByFile[goFile]; cmd != nil {
			kill(cmd, goFile)
		}

		exePath := filepath.Join(p.tmpDir, goFile+".exe")
		build(goFile, exePath)
		p.serviceCmdsByFile[goFile] = start(goFile, exePath, args)
	}
}

func (p *Project) runExecutables() {
	for goFile, args := range p.Executables {
		exePath := filepath.Join(p.tmpDir, goFile+".exe")
		build(goFile, exePath)
		run(goFile, exePath, args)
	}
}
