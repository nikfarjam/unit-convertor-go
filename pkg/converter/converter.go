package converter

import (
	"fmt"
	"math"
	"strings"
)

type ConverterRequest struct {
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

func NewConverterRequest(value float64, from string, to string) *ConverterRequest {
	return &ConverterRequest{
		Value: value,
		From:  from,
		To:    to,
	}
}

type ConverterResponse struct {
	Value float64 `json:"value"`
	Unit  string  `json:"Unit"`
}

func ConvertUnit(req ConverterRequest) (ConverterResponse, error) {
	from := strings.ToUpper(req.From)
	to := strings.ToUpper(req.To)
	switch fmt.Sprintf("%v_%v", from, to) {
	case "CELSIUS_FAHRENHEIT":
		return ConverterResponse{
			Value: roundFloat((req.Value*9/5)+32, 2),
			Unit:  to,
		}, nil
	case "FAHRENHEIT_CELSIUS":
		return ConverterResponse{
			Value: roundFloat((req.Value-32)*5/9, 2),
			Unit:  to,
		}, nil
	default:
		return ConverterResponse{},
			fmt.Errorf("error: form %s and to %s units are not valid", from, to)
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
