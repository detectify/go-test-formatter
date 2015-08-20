package tests

import (
	"testing"
)

func TestCount(t *testing.T) {
	stack := &SuiteStack{}

	stack.Push(&Suite{})

	expected := 1
	actual := stack.Count()

	if expected != actual {
		t.Errorf("Expected %d, have %d", expected, actual)
	}
}

func TestPop(t *testing.T) {
	stack := &SuiteStack{}

	suite := &Suite{}

	stack.Push(suite)

	expectedCount := 0
	actualSuite := stack.Pop()
	actualCount := len(stack.nodes)

	if expectedCount != actualCount {
		t.Errorf("Expected %d, have %d", expectedCount, actualCount)
	}

	if suite != actualSuite {
		t.Errorf("Expected %p, have %p", suite, actualSuite)
	}
}

func TestPush(t *testing.T) {
	stack := &SuiteStack{}

	stack.Push(&Suite{})

	expected := 1
	actual := len(stack.nodes)

	if expected != actual {
		t.Errorf("Expected %d, have %d", expected, actual)
	}
}
