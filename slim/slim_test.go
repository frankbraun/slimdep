package slim_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/slimdep/slim"
	"github.com/frankbraun/slimdep/util"
	"github.com/frankbraun/slimdep/util/call"
)

func init() {
	log.Std = log.NewStd(os.Stdout)
}

func TestSlimDown(t *testing.T) {
	tests := []string{
		"undefined",
		"interface",
	}
	gopath, err := util.GetGOPATH()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		base := filepath.Join(gopath, "src", "github.com", "frankbraun", "slimdep", "slim", "testdata")
		src := filepath.Join(base, test+"_vendor")
		dst := filepath.Join(base, test, "vendor")
		if err := file.CopyDir(src, dst); err != nil {
			t.Error(err)
		} else {
			rootPkg := "github.com/frankbraun/slimdep/slim/testdata/" + test
			stderr, err := call.GoBuild(rootPkg, false, "", "")
			if err != nil {
				fmt.Fprint(os.Stderr, stderr.String())
				t.Error(err)
			} else {
				err = slim.Down(rootPkg, false, false)
				if err != nil {
					t.Error(err)
				}
			}
			os.RemoveAll(dst)
		}
	}
}
