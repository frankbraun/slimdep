// Package slim implements functionality to slim down vendor directories.
package slim

import (
	"bytes"
	"errors"

	"github.com/frankbraun/codechain/tree"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/slimdep/code"
	"github.com/frankbraun/slimdep/util"
	"github.com/frankbraun/slimdep/util/call"
)

// Down slims down the vendor directory for the given rootPkg.
func Down(rootPkg string, recursiveBuild bool) error {
	log.Println("slim.Down()")
	var err error
	err = setupDirs(rootPkg)
	if err != nil {
		return err
	}
	vendorDir, hiddenVendorDir, _, err := derivePaths(rootPkg)
	if err != nil {
		return err
	}
	vendorTree, err := code.ReadDir(hiddenVendorDir)
	if err != nil {
		return err
	}
	lastTreeHash, err := tree.Hash(vendorDir, nil)
	if err != nil {
		return err
	}
	var symbols []string
	log.Println("enter loop")
	iterations := 0
	for {
		iterations++
		log.Printf("iteration %d", iterations)
		// clean tree
		log.Println("clean tree")
		err = vendorTree.Clean(symbols)
		if err != nil {
			return teardownDirs(err, rootPkg)
		}
		// write clean version
		log.Printf("write clean version %d", iterations)
		err = vendorTree.Write()
		if err != nil {
			return teardownDirs(err, rootPkg)
		}
		// call goimports
		log.Println("call goimports")
		err = call.Goimports(vendorDir)
		if err != nil {
			return teardownDirs(err, rootPkg)
		}
		// make sure the tree hash changed
		treeHash, err := tree.Hash(vendorDir, nil)
		if err != nil {
			return teardownDirs(err, rootPkg)
		}
		if bytes.Equal(lastTreeHash[:], treeHash[:]) {
			err = errors.New("tree hashes didn't change")
			return teardownDirs(err, rootPkg)
		}
		lastTreeHash = treeHash
		// compile again
		log.Println("build again")
		buildErr := false
		stderr, err := call.GoBuild(rootPkg, recursiveBuild)
		if err != nil {
			buildErr = true
			log.Println("build failed")
			newSymbols, err := parseBuildError(stderr)
			if err != nil {
				return teardownDirs(err, rootPkg)
			}
			symbols = append(symbols, newSymbols...)
			symbols = util.UniqueStrings(symbols)
		}
		// compile tests again
		stderr, err = call.GoTest(rootPkg, recursiveBuild)
		if err != nil {
			log.Println("test compile failed")
			newSymbols, err := parseBuildError(stderr)
			if err != nil {
				return teardownDirs(err, rootPkg)
			}
			symbols = append(symbols, newSymbols...)
			symbols = util.UniqueStrings(symbols)
		} else if !buildErr {
			log.Printf("done after %d iterations!", iterations)
			return teardownDirs(nil, rootPkg) // done
		}
		// read vendor tree again
		vendorTree, err = code.ReadDir(hiddenVendorDir)
		if err != nil {
			return teardownDirs(err, rootPkg)
		}
	}
}
