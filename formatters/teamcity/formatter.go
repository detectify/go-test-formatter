package teamcity

import (
	"fmt"
	"io"
	"strings"

	"github.com/detectify/go-test-formatter/formatters"
	"github.com/detectify/go-test-formatter/tests"
)

const (
	teamcityTag       = "##teamcity"
	testSuiteStarted  = teamcityTag + "[testSuiteStarted name='%s']\n"
	testSuiteFinished = teamcityTag + "[testSuiteFinished name='%s']\n"
	testStarted       = teamcityTag + "[testStarted name='%s']\n"
	testFailed        = teamcityTag + "[testFailed name='%s' message='%s']\n"
	testIgnored       = teamcityTag + "[testIgnored name='%s']\n"
	testFinished      = teamcityTag + "[testFinished name='%s' duration='%s']\n"
)

// A Formatter formats test results into something that TeamCity can interpret.
type Formatter struct {
}

// Format formats a collection of test suites into something TeamCity can
// interpret.
func (f *Formatter) Format(packages []*tests.Package, writer io.Writer) error {
	for _, pack := range packages {
		err := f.printPackage(pack, writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Formatter) printPackage(pack *tests.Package, writer io.Writer) error {
	packName := escape(pack.Name)

	if pack.HasTestCases() {
		for _, test := range pack.Tests {
			err := f.printTestSuite(test, packName, writer)

			if err != nil {
				return err
			}
		}
		return nil
	}

	fmt.Fprintf(writer, testSuiteStarted, packName)

	for _, test := range pack.Tests {
		err := f.printTest(test, writer)

		if err != nil {
			return err
		}
	}

	fmt.Fprintf(writer, testSuiteFinished, packName)

	return nil
}

func (f *Formatter) printTest(test *tests.Test, writer io.Writer) error {
	name := escape(test.Name)

	fmt.Fprintf(writer, testStarted, name)

	if test.Failed {
		message := escape(test.Message)

		fmt.Fprintf(writer, testFailed, name, message)
	} else if test.Skipped {
		fmt.Fprintf(writer, testIgnored, name)
	}

	fmt.Fprintf(writer, testFinished, name, test.Time)

	return nil
}

func (f *Formatter) printTestSuite(test *tests.Test, packName string, writer io.Writer) error {
	name := fmt.Sprintf("%s/%s", packName, escape(test.Name))

	fmt.Fprintf(writer, testSuiteStarted, name)

	for _, testCase := range test.Cases {
		err := f.printTestCase(testCase, writer)

		if err != nil {
			return err
		}
	}

	fmt.Fprintf(writer, testSuiteFinished, name)

	return nil
}

func (f *Formatter) printTestCase(test *tests.TestCase, writer io.Writer) error {
	name := escape(test.Name)

	fmt.Fprintf(writer, testStarted, name)

	if test.Failed {
		message := escape(test.Message)

		fmt.Fprintf(writer, testFailed, name, message)
	} else if test.Skipped {
		fmt.Fprintf(writer, testIgnored, name)
	}

	fmt.Fprintf(writer, testFinished, name, test.Time)

	return nil
}

func escape(output string) string {
	output = strings.Replace(output, "|", "||", -1)
	output = strings.Replace(output, "'", "|'", -1)
	output = strings.Replace(output, "\n", "|n", -1)
	output = strings.Replace(output, "\r", "|r", -1)
	output = strings.Replace(output, "[", "|[", -1)
	output = strings.Replace(output, "]", "|]", -1)

	return output
}

// New creates a new TeamCity formatter.
func New() *Formatter {
	formatter := &Formatter{}

	return formatter
}

func init() {
	formatters.Register("teamcity", func() (formatters.Formatter, error) {
		return New(), nil
	})
}
