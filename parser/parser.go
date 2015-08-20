package parser

import (
	"bufio"
	"fmt"
	"github.com/detectify/go-test-formatter/tests"
	"io"
	"regexp"
	"strings"
)

const (
	startPattern       = "^=== RUN:?[[:space:]]+([a-zA-Z_][^[:space:]]*)"
	endPattern         = "^--- (PASS|FAIL|SKIP):[[:space:]]+([a-zA-Z_][^[:space:]]*) \\((\\d+(.\\d+)?)"
	suitePattern       = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(\\d+.\\d+)"
	noFilesPattern     = "^\\?.*\\[no test files\\]$"
	buildFailedPattern = "^FAIL.*\\[(build|setup) failed\\]$"
)

// A Parser parses `go test -v` output into an object model.
type Parser struct {
}

// New creates a new parser.
func New() *Parser {
	return &Parser{}
}

// Parse parses `go test -v` output into an object model.
func (p *Parser) Parse(reader io.Reader) ([]*tests.Suite, error) {
	findStart := regexp.MustCompile(startPattern).FindStringSubmatch
	findEnd := regexp.MustCompile(endPattern).FindStringSubmatch
	findSuite := regexp.MustCompile(suitePattern).FindStringSubmatch
	isNoFiles := regexp.MustCompile(noFilesPattern).MatchString
	isBuildFailed := regexp.MustCompile(buildFailedPattern).MatchString
	isExit := regexp.MustCompile("^exit status -?\\d+").MatchString

	var suites []*tests.Suite

	suiteStack := &tests.SuiteStack{}

	var currentTest *tests.Test
	var currentSuite *tests.Suite
	var output []string

	handlePanic := func() {
		currentTest.Failed = true
		currentTest.Skipped = false
		currentTest.Time = "N/A"
		currentSuite.Tests = append(currentSuite.Tests, currentTest)
		currentTest = nil
	}

	appendError := func() error {
		if len(output) > 0 && currentSuite != nil && len(currentSuite.Tests) > 0 {
			message := strings.Join(output, "\n")

			lastIndex := len(currentSuite.Tests) - 1

			if currentSuite.Tests[lastIndex].Message == "" {
				currentSuite.Tests[lastIndex].Message = message
			} else {
				currentSuite.Tests[lastIndex].Message += "\n" + message
			}
		}

		output = make([]string, 0)

		return nil
	}

	scanner := bufio.NewScanner(reader)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()

		// Skip folders without tests.
		if isNoFiles(line) {
			continue
		}

		if isBuildFailed(line) {
			return nil, fmt.Errorf("%d: package build failed: %s", lineNumber, line)
		}

		if currentSuite == nil {
			currentSuite = &tests.Suite{}
		}

		tokens := findStart(line)

		if tokens != nil {
			if currentTest != nil {
				if suiteStack.Count() == 0 {
					suiteStack.Push(currentSuite)

					currentSuite = &tests.Suite{Name: currentTest.Name}
				} else {
					handlePanic()
				}
			}

			if e := appendError(); e != nil {
				return nil, e
			}

			currentTest = &tests.Test{
				Name: tokens[1],
			}

			continue
		}

		tokens = findEnd(line)

		if tokens != nil {
			if currentTest == nil {
				if suiteStack.Count() > 0 {
					previousSuite := suiteStack.Pop()

					suites = append(suites, currentSuite)

					currentSuite = previousSuite

					continue
				} else {
					return nil, fmt.Errorf("%d: orphan end test", lineNumber)
				}
			}

			if tokens[2] != currentTest.Name {
				err := fmt.Errorf("%d: name mismatch (try disabling parallel mode)", lineNumber)

				return nil, err
			}

			currentTest.Failed = (tokens[1] == "FAIL") // || (failOnRace && hasDatarace(output))
			currentTest.Skipped = tokens[1] == "SKIP"
			currentTest.Passed = tokens[1] == "PASS"
			currentTest.Time = tokens[3]
			currentTest.Message = strings.Join(output, "\n")

			currentSuite.Tests = append(currentSuite.Tests, currentTest)

			currentTest = nil

			output = make([]string, 0)

			continue
		}

		tokens = findSuite(line)

		if tokens != nil {
			if currentTest != nil {
				handlePanic()
			}

			if e := appendError(); e != nil {
				return nil, e
			}

			currentSuite.Name = tokens[2]
			currentSuite.Time = tokens[3]

			suites = append(suites, currentSuite)

			currentSuite = nil

			continue
		}

		if isExit(line) || line == "FAIL" || line == "PASS" {
			continue
		}

		output = append(output, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}
