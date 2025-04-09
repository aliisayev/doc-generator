package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type FormData struct {
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	MiddleName string   `json:"middle_name"`
	BirthDate  string   `json:"birth_date"`
	Phone      string   `json:"phone"`
	Gender     string   `json:"gender"`
	Answers    []string `json:"answers"`
}

func indexSafe(slice []string, i int) string {
	if i >= 0 && i < len(slice) {
		return slice[i]
	}
	return ""
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	http.HandleFunc("/submit", handleSubmit)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusBadRequest)
		return
	}

	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if strings.HasPrefix(d.Name(), "contract_") && strings.HasSuffix(d.Name(), ".html") {
			tmpl, err := template.New(d.Name()).
				Funcs(template.FuncMap{"indexSafe": indexSafe}).
				ParseFiles(path)
			if err != nil {
				log.Printf("❌ Ошибка шаблона %s: %v", d.Name(), err)
				return nil
			}

			var filledHTML bytes.Buffer
			err = tmpl.Execute(&filledHTML, data)
			if err != nil {
				log.Printf("❌ Ошибка рендера %s: %v", d.Name(), err)
				return nil
			}

			// PDF генерация
			pdfg, err := wkhtml.NewPDFGenerator()
			if err != nil {
				log.Printf("❌ wkhtmltopdf не найден: %v", err)
				return nil
			}

			page := wkhtml.NewPageReader(&filledHTML)
			page.EnableLocalFileAccess.Set(true)
			pdfg.AddPage(page)

			err = pdfg.Create()
			if err != nil {
				log.Printf("❌ Ошибка создания PDF для %s: %v", d.Name(), err)
				return nil
			}

			// Сохраняем PDF в ZIP
			pdfName := strings.TrimSuffix(d.Name(), ".html") + ".pdf"
			pdfFile, err := zipWriter.Create(pdfName)
			if err != nil {
				log.Printf("❌ Ошибка создания PDF-файла в ZIP: %v", err)
				return nil
			}

			_, err = io.Copy(pdfFile, bytes.NewReader(pdfg.Bytes()))
			if err != nil {
				log.Printf("❌ Ошибка записи PDF в ZIP: %v", err)
				return nil
			}

			log.Printf("✅ Добавлен PDF: %s", pdfName)
		}

		return nil
	})

	zipWriter.Close()

	if err != nil {
		http.Error(w, "Ошибка генерации PDF", http.StatusInternalServerError)
		return
	}

	// Отдаём архив пользователю
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=contracts.zip")
	w.Write(zipBuffer.Bytes())
}
