package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/laynest1/student_finance/internal/models"
)

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	transactions := []models.Transaction{
		{ID: 1, Amount: 300., Category: "няма", Date: "17-05-2006", Description: "взгрел кирюху на др"},
		{ID: 2, Amount: 1600., Category: "качалка", Date: "31-08-2026", Description: "покачаться на месяц"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
