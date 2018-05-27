// Package code implement functionality to handle code trees.
package code

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/frankbraun/codechain/util"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/slimdep/internal/def"
)

// Tree defines a code tree.
type Tree struct {
	fset    *token.FileSet
	allPkgs map[string]map[string]*ast.Package
}

// ReadDir reads a code tree from given path.
func ReadDir(path string) (*Tree, error) {
	var tree Tree
	tree.fset = token.NewFileSet()
	tree.allPkgs = make(map[string]map[string]*ast.Package)
	mode := parser.ParseComments
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			pkgs, err := parser.ParseDir(tree.fset, path, nil, mode)
			if err != nil {
				return err
			}
			if len(pkgs) > 0 {
				tree.allPkgs[path] = pkgs
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &tree, nil
}

func cleanFile(fset *token.FileSet, file *ast.File, symbols []string) {
	// save comment map
	cmap := ast.NewCommentMap(fset, file, file.Comments)
	// clean declarations
	var decls []ast.Decl
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			//log.Printf("%s.%s\n", file.Name, d.Name.Name)
			if d.Name.Name == "init" || util.ContainsString(symbols, d.Name.Name) {
				//log.Printf("found: %s\n", d.Name.Name)
				decls = append(decls, decl)
			}
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				decls = append(decls, decl)
			case token.CONST:
				fallthrough
			case token.VAR:
				var match bool
			outer:
				for _, spec := range d.Specs {
					valueSpec := spec.(*ast.ValueSpec)
					for _, name := range valueSpec.Names {
						log.Printf("%s.%s\n", file.Name, name.Name)
						if util.ContainsString(symbols, name.Name) {
							log.Printf("found: %s\n", name.Name)
							match = true
							break outer
						}
					}
				}
				if match {
					decls = append(decls, decl)
				}
			case token.TYPE:
				var keptSpecs []ast.Spec
				for _, spec := range d.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					if util.ContainsString(symbols, typeSpec.Name.Name) {
						keptSpecs = append(keptSpecs, spec)
					}
				}
				if len(keptSpecs) > 0 {
					d.Specs = keptSpecs
					decls = append(decls, decl)
				}
			}
		}
	}
	file.Decls = decls
	// comments
	file.Comments = cmap.Filter(file).Comments()
}

func cleanPackage(fset *token.FileSet, pkg *ast.Package, symbolMap map[string][]string) *ast.Package {
	for _, file := range pkg.Files {
		cleanFile(fset, file, symbolMap[pkg.Name])
	}
	return pkg
}

func createSymbolMap(symbols []string) (map[string][]string, error) {
	symbolMap := make(map[string][]string)
	for _, symbol := range symbols {
		parts := strings.Split(symbol, ".")
		if len(parts) != 2 {
			return nil, fmt.Errorf("code: cannot parse symbol: %s", symbol)
		}
		pkg := parts[0]
		symbolMap[pkg] = append(symbolMap[pkg], parts[1])
		// TODO: hack
		if strings.HasPrefix(pkg, "go-") {
			pkg = strings.TrimPrefix(pkg, "go-")
			symbolMap[pkg] = append(symbolMap[pkg], parts[1])
		}
	}
	return symbolMap, nil
}

// Clean the given tree from all symbols except the given ones.
func (t *Tree) Clean(symbols []string) error {
	symbolMap, err := createSymbolMap(symbols)
	if err != nil {
		return err
	}
	for path, pkgs := range t.allPkgs {
		cleanPkgs := make(map[string]*ast.Package)
		for name, pkg := range pkgs {
			if name == "main" {
				// skip main packages
				continue
			}
			cleanPkgs[name] = cleanPackage(t.fset, pkg, symbolMap)
		}
		t.allPkgs[path] = cleanPkgs
	}
	return nil
}

func (t *Tree) Write() error {
	for _, pkgs := range t.allPkgs {
		for _, pkg := range pkgs {
			for filename, file := range pkg.Files {
				/*
					if len(file.Decls) == 0 {
						if err := os.RemoveAll(filename); err != nil {
							return err
						}
						continue
					}
				*/
				filename = strings.Replace(filename, def.HiddenVendorDir, def.VendorDir, 1)
				f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0755)
				if err != nil {
					return err
				}
				if err := format.Node(f, t.fset, file); err != nil {
					return err
				}
				if err := f.Close(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
