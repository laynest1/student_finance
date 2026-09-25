package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/laynest1/student_finance/internal/database"
    "github.com/laynest1/student_finance/internal/models"
)


func GetTransaction(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    rows, err := database.DB.Query(r.Context(),
        "SELECT id, amount, category, date, description FROM transactions")
    if err != nil {
        log.Printf("GetTransactions query error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    transactions := []models.Transaction{}
    for rows.Next() {
        var t models.Transaction
        if err := rows.Scan(&t.ID, &t.Amount, &t.Category, &t.Date, &t.Description); err != nil {
            log.Printf("GetTransactions scan error: %v", err)
            http.Error(w, "internal server error", http.StatusInternalServerError)
            return
        }
        transactions = append(transactions, t)
    }

    if err := rows.Err(); err != nil {
        log.Printf("GetTransactions rows error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(transactions)
}


func CreateTransaction(w http.ResponseWriter, r *http.Request) {
    var newTransaction models.Transaction

    if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // валидация
    if newTransaction.Amount <= 0 {
        http.Error(w, "amount must be > 0", http.StatusBadRequest)
        return
    }
    if newTransaction.Category == "" {
        http.Error(w, "category is required", http.StatusBadRequest)
        return
    }
    if newTransaction.Date == "" {
        http.Error(w, "date is required", http.StatusBadRequest)
        return
    }

    query := `INSERT INTO transactions (amount, category, date, description)
              VALUES ($1, $2, $3, $4) RETURNING id`

    err := database.DB.QueryRow(r.Context(), query,
        newTransaction.Amount,
        newTransaction.Category,
        newTransaction.Date,
        newTransaction.Description,
    ).Scan(&newTransaction.ID)

    if err != nil {
        log.Printf("CreateTransaction insert error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newTransaction)
}