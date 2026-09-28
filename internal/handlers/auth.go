package handlers

import(
	"encoding/json"
	"net/http"
	"time"
	"log"

	"github.com/laynest1/student_finance/internal/auth"
	"github.com/laynest1/student_finance/internal/database"
)

type RegisterRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`

}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err!= nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username и password обязательны", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "пароль должен быть минимум 6 символов", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	query := `INSERT INTO users (username, password, created_at) VALUES ($1, $2, $3) RETURNING id`


	var userID int
	err = database.DB.QueryRow(r.Context(), query,
		req.Username,
		hashedPassword,
		time.Now(),
	).Scan(&userID)

	if err != nil {
		log.Printf("Register error: %v", err)  // <-- Добавили логирование
		http.Error(w, "имя занято: "+err.Error(), http.StatusConflict)  // <-- Показываем реальную ошибку
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Пользователь успешно зарегистрирован",
		"user_id": userID,
	})



}


type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

func Login (w http.ResponseWriter, r* http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username и password обязательны", http.StatusBadRequest)
		return
	}

	var userID int
	var storedHash string
	query := `SELECT id, password FROM users WHERE username = $1`
	err := database.DB.QueryRow(r.Context(), query, req.Username).Scan(&userID, &storedHash)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPass(req.Password, storedHash) {
		http.Error(w, "пароль неверный, покушай ...", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token:  token,
		UserID: userID,
	})
}

func AllUsers (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application.json") 

	rows, err := database.DB.Query(r.Context(), "SELECT id, username, created_at FROM users")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var username string
		var createdAt time.Time
		err := rows.Scan(&id, &username, &createdAt)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		users = append(users, map[string]interface{}{
			"id":         id,
			"username":   username,
			"created_at": createdAt,
		})
	}

	json.NewEncoder(w).Encode(users)
}