package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

// CalcRequest struct
type CalcRequest struct {
	Method string    `json:"method"`
	Args   []float64 `json:"args"`
}

const (
	MethodSumm = "summ"
	MethodDiv  = "div"
	MethodPow  = "pow"
)

// Errors definition
var (
	ErrorMethodNotImplemented = errors.New("method not implemented")
	ErrorNotEnoughArguments   = errors.New("not enough arguments")
	ErrorDivisionByZero       = errors.New("division by zero")
)

// Flag variable
var (
	port = flag.Int("port", 8080, "Listening port (default 8080)")
)

func main() {
	flag.Parse()

	fmt.Printf("Go level two. HomeWork-01\n-------------------------\n")
	fmt.Printf("Server listening on http://localhost:%v/calc\n", *port)

	http.HandleFunc("/calc", calcHandler)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%v", *port), nil))
}

// calcHandler function
func calcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(CalcRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		writeResponse(w, http.StatusBadRequest, "Error: JSON decoder: %v", err)
		return
	}

	result, err := Calculate(request)
	if err != nil {
		log.Printf("Error: function Calculate: %v", err)

		switch err {
		case ErrorMethodNotImplemented:
			writeResponse(w, http.StatusNotImplemented, "Error: method %v not implemented", request.Method)
		case ErrorDivisionByZero:
			writeResponse(w, http.StatusBadRequest, "Error: division by zero in method <%v>", request.Method)
		case ErrorNotEnoughArguments:
			writeResponse(w, http.StatusBadRequest, "Error: not enough arguments for method <%v>", request.Method)
		default:
			writeResponse(w, http.StatusInternalServerError, "Error: calculate %v", err)
		}
		return
	}

	writeResponse(w, http.StatusOK, "Result: %v", result)
}

// Calculate function
func Calculate(request *CalcRequest) (result float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: function Calculate: %v", r)
			err = fmt.Errorf("panic at %v", time.Now().UTC())
		}
	}()

	switch strings.ToLower(request.Method) {
	case MethodSumm:
		result, err = calcSumm(request.Args)
	case MethodDiv:
		result, err = calcDiv(request.Args)
	case MethodPow:
		result, err = calcPow(request.Args)
	default:
		err = ErrorMethodNotImplemented
	}
	return
}

// calcSumm function
func calcSumm(args []float64) (result float64, err error) {
	for _, arg := range args {
		result += arg
	}
	return
}

// calcDiv function
func calcDiv(args []float64) (result float64, err error) {
	if len(args) < 2 {
		err = ErrorNotEnoughArguments
	} else if args[1] == 0 {
		err = ErrorDivisionByZero
	} else {
		result = args[0] / args[1]
	}
	return
}

// calcPow function
func calcPow(args []float64) (result float64, err error) {
	result = math.Pow(args[0], args[1])
	// А вот здесь мы не подумали и запросто словим панику,
	// если количество элементов в слайсе request.Args окажется меньше двух
	return
}

// HELPERS

// writeResponse function
func writeResponse(w http.ResponseWriter, status int, format string, value interface{}) {
	w.WriteHeader(status)
	fmt.Fprintf(w, format, value)
}
