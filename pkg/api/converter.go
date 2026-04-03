package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/nikfarjam/unit-convertor-go/pkg/converter"
)

func ConverterHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)
	dec := json.NewDecoder(r.Body)
	req := &converter.ConverterRequest{}

	if err := dec.Decode(req); err != nil {
		slog.Error("Error: bad request", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	slog.Debug("Request received", "from", req.From, "to", req.To)

	if strings.ToUpper(req.From) != "CELSIUS" && strings.ToUpper(req.From) != "FAHRENHEIT" {
		slog.Error("Error: invalid unit", "unit", req.From)
		http.Error(w, "invalid from", http.StatusBadRequest)
		return
	}

	if strings.ToUpper(req.To) != "CELSIUS" && strings.ToUpper(req.To) != "FAHRENHEIT" {
		slog.Error("Error: invalid unit", "unit", req.To)
		http.Error(w, "invalid to", http.StatusBadRequest)
		return
	}

	resp, err := converter.ConvertUnit(*req)
	if err != nil {
		slog.Error("Error: not able to process request", "error", err)
		http.Error(w, "not able to process request", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		slog.Error("Error: not able to encode response", "error", err)
		http.Error(w, "not able to process request", http.StatusInternalServerError)
		return
	}
}
