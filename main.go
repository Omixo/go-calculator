package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Response structure
type Response struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

// /calculate API handler
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	// Allow CORS for browser fetch()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Parse query params
	num1Str := r.URL.Query().Get("num1")
	num2Str := r.URL.Query().Get("num2")
	op := r.URL.Query().Get("operation")

	num1, err1 := strconv.ParseFloat(num1Str, 64)
	num2, err2 := strconv.ParseFloat(num2Str, 64)

	if err1 != nil || err2 != nil {
		json.NewEncoder(w).Encode(Response{Error: "Invalid number input"})
		return
	}

	var result float64
	switch op {
	case "add":
		result = num1 + num2
	case "sub":
		result = num1 - num2
	case "mul":
		result = num1 * num2
	case "div":
		if num2 == 0 {
			json.NewEncoder(w).Encode(Response{Error: "Cannot divide by zero"})
			return
		}
		result = num1 / num2
	default:
		json.NewEncoder(w).Encode(Response{Error: "Invalid operation"})
		return
	}

	json.NewEncoder(w).Encode(Response{Result: result})
}

func main() {
	// Serve index.html from static/ folder
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Register /calculate route
	http.HandleFunc("/calculate", calculateHandler)

	fmt.Println("✅ Server running at: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
