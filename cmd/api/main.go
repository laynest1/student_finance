package main

import (
	"net/http"
	"github.com/laynest1/student_finance/internal/handlers"
)

func main() {
	http.HandleFunc("/api/transactions", handlers.GetTransaction)
	http.ListenAndServe(":7777", nil)
}