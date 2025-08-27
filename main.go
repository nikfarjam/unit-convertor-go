package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/nikfarjam/unit-convertor-go/pkg/converter"
)

func main() {
	http.HandleFunc("/converter", converterHandler)
	addr := ":9090"
	log.Printf("Server is running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func converterHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	req := &converter.ConverterRequest{}

	if err := dec.Decode(req); err != nil {
		log.Printf("Error: bad request %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if strings.ToUpper(req.From) != "CELSIUS" && strings.ToUpper(req.From) != "FAHRENHEIT" {
		log.Printf("Error: invalid unit")
		http.Error(w, "invalid from", http.StatusBadRequest)
		return
	}

	if strings.ToUpper(req.To) != "CELSIUS" && strings.ToUpper(req.To) != "FAHRENHEIT" {
		log.Printf("Error: invalid unit")
		http.Error(w, "invalid from", http.StatusBadRequest)
		return
	}

	resp, err := converter.ConvertUnit(*req)
	if err != nil {
		log.Printf("not able to process request %v", err)
		http.Error(w, "not able to process request", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Printf("Error: not able to encode response %s", err)
		http.Error(w, "not able to process request", http.StatusInternalServerError)
		return
	}
}
