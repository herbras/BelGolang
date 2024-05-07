package main

import (
	"testing"
)

// Test for variable declarations
func TestVariableDeclarations(t *testing.T) {
	// Test positiveNumber
	if positiveNumber != 89 {
		t.Errorf("positiveNumber = %d; expected %d", positiveNumber, 89)
	}

	// Test negativeNumber
	if negativeNumber != -1243423644 {
		t.Errorf("negativeNumber = %d; expected %d", negativeNumber, -1243423644)
	}

	// Test exist
	if exist != true {
		t.Errorf("exist = %t; expected %t", exist, true)
	}
}

// Test for constant declarations
func TestConstantDeclarations(t *testing.T) {
	// Test string constants
	if it != "ibrah" {
		t.Errorf("it = %s; expected %s", it, "ibrah")
	}

	if tof != "him" {
		t.Errorf("tof = %s; expected %s", tof, "him")
	}

	// Test boolean constant
	if er != true {
		t.Errorf("er = %t; expected %t", er, true)
	}

	// Test numeric constant
	if numeric != 1 {
		t.Errorf("numeric = %d; expected %d", numeric, 1)
	}

	// Test float constant
	if floatNum != 2.2 {
		t.Errorf("floatNum = %f; expected %f", floatNum, 2.2)
	}
}

func TestTypeExample(t *testing.T) {
	typeExample()
}

