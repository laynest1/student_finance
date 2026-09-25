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
	log.Println("✅ Таблица transactions готова!")
}