package aph

import (
	"log"
	"os/exec"
)

func build(inPath, outPath string) {
	if out, err := exec.Command("go", "build", "-o", outPath, inPath).CombinedOutput(); err != nil {
		log.Fatalf("ERROR: while building %s: %s -> %s", inPath, err, out)
	}
}
