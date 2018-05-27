// Package call implements wrapper functions to call Go binaries.
package call

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/frankbraun/codechain/util/log"
)

// GoBuild runs `go build rootPkg/...`.
func GoBuild(rootPkg string, recursiveBuild bool) (*bytes.Buffer, error) {
	args := []string{"build"}
	if recursiveBuild {
		args = append(args, filepath.Join(rootPkg, "..."))
	} else {
		args = append(args, rootPkg)
	}
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	log.Println("go " + strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		return &stderr, err
	}
	return nil, nil
}

// GoTest runs `go test -run UZy65SSLuQXemWfgjK4EO4WqVHxTLbCR rootPkg/...`.
func GoTest(rootPkg string, recursiveBuild bool) (*bytes.Buffer, error) {
	args := []string{
		"test",
		"-run", "UZy65SSLuQXemWfgjK4EO4WqVHxTLbCR",
	}
	if recursiveBuild {
		args = append(args, filepath.Join(rootPkg, "..."))
	} else {
		args = append(args, rootPkg)
	}
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	log.Println("go " + strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		return &stderr, err
	}
	return nil, nil
}

// Goimports runs `goimports -w path`.
func Goimports(path string) error {
	cmd := exec.Command("goimports", "-w", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
