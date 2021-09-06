package calculate

import "testing"

func TestDivisionPositive(t *testing.T) {
	actualValue := DivisionNumber(4, 2)
	expectedValue := 2
	if actualValue != expectedValue {
		t.Errorf("actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestDivisionNegative(t *testing.T) {
	actualValue := DivisionNumber(-4, -2)
	expectedValue := 2
	if actualValue != expectedValue {
		t.Errorf("actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestDivisionPositiveNegative(t *testing.T) {
	actualValue := DivisionNumber(4, -2)
	expectedValue := -2
	if actualValue != expectedValue {
		t.Errorf("actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestDivisionNegativePositive(t *testing.T) {
	actualValue := DivisionNumber(-4, 2)
	expectedValue := -2
	if actualValue != expectedValue {
		t.Errorf("actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestDivisionZeroNumerator(t *testing.T) {
	actualValue := DivisionNumber(0, 4)
	expectedValue := 0
	if actualValue != expectedValue {
		t.Errorf("actual value = %d and expected value = %d", actualValue, expectedValue)
	}
}

func TestDivisionZeroDenominator(t *testing.T) {
	actualValue := DivisionNumber(4, 0)
	if actualValue != "division by zero" {
		t.Error("return value not correct")
	}
}
