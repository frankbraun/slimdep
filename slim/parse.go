package slim

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/frankbraun/codechain/util/log"
)

const ident = "\\p{L}[\\p{L}\\p{N}]*" // Go identifier

var (
	constInitializerRE      = regexp.MustCompile("const initializer .* is not a constant")
	cannotUseRE             = regexp.MustCompile("cannot use .* as type .* in")
	missingMethodRE         = regexp.MustCompile(fmt.Sprintf("missing (%s) method", ident))
	missingMethodOtherPkgRE = regexp.MustCompile(fmt.Sprintf("\t(%s)\\.%s does not implement %s\\.Interface \\(missing %s method\\)",
		ident, ident, ident, ident))
	missingMethodSamePkgRE = regexp.MustCompile(fmt.Sprintf("\\(type \\*?(%s)\\.%s has no field or method (%s)\\)",
		ident, ident, ident))
	missingMethodNoPkgRE = regexp.MustCompile(fmt.Sprintf("\\(type \\*?%s has no field or method (%s)",
		ident, ident))
)

// parseBuildError returns list of undefined symbols parsed from stderr.
func parseBuildError(stderr *bytes.Buffer) ([]string, error) {
	var (
		symbols []string
		prev    string
	)
	s := bufio.NewScanner(stderr)
	for s.Scan() {
		line := s.Text()

		// comments
		if strings.HasPrefix(line, "#") {
			log.Printf("match: %s\n", line)
			continue
		}

		// too many errors
		if strings.Contains(line, "too many errors") {
			log.Printf("match: %s\n", line)
			continue
		}

		// interface with no methods (previous line should contain missing type)
		if strings.Contains(line, "(type interface {} is interface with no methods)") {
			log.Printf("match: %s\n", line)
			continue
		}

		// const initializer
		if constInitializerRE.MatchString(line) {
			log.Printf("match: %s\n", line)
			continue
		}

		// cannot use line
		if cannotUseRE.MatchString(line) {
			log.Printf("match: %s\n", line)
			// skip it, will be processed in next line
			prev = line
			continue
		}

		// missing method
		match := missingMethodRE.FindStringSubmatch(line)
		if match != nil {
			log.Printf("match: %s\n", line)
			ident := match[1]
			// determine package
			var pkg string
			match := missingMethodOtherPkgRE.FindStringSubmatch(line)
			if match != nil {
				pkg = match[1]
			} else {
				parts := strings.Split(prev, ": ")
				pkg = filepath.Base(filepath.Dir(parts[0]))
			}
			fullSymbol := pkg + "." + ident
			symbols = append(symbols, fullSymbol)
			prev = line
			continue
		}

		match = missingMethodSamePkgRE.FindStringSubmatch(line)
		if match != nil {
			log.Printf("match: %s\n", line)
			fullSymbol := match[1] + "." + match[2]
			symbols = append(symbols, fullSymbol)
			prev = line
			continue
		}

		match = missingMethodNoPkgRE.FindStringSubmatch(line)
		if match != nil {
			log.Printf("match: %s\n", line)
			ident := match[1]
			parts := strings.Split(line, ": ")
			pkg := filepath.Base(filepath.Dir(parts[0]))
			fullSymbol := pkg + "." + ident
			symbols = append(symbols, fullSymbol)
			prev = line
			continue
		}

		// undefined symbol
		parts := strings.Split(line, ": ")
		if len(parts) == 3 && parts[1] == "undefined" {
			log.Printf("match: %s\n", line)
			fullSymbol := parts[2]
			particles := strings.Split(fullSymbol, ".")
			l := len(particles)
			if l < 2 {
				// we are dealing with an internal symbol, try to figure out
				// package name from path
				pkg := filepath.Base(filepath.Dir(parts[0]))
				fullSymbol = pkg + "." + fullSymbol
			} else if l > 2 {
				pkg := strings.Join(particles[:l-1], ".")
				pkg = strings.TrimFunc(pkg, func(r rune) bool {
					return r == '"'
				})
				pkg = filepath.Base(pkg)
				fullSymbol = pkg + "." + particles[l-1]
			}
			symbols = append(symbols, fullSymbol)
		} else {
			return nil, fmt.Errorf("slim: cannot parse error: %s", line)
		}

		prev = line
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return symbols, nil
}
