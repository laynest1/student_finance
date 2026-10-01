package handlers


import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/laynest1/student_finance/internal/auth"
    "github.com/laynest1/student_finance/internal/database"
    "github.com/laynest1/student_finance/internal/models"
)


func CreateCategory(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "неверный или недействительный токен", http.StatusUnauthorized)
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if category.Name == "" {
        http.Error(w, "название категории обязательно", http.StatusBadRequest)
        return
    }

    if category.Percentage <= 0 || category.Percentage > 100 {
        http.Error(w, "процент должен быть от 1 до 100", http.StatusBadRequest)
        return
    }

	queryChek := `SELECT COALESCE(SUM(percentage), 0) FROM categories WHERE user_id = $1`

	var totalPercentage int
	database.DB.QueryRow(r.Context(), queryChek, userID).Scan(&totalPercentage)
	if totalPercentage + category.Percentage > 100 {
		http.Error(w, "сумма процентов не может превышать 100", http.StatusBadRequest)
        return
    }

	query := `INSERT INTO categories(user_id, name, percentage) VALUES ($1, $2, $3) RETURNING id`

	err = database.DB.QueryRow(r.Context(), query, userID, category.Name, category.Percentage).Scan(&category.ID)
	if err != nil {
        log.Printf("CreateCategory error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

	category.UserID = userID
    category.CreatedAt = time.Now()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(category)


}

func GetCategories(w http.ResponseWriter, r *http.Request) {
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

	query := `SELECT id, user_id, name, percentage, created_at 
              FROM categories 
              WHERE user_id = $1 
              ORDER BY name`

	rows, err := database.DB.Query(r.Context(), query, userID)
	if err != nil {
		log.Printf("error get catogories : %v", err)
		http.Error(w, "internal serves error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []models.Category{}

	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Percentage, &c.CreatedAt)
		if err != nil {
			log.Printf("GetCategories scan error: %v", err)
            http.Error(w, "internal server error", http.StatusInternalServerError)
            return
        }
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)





	}

	func DeleteCategory(w http.ResponseWriter, r *http.Request) {
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
	
		urlPath := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		if urlPath == "" || urlPath == r.URL.Path {
			http.Error(w, "не указан id", http.StatusBadRequest)
			return
		}
		categoryID, err := strconv.Atoi(urlPath)
		if err != nil {
			http.Error(w, "неверный id", http.StatusBadRequest)
			return
		}
	
		query := `DELETE FROM categories WHERE user_id = $1 AND id = $2`
		result, err := database.DB.Exec(r.Context(), query, userID, categoryID)
		if err != nil {
			log.Printf("DeleteCategory error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	
		if result.RowsAffected() == 0 {
			http.Error(w, "категория не найдена", http.StatusNotFound)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "категория удалена"})
	}




	


