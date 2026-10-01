package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "strconv"
    "github.com/laynest1/student_finance/internal/auth"

    "github.com/laynest1/student_finance/internal/database"
    "github.com/laynest1/student_finance/internal/models"
)


func GetTransaction(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "требуется авторизация", http.StatusUnauthorized)
        return
    }
    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == authHeader {
        http.Error(w, "неверный формат токена", http.StatusUnauthorized)
        return
    }

    userID, err := auth.ValidateToken(token)
    if err != nil {
        http.Error(w, "недействительный токен", http.StatusUnauthorized)
        return
    }

    query := `SELECT id, user_id, amount, category, date, description 
	          FROM transactions 
	          WHERE user_id = $1
	          ORDER BY id DESC`

    rows, err := database.DB.Query(r.Context(), query, userID)
    if err != nil {
        log.Printf("GetTransaction query error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    transactions := []models.Transaction{}

    for rows.Next(){
        var t models.Transaction
        err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Category, &t.Date, &t.Description)
        if err != nil {
            log.Printf("get transaction scan error: %v", err)
            http.Error(w, "internal server error", http.StatusInternalServerError)
            return

        }
        transactions = append(transactions, t)
    
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(transactions)
}


func CreateTransaction(w http.ResponseWriter, r *http.Request) {

 
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "требуется авторизация", http.StatusUnauthorized)
        return
    }
    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == authHeader {
        http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
        return
    }

    userID, err := auth.ValidateToken(token)
    if err != nil{
        http.Error(w, "Недействительный токен", http.StatusUnauthorized)
        return
    }

    var newTransaction  models.Transaction
    err = json.NewDecoder(r.Body).Decode(&newTransaction)
    if err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }
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
    query := `INSERT INTO transactions (user_id, amount, category, date, description)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`


    err = database.DB.QueryRow(r.Context(), query, userID, newTransaction.Amount,
    newTransaction.Category, newTransaction.Date, newTransaction.Description).Scan(&newTransaction.ID)
    if err != nil {
		log.Printf("CreateTransaction insert error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    newTransaction.UserID = userID
    json.NewEncoder(w).Encode(newTransaction)

}



func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "требуется авторизация", http.StatusUnauthorized)
        return
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == authHeader {
        http.Error(w, "неверный формат токена", http.StatusUnauthorized)
        return
    }

    userID, err := auth.ValidateToken(token)
    if err != nil {
        http.Error(w, "недействительный токен", http.StatusUnauthorized)
        return
    }

    urlPath := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
    if urlPath == "" || urlPath == r.URL.Path {
        http.Error(w, "не указан айди", http.StatusBadRequest)
        return
    }

    tranactionID, err := strconv.Atoi(urlPath)
    if err != nil {
        http.Error(w, "не верный id", http.StatusBadRequest)
        return
    }

    query := `DELETE FROM transactions WHERE user_id = $1 AND id = $2`



    result, err := database.DB.Exec(r.Context(), query, userID, tranactionID)
    if err != nil {
        log.Printf("delete transaction error: %v", err)
        http.Error(w,"internal server error", http.StatusInternalServerError)
        return
    }

    delCnt := result.RowsAffected()
    
    if delCnt == 0 {
        http.Error(w,"транзакция не найдена или не принадлежит вам", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message" : "транзакция удалена",
    })


}


func GetStats(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "требуется авторизация", http.StatusUnauthorized)
        return
    }
    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == authHeader {
        http.Error(w, "неверный формат токена", http.StatusUnauthorized)
        return
    }

    userID, err:= auth.ValidateToken(token)
    if err != nil {
        http.Error(w, "неверный токен", http.StatusUnauthorized)
        return
    }
    queryTotal := `SELECT COALESCE(SUM(amount), 0), COUNT(*) from transactions where user_id = $1`


    var sum float64
    var cnt int
    err = database.DB.QueryRow(r.Context(), queryTotal, userID).Scan(&sum, &cnt)
    if err != nil {
		log.Printf("GetStats total error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

    queryGroup := `SELECT category, sum(amount) FROM transactions where user_id = $1 GROUP BY category`

    rows, err := database.DB.Query(r.Context(), queryGroup, userID)

    if err != nil {
        log.Printf("Stats category error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    defer rows.Close()

    byCategory := make(map[string]float64)
    for rows.Next() {
        var category string
        var sum float64
        err := rows.Scan(&category, &sum)
        if err != nil {
            log.Printf("stats scan category error: %v", err)
            http.Error(w, "internal server error", http.StatusInternalServerError)
            return
        }
        byCategory[category] = sum
    }
    stats := map[string]interface{}{
        "total_amount": sum,
        "transaction_count": cnt,
        "by_category": byCategory,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)



}

