package converter

import (
	"fmt"
	"math"
	"strings"
)

func ConvertUnit(from string, to string, value float64) (float64, error) {
	switch fmt.Sprintf("%v_%v", strings.ToUpper(from), strings.ToUpper(to)) {
	case "CELSIUS_FAHRENHEIT":
		return roundFloat((value*9/5)+32, 2), nil
	case "FAHRENHEIT_CELSIUS":
		return roundFloat((value-32)*5/9, 2), nil
	case "CELSIUS_CELSIUS":
		return roundFloat(value, 2), nil
	case "FAHRENHEIT_FAHRENHEIT":
		return roundFloat(value, 2), nil
	default:
		return 0, fmt.Errorf("error: from %s and to %s units are not valid", from, to)
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
