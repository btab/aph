package aph

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

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
	log.Printf("running %s", name)

	out, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		log.Printf("ERROR: while running %s: %s", name, out)
		return false
	}

	return true
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
	log.Printf("starting %s", name)

	cmd := exec.Command(path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		log.Fatalf("ERROR: while starting %s: %s", name, err)
	}

	return cmd
}
