package tests

// A Result is the result of a test run.
type Result struct {
	Suites   []*Suite
	Assembly string
	RunDate  string
	RunTime  string
	Time     string
	Total    int
	Passed   int
	Failed   int
	Skipped  int
}
