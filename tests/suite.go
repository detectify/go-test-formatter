package tests

// Package represents a collection of tests within a package
type Package struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}

// HasTestCases indicates whether the package has test cases
func (p *Package) HasTestCases() bool {
	for _, test := range p.Tests {
		if len(test.Cases) > 0 {
			return true
		}
	}
	return false
}
