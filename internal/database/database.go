package database

import (
	"fmt"
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"os"

)
var DB *pgxpool.Pool


func Init() {
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Не удалось загрузить .env файл:", err)
	}
	connStr := os.Getenv("DATABASE_URL")



	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("не удалось подключиться к бд", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("бд недоступна ", err)
	}

	fmt.Println("подключились к бд")

}