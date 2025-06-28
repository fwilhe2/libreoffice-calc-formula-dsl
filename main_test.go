package main

import (
	"fmt"
	"maps"
	"testing"
)

func assert(t *testing.T, condition bool, message string) {
	if !condition {
		t.Error(message)
	}
}

func TestParseDsl(t *testing.T) {
	expectedConstants := make(map[string]string)
	expectedConstants["foo"] = "1"
	expectedFormulas := make(map[string]Formula)
	expectedFormulas["bar"] = Formula{Args: []string{"baz"}, Body: Variable{Name: "foo"}}

	source := "let foo = 1\ndefine bar(baz) = foo"
	actualConstants, actualFormulas := parseDSL(source)

	assert(t, maps.Equal(actualConstants, expectedConstants), fmt.Sprintf("Expected %v to equal %v", actualConstants, expectedConstants))
	assert(t, actualFormulas["bar"].Args[0] == expectedFormulas["bar"].Args[0], fmt.Sprintf("Expected %v to equal %v", actualFormulas["bar"].Args[0], expectedFormulas["bar"].Args[0]))
	assert(t, actualFormulas["bar"].Body.String() == expectedFormulas["bar"].Body.String(), fmt.Sprintf("Expected %v to equal %v", actualFormulas["bar"].Body.String(), expectedFormulas["bar"].Body.String()))
	// assert(t, maps.Equal(actualFormulas, expectedFormulas), "foo")
}
