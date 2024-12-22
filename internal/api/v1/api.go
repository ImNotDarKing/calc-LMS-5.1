package v1

import (
	"bytes"
	"calc-LMS/internal/calculator"
	"encoding/json"
	"io"
	"math"
	"net/http"
)


type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}


func Calculate(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			http.Error(w, `{"error": "Internal Server Error}`, http.StatusInternalServerError)
		}
	}()
	if r.Method != "POST" {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := calculator.Calc(req.Expression)
	response := Response{}

	if err != nil || math.IsNaN(result) {
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	} else {
		response.Result = result
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
