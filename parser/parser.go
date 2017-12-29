package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/detectify/go-test-formatter/tests"
)

var (
	findStart     = regexp.MustCompile("^=== RUN:?[[:space:]]+([a-zA-Z_][^[:space:]]*)").FindStringSubmatch
	findEnd       = regexp.MustCompile("^--- (PASS|FAIL|SKIP):[[:space:]]+([a-zA-Z_][^[:space:]]*) \\((\\d+(.\\d+)?)").FindStringSubmatch
	findPackage   = regexp.MustCompile("^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(\\d+.\\d+)").FindStringSubmatch
	isStatus      = regexp.MustCompile("^(PASS|FAIL|SKIP)$").MatchString
	isNoFiles     = regexp.MustCompile("^\\?.*\\[no test files\\]$").MatchString
	isBuildFailed = regexp.MustCompile("^FAIL.*\\[(build|setup) failed\\]$").MatchString
	isExit        = regexp.MustCompile("^exit status -?\\d+").MatchString
)

// A Parser parses `go test -v` output into an object model.
type Parser struct {
	packages     []*tests.Package
	packageTests []*tests.Test
	current      *tests.Test
}

// New creates a new parser.
func New() *Parser {
	return &Parser{
		packages:     make([]*tests.Package, 0),
		packageTests: make([]*tests.Test, 0),
	}
}

// Parse parses `go test -v` output into an object model.
func (p *Parser) Parse(reader io.Reader) ([]*tests.Package, error) {
	scanner := bufio.NewScanner(reader)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || isNoFiles(line) || isExit(line) || isStatus(line) {
			continue
		}

		if isBuildFailed(line) {
			return nil, fmt.Errorf("%d: package build failed: %s", lineNumber, line)
		}

		tokens := findStart(line)
		if tokens != nil {
			p.parseStart(tokens)
			continue
		}

		tokens = findEnd(line)
		if tokens != nil {
			if !p.parseEnd(tokens) {
				return nil, fmt.Errorf("%d: orphan/invalid end test", lineNumber)
			}
			continue
		}

		tokens = findPackage(line)
		if tokens != nil {
			p.parsePackage(tokens)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	packages := p.packages

	return packages, nil
}

func (p *Parser) parseStart(tokens []string) {
	if p.current == nil || p.current.Time != "" {
		// there is no test or the test has ended
		// assume test case of main
		p.current = &tests.Test{
			Name:  tokens[1],
			Cases: make(map[string]*tests.TestCase, 0),
		}
		p.packageTests = append(p.packageTests, p.current)
		return
	}

	// assume this is a test case
	p.current.Cases[tokens[1]] = &tests.TestCase{
		Name: tokens[1],
	}
}

func (p *Parser) parseEnd(tokens []string) bool {
	if p.current == nil {
		// there is no test
		return false
	}

	if tokens[2] == p.current.Name {
		// ends the test
		p.current.Failed = tokens[1] == "FAIL"
		p.current.Skipped = tokens[1] == "SKIP"
		p.current.Passed = tokens[1] == "PASS"
		p.current.Time = tokens[3]
		return true
	}

	for name, test := range p.current.Cases {
		if tokens[2] == name {
			// ends a test case
			test.Failed = tokens[1] == "FAIL"
			test.Skipped = tokens[1] == "SKIP"
			test.Passed = tokens[1] == "PASS"
			test.Time = tokens[3]
			return true
		}
	}

	// name does not match any test/test case
	return false
}

func (p *Parser) parsePackage(tokens []string) {
	p.packages = append(p.packages, &tests.Package{
		Name:   tokens[2],
		Time:   tokens[3],
		Status: tokens[1],
		Tests:  p.packageTests,
	})

	p.packageTests = make([]*tests.Test, 0)
	p.current = nil
}
