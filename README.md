## Slimdep — prune vendored Go code with blind tree shaking — α release

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/frankbraun/slimdep) [![Build Status](https://img.shields.io/travis/frankbraun/slimdep.svg?style=flat-square)](https://travis-ci.org/frankbraun/slimdep) [![Go Report Card](https://goreportcard.com/badge/github.com/frankbraun/slimdep?style=flat-square)](https://goreportcard.com/report/github.com/frankbraun/slimdep)

### Installation

```
go get -u -v github.com/frankbraun/slimdep
```

### Mode of operation

`slimdep` assumes that all dependencies are vendored into the `vendor` folder.

- Use `go/parser` to parse into AST and then `ast.Walk` to process it,
  analog to gofmt.
- Remove all functions except for `init` (maybe even remove types, variables,
  consts).
- Compile.
- Loop: analyse error message -> introduce missing function -> compile again.
- To support other platforms: cross compile (option `-a`).
