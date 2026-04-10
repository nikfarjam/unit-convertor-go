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

func TestSameUnitConversion(t *testing.T) {

	tests := []struct {
		unit  string
		value float64
	}{
		{"celsius", 25},
		{"Celsius", 23.56},
		{"CELSIUS", 25.67},
		{"fahrenheit", 77},
		{"Fahrenheit", 77.12},
		{"FAHRENHEIT", 77.65},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%v", tt.unit, tt.value), func(t *testing.T) {
			result, err := ConvertUnit(tt.unit, tt.unit, tt.value)
			if err != nil {
				t.Fatal(err)
			}
			if result != tt.value {
				t.Errorf("for %s_%v: expected %v, got %v", tt.unit, tt.value, tt.value, result)
			}
		})
	}
}
