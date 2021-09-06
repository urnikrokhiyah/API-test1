package calculate

import "testing"

func TestSubstractPositiveNumber(t *testing.T) {
	actualValue := SubstractionNumber(1, 9)
	expectedValue := -8
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestSubstractNegativeNumber(t *testing.T) {
	actualValue := SubstractionNumber(-1, -9)
	expectedValue := 8
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestSubstractPositiveNegativeNumber(t *testing.T) {
	actualValue := SubstractionNumber(1, -9)
	expectedValue := 10
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestSubstractNegativePositiveNumber(t *testing.T) {
	actualValue := SubstractionNumber(-1, 9)
	expectedValue := -10
	if actualValue != expectedValue {
		t.Errorf("Actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}
