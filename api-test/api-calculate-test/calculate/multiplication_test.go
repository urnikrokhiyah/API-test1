package calculate

import "testing"

func TestMultiplePositive(t *testing.T) {
	actualValue := MultiplicationNumber(1, 9)
	expectedValue := 9
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestMultipleNegative(t *testing.T) {
	actualValue := MultiplicationNumber(-1, -9)
	expectedValue := 9
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestMultiplePositiveNegative(t *testing.T) {
	actualValue := MultiplicationNumber(1, -9)
	expectedValue := -9
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestMultipleNegativePositive(t *testing.T) {
	actualValue := MultiplicationNumber(-1, 9)
	expectedValue := -9
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestMultipleWithZero(t *testing.T) {
	actualValue := MultiplicationNumber(1, 0)
	expectedValue := 0
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}
