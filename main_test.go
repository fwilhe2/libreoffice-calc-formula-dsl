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

func TestParseCsv(t *testing.T) {
	expectedConstants := make(map[string]string)
	expectedConstants["foo"] = "1"
	expectedFormulas := make(map[string]Formula)
	expectedFormulas["bar"] = Formula{Args: []string{"baz"}, Body: Variable{Name: "foo"}}

	source := "let foo = 1\ndefine bar(baz) = foo"
	actualConstants, actualFormulas := parseDSL(source)

	fmt.Println("constants:", actualConstants)
	fmt.Println("constants:", expectedConstants)
	fmt.Println("formulas:", actualFormulas)
	fmt.Println("formulas:", expectedFormulas)

	assert(t, maps.Equal(actualConstants, expectedConstants), "foo")
	// assert(t, maps.Equal(actualFormulas, expectedFormulas), "foo")
}
