package converter

import (
	"fmt"
	"testing"
)

func TestCelsiusFahrenheit(t *testing.T) {
	req := NewConverterRequest(97, "celsius", "FAHRENHEIT")

	resp, err := ConvertUnit(*req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Unit != "FAHRENHEIT" {
		t.Fatalf("Expected resp.Unit to be 'FAHRENHEIT' but it was %v", resp.Unit)
	}
	if resp.Value != 206.6 {
		t.Fatalf("Expected resp.Unit to be '206.6' but it was %v", resp.Value)
	}
}

func TestFahrenheitCelsius(t *testing.T) {
	req := NewConverterRequest(40, "FAHRENHEIT", "celsius")

	resp, err := ConvertUnit(*req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Unit != "CELSIUS" {
		t.Fatalf("Expected resp.Unit to be 'CELSIUS' but it was %v", resp.Unit)
	}
	if !almostEqual(resp.Value, 4.44) {
		t.Fatalf("Expected resp.Unit to be '4.44' but it was %v", resp.Value)
	}
}

func TestInvalidInput(t *testing.T) {
	req := NewConverterRequest(40, "test", "Invalid")

	_, err := ConvertUnit(*req)
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
