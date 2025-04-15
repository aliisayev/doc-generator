package main

import (
	"database/sql"
	"time"
)

// Contract — модель строки в таблице contracts
type Contract struct {
	ID        int
	CreatedAt time.Time
	Sent      bool
	SentAt    sql.NullTime
	PDFPath   string
	PhotoPath string
	JSONData  string
	Form      FormData // 💥 распарсенные данные из JSON
}

// FormData — структура формы, которую заполняет пользователь
type FormData struct {
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	MiddleName  string   `json:"middle_name"`
	BirthDate   string   `json:"birth_date"`
	Phone       string   `json:"phone"`
	Gender      string   `json:"gender"`
	Email       string   `json:"email"`
	Address     string   `json:"address"`
	Citizenship string   `json:"citizenship"`
	Photo       string   `json:"photo"`
	Answers     []string `json:"answers"`
	Today       string   `json:"today"`
}
