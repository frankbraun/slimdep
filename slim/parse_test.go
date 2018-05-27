package slim

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"
)

var matchIdent = regexp.MustCompile(ident)

func TestMatchIdent(t *testing.T) {
	identifiers := []string{
		"ABC",
		"_ABC",
		"ABC123",
	}
	for _, identifier := range identifiers {
		if !matchIdent.MatchString(identifier) {
			t.Errorf("%s should match identifier", identifier)
		}
	}
}

func TestParseBuildErrors(t *testing.T) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", "build_errors.txt"))
	if err != nil {
		t.Fatal(err)
	}
	symbols, err := parseBuildError(bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("symbols:")
	for _, symbol := range symbols {
		fmt.Println(symbol)
	}
}
