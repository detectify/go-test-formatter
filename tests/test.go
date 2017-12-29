package tests

// Test is a single executed test
type Test struct {
	Name    string
	Time    string
	Message string
	Failed  bool
	Skipped bool
	Passed  bool
	Cases   map[string]*TestCase
}

// TestCase is a single executed test case
type TestCase struct {
	Name    string
	Time    string
	Message string
	Failed  bool
	Skipped bool
	Passed  bool
}
