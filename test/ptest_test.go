package test

import (
	"testing"
)

func add(i1, i2 int) (result int) {
	result = i1 + i2
	return 
}

func TestAdd(t *testing.T) {
	result := add(3, 4)
	expected := 7
	if result != expected {
		t.Errorf("Add function returned incorrect result, got: %d, want: %d", result, expected)
	}
}

