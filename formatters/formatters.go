package formatters

import (
	"fmt"
)

// FormatterCreator creates formatters.
type FormatterCreator func() (Formatter, error)

var registeredFormatters = make(map[string]FormatterCreator)

// Register registers a named formatter.
func Register(name string, creator FormatterCreator) {
	registeredFormatters[name] = creator
}

// Find finds a registered formatter by name.
func Find(name string) (Formatter, error) {
	if creator, ok := registeredFormatters[name]; ok {
		return creator()
	}

	return nil, fmt.Errorf("'%s' formatter not found", name)
}
