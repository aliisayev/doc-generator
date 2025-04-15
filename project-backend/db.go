package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./contracts.db")
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS contracts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		json_data TEXT,
		pdf_path TEXT,
		photo_path TEXT,
		sent BOOLEAN DEFAULT FALSE,
		sent_at DATETIME
	);
	`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("❌ Ошибка создания таблицы:", err)
	}

	log.Println("✅ База данных инициализирована")
}
