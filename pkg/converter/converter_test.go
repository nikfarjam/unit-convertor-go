package converter

import (
	"fmt"
	"testing"
)

func TestCelsiusFahrenheit(t *testing.T) {
	result, err := ConvertUnit("celsius", "FAHRENHEIT", 97)
	if err != nil {
		t.Fatal(err)
	}

	if result != 206.6 {
		t.Fatalf("Expected resp.Unit to be '206.6' but it was %v", result)
	}
}

func TestFahrenheitCelsius(t *testing.T) {
	result, err := ConvertUnit("FAHRENHEIT", "celsius", 40)
	if err != nil {
		t.Fatal(err)
	}

	if result != 4.44 {
		t.Fatalf("Expected resp.Unit to be '4.44' but it was %v", result)
	}
}

func TestInvalidInput(t *testing.T) {
	_, err := ConvertUnit("test", "Invalid", 40)
	if err == nil {
		t.Fatal("When units are not valid ConvertUnit must return error")
	}
}

func almostEqual(v1, v2 float64) bool {
	fmt.Printf("Diff: %v", Abs(v2-v1))
	return Abs(v2-v1) < 0.001
}

func Abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
