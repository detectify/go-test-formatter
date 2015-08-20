package tests

// A Suite represents a collection of tests.
type Suite struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}
