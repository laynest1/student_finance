package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"github.com/laynest1/student_finance/internal/models"
)


var (
	mu sync.Mutex
	transactions = []models.Transaction{
		{ID: 1, Amount: 300., Category: "няма", Date: "17-05-2006", Description: "взгрел кирюху на др"},
		{ID: 2, Amount: 1600., Category: "качалка", Date: "31-08-2026", Description: "покачаться на месяц"},
	}
)
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}


func CreateTransaction (w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var newTransaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&newTransaction)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTransaction.ID = int(len(transactions) + 1)
	transactions = append(transactions, newTransaction)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTransaction)

}