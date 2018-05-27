// slimdep prunes vendored code via blind tree shaking.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/slimdep/slim"
	"github.com/frankbraun/slimdep/util/call"
)

func slimDep(rootPkg string, recursiveBuild bool) error {
	// make sure the root directory compiles before we start
	log.Println("build root directory")
	if stderr, err := call.GoBuild(rootPkg, recursiveBuild); err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
		return err
	}
	// make sure the tests compile before we start
	log.Println("build tests")
	if stderr, err := call.GoTest(rootPkg, recursiveBuild); err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
		return err
	}
	// slim it down
	log.Println("slim it down")
	return slim.Down(rootPkg, recursiveBuild)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags] root_package\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	recursive := flag.Bool("r", false, "Build packages recursively (append '/...')")
	verbose := flag.Bool("v", false, "Be verbose")
	flag.Usage = usage
	flag.Parse()
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if flag.NArg() == 0 {
		usage()
	}
	if err := slimDep(flag.Arg(0), *recursive); err != nil {
		fatal(err)
	}
}
