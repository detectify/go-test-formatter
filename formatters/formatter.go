package formatters

import (
	"github.com/detectify/go-test-formatter/tests"
	"io"
)

// A Formatter formats a collection of test suites into some kind
// of text format.
type Formatter interface {
	Format([]*tests.Suite, io.Writer) error
}
