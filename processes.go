package aph

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func execute(name, path, verb string, args []string, f func(*exec.Cmd) error, wait bool) (*exec.Cmd, bool) {
	log.Printf("%s %s", verb, name)

	cmd := exec.Command(path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := f(cmd); err != nil {
		log.Printf("ERROR: while %s %s: %s", verb, name, err)
		return cmd, false
	}

	return cmd, true
}

func kill(cmd *exec.Cmd, name string) {
	log.Printf("killing %s", name)

	if err := cmd.Process.Kill(); err != nil {
		log.Printf("ERROR: while killing old process for %s: %s", name, err)
	}

	if err := cmd.Wait(); err != nil && err.Error() != "exit status 1" {
		log.Printf("ERROR: while waiting on old, killed process for %s: %s", name, err)
	}
}

func run(name, path string, args []string) bool {
	f := func(cmd *exec.Cmd) error { return cmd.Run() }
	_, ok := execute(name, path, "running", args, f, true)
	return ok
}

func runForEach(name, path string, commonArgs, eachArgs []string) bool {
	for _, a := range eachArgs {
		if !run(fmt.Sprintf("%s (%s)", name, a), path, append(commonArgs, a)) {
			return false
		}
	}

	return true
}

func start(name, path string, args []string) *exec.Cmd {
	f := func(cmd *exec.Cmd) error { return cmd.Start() }
	cmd, _ := execute(name, path, "starting", args, f, false)
	return cmd
}
