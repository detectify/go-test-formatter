package tests

// A SuiteStack adds stack operations to an array of Suites.
type SuiteStack struct {
	nodes []*Suite
}

// Count returns the number of elements currently in the stack.
func (s *SuiteStack) Count() int {
	return len(s.nodes)
}

// Pop deletes the last element of the stack and returns it.
func (s *SuiteStack) Pop() *Suite {
	if s.Count() == 0 {
		return nil
	}

	var suite *Suite

	suite, s.nodes = s.nodes[len(s.nodes)-1], s.nodes[:len(s.nodes)-1]

	return suite
}

// Push adds a new element to the end of the stack.
func (s *SuiteStack) Push(suite *Suite) {
	s.nodes = append(s.nodes, suite)
}
