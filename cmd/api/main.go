package main

import (
	"net/http"
	"github.com/laynest1/student_finance/internal/handlers"
	"github.com/laynest1/student_finance/internal/database"
	

)

func main() {

	database.Init()
	database.Migrate()
	http.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request){
		switch r.Method {
			case http.MethodGet :
				handlers.GetTransaction(w,r)
			case http.MethodPost :
				handlers.CreateTransaction(w,r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)


	}

	})
	http.ListenAndServe(":7777", nil)
}