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

func TestCompileFormula(t *testing.T) {
	constants := make(map[string]string)
	constants["foo"] = "1"
	constants["towel"] = "42"
	formulas := make(map[string]Formula)
	formulas["bar"] = Formula{Args: []string{"baz"}, Body: BinaryOp{Left: Variable{Name: "foo"}, Right: Variable{Name: "towel"}, Operator: "+"}}

	compiled := compileFormula("bar", []string{"baz"}, constants, formulas)

	assert(t, compiled == "=(1+42)", "Expected result =(1+42), got: "+compiled)
}
