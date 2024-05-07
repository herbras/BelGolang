package main

import (
	"os"
	"testing"
)

// Test for add function
func TestAdd(t *testing.T) {
	result := add(5, 3)
	expected := 8
	if result != expected {
		t.Errorf("add(5, 3) = %d; expected %d", result, expected)
	}

	result = add(-5, 3)
	expected = -2
	if result != expected {
		t.Errorf("add(-5, 3) = %d; expected %d", result, expected)
	}
}

// Test for power function
func TestPower(t *testing.T) {
	result := power(2, 3)
	expected := 8
	if result != expected {
		t.Errorf("power(2, 3) = %d; expected %d", result, expected)
	}

	result = power(5, 0)
	expected = 1
	if result != expected {
		t.Errorf("power(5, 0) = %d; expected %d", result, expected)
	}
}

// Test for factorial function
func TestFactorial(t *testing.T) {
	result := factorial(5)
	expected := 120
	if result != expected {
		t.Errorf("factorial(5) = %d; expected %d", result, expected)
	}

	result = factorial(0)
	expected = 1
	if result != expected {
		t.Errorf("factorial(0) = %d; expected %d", result, expected)
	}
}

// Test for findLargest function
func TestFindLargest(t *testing.T) {
	arr := []int{3, 7, 2, 9, 1, 5}
	result := findLargest(arr)
	expected := 9
	if result != expected {
		t.Errorf("findLargest(%v) = %d; expected %d", arr, result, expected)
	}

	arr = []int{-3, -7, -2, -1, -5}
	result = findLargest(arr)
	expected = -1
	if result != expected {
		t.Errorf("findLargest(%v) = %d; expected %d", arr, result, expected)
	}
}

// Test for countLines function
func TestCountLines(t *testing.T) {
	// Create a temporary test file
	filePath := "test_sample.txt"
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(filePath) // Clean up after test

	// Write test content
	content := "Line 1\nLine 2\nLine 3\n"
	file.WriteString(content)
	file.Close()

	// Test the function
	lineCount, err := countLines(filePath)
	if err != nil {
		t.Fatalf("countLines failed: %v", err)
	}

	expected := 3
	if lineCount != expected {
		t.Errorf("countLines(%s) = %d; expected %d", filePath, lineCount, expected)
	}
}
