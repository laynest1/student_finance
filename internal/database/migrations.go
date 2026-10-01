package database

import (
	"context"
	"log"
)

func Migrate() {
	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(context.Background(), usersQuery)
	if err != nil {
		log.Fatal("ошибка миграции users:", err)
	}
	log.Println("таблица users создана")


	transactionsQuery := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		amount FLOAT NOT NULL,
		category VARCHAR(100) NOT NULL,
		date VARCHAR(20) NOT NULL,
		description TEXT NOT NULL
	);`

	_, err = DB.Exec(context.Background(), transactionsQuery)
	if err != nil {
		log.Fatal("Ошибка миграции transactions:", err)
	}
	log.Println("таблица transactions создана")

	_, err = DB.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            name VARCHAR(100) NOT NULL,
            percentage INTEGER NOT NULL CHECK (percentage > 0 AND percentage <= 100),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatalf("Ошибка создания таблицы categories: %v", err)
    }
    log.Println("таблица categories создана")
}