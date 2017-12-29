package formatters

import (
	"io"

	"github.com/detectify/go-test-formatter/tests"
)

// A Formatter formats a collection of test suites into some kind
// of text format.
type Formatter interface {
	Format([]*tests.Package, io.Writer) error
}
