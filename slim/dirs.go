package slim

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/slimdep/internal/def"
	"github.com/frankbraun/slimdep/util"
)

func derivePaths(rootPkg string) (
	vendorDir string,
	hiddenVendorDir string,
	failedVendorDir string,
	err error,
) {
	var gopath string
	gopath, err = util.GetGOPATH()
	if err != nil {
		return
	}
	vendorDir = filepath.Join(gopath, "src", rootPkg, def.VendorDir)
	hiddenVendorDir = filepath.Join(gopath, "src", rootPkg, def.HiddenVendorDir)
	failedVendorDir = filepath.Join(gopath, "src", rootPkg, def.FailedVendorDir)
	return
}

func setupDirs(rootPkg string) error {
	log.Println("setupDirs()")

	vendorDir, hiddenVendorDir, failedVendorDir, err := derivePaths(rootPkg)
	if err != nil {
		return err
	}

	exists, err := file.Exists(vendorDir)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("vendor folder does not exist: %s", vendorDir)
	}

	exists, err = file.Exists(failedVendorDir)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("remove directory %d", failedVendorDir)
		if err := os.RemoveAll(failedVendorDir); err != nil {
			return err
		}
	}

	exists, err = file.Exists(hiddenVendorDir)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("hidden vendor folder already exists: %s", hiddenVendorDir)
	}
	return file.CopyDir(vendorDir, hiddenVendorDir)
}

func teardownDirs(slimDownErr error, rootPkg string) error {
	log.Println("teardownDirs()")

	vendorDir, hiddenVendorDir, failedVendorDir, err := derivePaths(rootPkg)
	if err != nil {
		return err
	}

	if slimDownErr != nil {
		log.Println("failed, revert directories")
		if err := os.Rename(vendorDir, failedVendorDir); err != nil {
			return err
		}
		if err := os.Rename(hiddenVendorDir, vendorDir); err != nil {
			return err
		}
		return slimDownErr
	}
	log.Println("success, remove hidden vendor dir")
	return os.RemoveAll(hiddenVendorDir)
}
