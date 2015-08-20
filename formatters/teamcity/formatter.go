package teamcity

import (
	"fmt"
	"github.com/detectify/go-test-formatter/formatters"
	"github.com/detectify/go-test-formatter/tests"
	"io"
	"strings"
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
func (f *Formatter) Format(suites []*tests.Suite, writer io.Writer) error {
	for _, suite := range suites {
		return f.printSuite(suite, writer)
	}

	return nil
}

func (f *Formatter) printSuite(suite *tests.Suite, writer io.Writer) error {
	suiteName := escape(suite.Name)

	fmt.Fprintf(writer, testSuiteStarted, suiteName)

	for _, test := range suite.Tests {
		err := f.printTest(test, writer)

		if err != nil {
			return err
		}
	}

	fmt.Fprintf(writer, testSuiteFinished, suiteName)

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
