package util

import (
	"math"
	"testing"
)

func TestConvertSlopeToDegrees(t *testing.T) {
	// Test 0 slope
	res := ConvertSlopeToDegrees(0)
	expected := 0.0
	if res != expected {
		t.Errorf("ConvertSlopeToDegrees(0) returned %f, expected %f", res, expected)
	}
	// Test inf slope??
	// Test positive slope 45 deg
	res = ConvertSlopeToDegrees(1)
	expected = 45.0
	if res != expected {
		t.Errorf("ConvertSlopeToDegrees(1) returned %f, expected %f", res, expected)
	}
	// Test negative slope 135 deg
	res = ConvertSlopeToDegrees(-1)
	expected = 135.0
	if res != expected {
		t.Errorf("ConvertSlopeToDegrees(1) returned %f, expected %f", res, expected)
	}
	// Test 30 deg
	res = ConvertSlopeToDegrees(1.0 / math.Sqrt(3))
	expected = 30.0
	if res != expected {
		t.Errorf("ConvertSlopeToDegrees(1/math.Sqrt(3)) returned %f, expected %f", res, expected)
	}
	// Test 60 deg
	res = ConvertSlopeToDegrees(math.Sqrt(3))
	expected = 60.0
	if res != expected {
		t.Errorf("ConvertSlopeToDegrees(math.Sqrt(3)) returned %f, expected %f", res, expected)
	}
	// Test 120 deg
	res = ConvertSlopeToDegrees(-math.Sqrt(3))
	expected = 120.0
	if res != expected {
		t.Errorf("ConvertSlopeToDegrees(-math.Sqrt(3)) returned %f, expected %f", res, expected)
	}
	// Test 150 deg
	res = ConvertSlopeToDegrees(-1.0 / math.Sqrt(3))
	expected = 150.0
	if res != expected {
		t.Errorf("ConvertSlopeToDegrees(-1/math.Sqrt(3)) returned %f, expected %f", res, expected)
	}
}
