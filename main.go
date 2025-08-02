package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
	Usage  string  `json:"usage,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	num1Str := r.URL.Query().Get("num1")
	num2Str := r.URL.Query().Get("num2")
	op := r.URL.Query().Get("operation")

	num1, err1 := strconv.ParseFloat(num1Str, 64)
	num2, err2 := strconv.ParseFloat(num2Str, 64)

	if err1 != nil || err2 != nil {
		json.NewEncoder(w).Encode(Response{Error: "Please provide valid numbers in 'num1' and 'num2'."})
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
			json.NewEncoder(w).Encode(Response{Error: "Division by zero is not allowed."})
			return
		}
		result = num1 / num2
	case "mod":
		result = math.Mod(num1, num2)
	case "pow":
		result = math.Pow(num1, num2)
	default:
		json.NewEncoder(w).Encode(Response{
			Error: "Invalid operation. Use add, sub, mul, div, mod, or pow.",
			Usage: "Example: /calculate?num1=10&num2=5&operation=pow",
		})
		return
	}

	// ✅ Log to terminal
	log.Printf("[%s] %s: %.2f %s %.2f = %.2f\n", time.Now().Format(time.RFC822), r.RemoteAddr, num1, op, num2, result)

	// ✅ Append to log.txt
	logLine := fmt.Sprintf("[%s] %.2f %s %.2f = %.2f\n", time.Now().Format("2006-01-02 15:04:05"), num1, op, num2, result)
	f, err := openLogFile()
	if err == nil {
		defer f.Close()
		f.WriteString(logLine)
	}

	// ✅ Return result
	json.NewEncoder(w).Encode(Response{Result: result})
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Usage: `Use /calculate?num1=10&num2=5&operation=add | sub | mul | div | mod | pow`,
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Error: "Route not found"})
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/calculate", calculateHandler)
	http.HandleFunc("/help", helpHandler)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {}) // ignore favicon
	http.HandleFunc("/404", notFoundHandler)

	fmt.Println("✅ Server running at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func openLogFile() (*os.File, error) {
	return os.OpenFile("history/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
