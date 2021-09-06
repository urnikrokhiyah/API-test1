package calculate

import "testing"

func TestAddPositiveNumber(t *testing.T) {
	actualValue := AdditionNumber(1, 9)
	expectedValue := 10

	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and Expected value = %d", actualValue, expectedValue)
	}
}

func TestAddNegativeNumber(t *testing.T) {
	actualValue := AdditionNumber(-1, -9)
	expectedValue := -10

	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and Expected value = %d", actualValue, expectedValue)
	}
}

func TestAddPositiveNegativeNumber(t *testing.T) {
	actualValue := AdditionNumber(1, -9)
	expectedValue := -8

	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and Expected value = %d", actualValue, expectedValue)
	}
}

func TestAddNegativePositiveNumber(t *testing.T) {
	actualValue := AdditionNumber(-1, 9)
	expectedValue := 8

	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and Expected value = %d", actualValue, expectedValue)
	}
}
