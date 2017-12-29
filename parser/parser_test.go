package parser

import (
	"bytes"
	"testing"

	"github.com/test-go/testify/assert"
)

// TestParserNoSuite tests parsing output without test suites
func TestParserNoSuite(t *testing.T) {
	packages, err := New().Parse(bytes.NewReader([]byte(
		`=== RUN Test1
		--- PASS: Test1 (0.00s)
		=== RUN Test2
		--- PASS: Test2 (0.01s)
		PASS
		ok      path.to/project/package    0.019s`)))

	assert.Nil(t, err)
	assert.NotNil(t, packages)
	assert.Equal(t, 1, len(packages))
	assert.Equal(t, "path.to/project/package", packages[0].Name)
	assert.Equal(t, "ok", packages[0].Status)
	assert.Equal(t, "0.019", packages[0].Time)
	assert.Equal(t, 2, len(packages[0].Tests))

	assert.Equal(t, "Test1", packages[0].Tests[0].Name)
	assert.Equal(t, "0.00", packages[0].Tests[0].Time)
	assert.True(t, packages[0].Tests[0].Passed)

	assert.Equal(t, "Test2", packages[0].Tests[1].Name)
	assert.Equal(t, "0.01", packages[0].Tests[1].Time)
	assert.True(t, packages[0].Tests[1].Passed)
}

// TestParserSuiteSequential tests parsing output with test suites and sequential test cases
func TestParserSuiteSequential(t *testing.T) {
	packages, err := New().Parse(bytes.NewReader([]byte(
		`?       path.to/project/package1      [no test files]
		?       path.to/project/package1/subpackage        [no test files]
		=== RUN   TestSuite1
		=== RUN   TestFunction1
		--- PASS: TestFunction1 (0.01s)
		=== RUN   TestFunction2
		2018-01-16T10:49:44Z [Debug] debug log
		2018-01-16T10:49:44Z [Error] error log
		--- PASS: TestFunction2 (0.00s)
		--- PASS: TestSuite1 (0.01s)
		PASS
		=== RUN   TestSuite2
		=== RUN   TestFunction3
		--- PASS: TestFunction3 (0.00s)
		=== RUN   TestFunction4
		2018-01-16T10:49:44Z [Debug] debug log
		2018-01-16T10:49:44Z [Error] error log
		--- PASS: TestFunction4 (0.00s)
		--- PASS: TestSuite2 (0.00s)
		PASS
		ok      path.to/project/package2     0.137s`)))

	assert.Nil(t, err)
	assert.NotNil(t, packages)
	assert.Equal(t, 1, len(packages))

	assert.Equal(t, "path.to/project/package2", packages[0].Name)
	assert.Equal(t, "ok", packages[0].Status)
	assert.Equal(t, "0.137", packages[0].Time)
	assert.Equal(t, 2, len(packages[0].Tests))

	assert.Equal(t, "TestSuite1", packages[0].Tests[0].Name)
	assert.Equal(t, "0.01", packages[0].Tests[0].Time)
	assert.True(t, packages[0].Tests[0].Passed)
	assert.Equal(t, 2, len(packages[0].Tests[0].Cases))

	assert.Equal(t, "TestFunction1", packages[0].Tests[0].Cases["TestFunction1"].Name)
	assert.Equal(t, "0.01", packages[0].Tests[0].Cases["TestFunction1"].Time)
	assert.True(t, packages[0].Tests[0].Cases["TestFunction1"].Passed)

	assert.Equal(t, "TestSuite2", packages[0].Tests[1].Name)
	assert.Equal(t, "0.00", packages[0].Tests[1].Time)
	assert.True(t, packages[0].Tests[1].Passed)
	assert.Equal(t, 2, len(packages[0].Tests[1].Cases))
}

// TestParserSuiteGrouped tests parsing output with test suites and grouped test cases
func TestParserSuiteGrouped(t *testing.T) {
	packages, err := New().Parse(bytes.NewReader([]byte(
		`?       path.to/project/package1      [no test files]
		=== RUN   TestSuite
		=== RUN   TestSuite/TestFunction1
		=== RUN   TestSuite/TestFunction2
		2018-01-16T10:49:44Z [Debug] debug log
		2018-01-16T10:49:44Z [Error] error log
		=== RUN   TestSuite/TestFunction3
		=== RUN   TestSuite/TestFunction4
		2018-01-16T10:49:44Z [Debug] debug log
		2018-01-16T10:49:44Z [Error] error log
		--- PASS: TestSuite (0.01s)
			--- PASS: TestSuite/TestFunction1 (0.00s)
			--- PASS: TestSuite/TestFunction2 (0.00s)
			--- PASS: TestSuite/TestFunction3 (0.00s)
			--- PASS: TestSuite/TestFunction4 (0.00s)
		PASS
		ok      path.to/project/package2     0.210s
		?       path.to/project/package3    [no test files]`)))

	assert.Nil(t, err)
	assert.NotNil(t, packages)
	assert.Equal(t, 1, len(packages))
	assert.Equal(t, "path.to/project/package2", packages[0].Name)
	assert.Equal(t, "ok", packages[0].Status)
	assert.Equal(t, "0.210", packages[0].Time)
	assert.Equal(t, 1, len(packages[0].Tests))

	assert.Equal(t, "TestSuite", packages[0].Tests[0].Name)
	assert.Equal(t, "0.01", packages[0].Tests[0].Time)
	assert.True(t, packages[0].Tests[0].Passed)
	assert.Equal(t, 4, len(packages[0].Tests[0].Cases))

	assert.Equal(t, "TestSuite/TestFunction1", packages[0].Tests[0].Cases["TestSuite/TestFunction1"].Name)
	assert.Equal(t, "0.00", packages[0].Tests[0].Cases["TestSuite/TestFunction1"].Time)
	assert.True(t, packages[0].Tests[0].Cases["TestSuite/TestFunction1"].Passed)
}
