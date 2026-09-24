package database

import (
	"log"
	"context"
)

func Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		amount FLOAT NOT NULL,
		category VARCHAR(100) NOT NULL,
		date VARCHAR(20) NOT NULL,
		description TEXT NOT NULL
	);`

	_, err := DB.Exec(context.Background(), query)

	if err != nil {
		log.Fatal("ошибка миграции", err)
	}

	log.Println("таблица создана")
}