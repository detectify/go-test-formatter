package tests

// A Test is a single executed test.
type Test struct {
	Name    string
	Time    string
	Message string
	Failed  bool
	Skipped bool
	Passed  bool
}
