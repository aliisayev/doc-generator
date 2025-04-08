package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

type FormData struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"`
	Phone      string `json:"phone"`
	Gender     string `json:"gender"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/submit", submitHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var data FormData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Ошибка при чтении формы", http.StatusBadRequest)
		return
	}

	templatePath := "templates/agreement_template.docx"
	outputPath := filepath.Join("output", fmt.Sprintf("contract_%s.docx", strings.ToLower(data.LastName)))

	err = fillTemplate(templatePath, outputPath, map[string]string{
		"{{Имя}}":      data.FirstName,
		"{{Фамилия}}":  data.LastName,
		"{{Отчество}}": data.MiddleName,
		"{{Дата}}":     data.BirthDate,
		"{{Телефон}}":  data.Phone,
		"{{Пол}}":      data.Gender,
	})

	if err != nil {
		http.Error(w, "Ошибка при создании документа", http.StatusInternalServerError)
		return
	}

	// Отдаём файл клиенту
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(outputPath))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	http.ServeFile(w, r, outputPath)
}

func fillTemplate(templatePath, outputPath string, data map[string]string) error {
	r, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		return err
	}
	doc := r.Editable()

	for placeholder, value := range data {
		doc.Replace(placeholder, value, -1)
	}

	os.MkdirAll("output", os.ModePerm)
	err = doc.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}
