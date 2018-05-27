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

// Platform defines a supported GOOS/GOARCH combinations.
type Platform struct {
	OS   string
	Arch string
}

// Platforms defines all possible GOOS/GOARCH combinations.
var Platforms = []Platform{
	{"android", "arm"},
	{"darwin", "386"},
	{"darwin", "amd64"},
	{"darwin", "arm"},
	{"darwin", "arm64"},
	{"dragonfly", "amd64"},
	{"freebsd", "386"},
	{"freebsd", "amd64"},
	{"freebsd", "arm"},
	{"linux", "386"},
	{"linux", "amd64"},
	{"linux", "arm"},
	{"linux", "arm64"},
	{"linux", "ppc64"},
	{"linux", "ppc64le"},
	{"linux", "mips"},
	{"linux", "mipsle"},
	{"linux", "mips64"},
	{"linux", "mips64le"},
	{"linux", "s390x"},
	{"netbsd", "386"},
	{"netbsd", "amd64"},
	{"netbsd", "arm"},
	{"openbsd", "386"},
	{"openbsd", "amd64"},
	{"openbsd", "arm"},
	{"plan9", "386"},
	{"plan9", "amd64"},
	{"solaris", "amd64"},
	{"windows", "386"},
	{"windows", "amd64"},
}

// GoBuild runs `go build rootPkg/...`.
func GoBuild(rootPkg string, recursiveBuild bool, OS, arch string) (*bytes.Buffer, error) {
	args := []string{"build"}
	if recursiveBuild {
		args = append(args, filepath.Join(rootPkg, "..."))
	} else {
		args = append(args, rootPkg)
	}
	cmd := exec.Command("go", args...)
	if OS != "" {
		cmd.Env = append(cmd.Env, "GOOS="+OS)
	}
	if arch != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
	}
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
