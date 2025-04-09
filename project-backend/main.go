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
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FormData struct {
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	MiddleName  string   `json:"middle_name"`
	BirthDate   string   `json:"birth_date"`
	Phone       string   `json:"phone"`
	Gender      string   `json:"gender"`
	Email       string   `json:"email"`
	Address     string   `json:"address"`
	Citizenship string   `json:"citizenship"` // ← вот это!
	Photo       string   `json:"photo"`
	Answers     []string `json:"answers"`
	Today       string
}

func indexSafe(slice []string, i int) string {
	if i >= 0 && i < len(slice) {
		return slice[i]
	}
	return ""
}

func genderWord(gender string) string {
	if gender == "Kişi" {
		return " oğlu"
	}
	return " qızı"
}

func formatToday() string {
	return time.Now().Format("02.01.2006")
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

	data.Today = formatToday()

	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if strings.HasPrefix(d.Name(), "contract_") && strings.HasSuffix(d.Name(), ".html") {
			tmpl, err := template.New(d.Name()).
				Funcs(template.FuncMap{
					"indexSafe":  indexSafe,
					"genderWord": genderWord,
				}).
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

			pdfg, err := wkhtml.NewPDFGenerator()
			if err != nil {
				log.Printf("❌ wkhtmltopdf не найден: %v", err)
				return nil
			}

			page := wkhtml.NewPageReader(&filledHTML)
			page.EnableLocalFileAccess.Set(true)
			pdfg.AddPage(page)

			pdfg.PageSize.Set(wkhtml.PageSizeA4)
			pdfg.MarginTop.Set(50)
			pdfg.MarginBottom.Set(40)
			pdfg.MarginLeft.Set(15)
			pdfg.MarginRight.Set(15)

			err = pdfg.Create()
			if err != nil {
				log.Printf("❌ Ошибка генерации PDF %s: %v", d.Name(), err)
				return nil
			}

			pdfName := strings.TrimSuffix(d.Name(), ".html") + ".pdf"
			writer, err := zipWriter.Create(pdfName)
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, bytes.NewReader(pdfg.Bytes()))
			if err != nil {
				return err
			}
		}

		return nil
	})

	zipWriter.Close()

	if err != nil {
		http.Error(w, "Ошибка генерации PDF", http.StatusInternalServerError)
		return
	}

	// 💬 Telegram — отправляем только ОДИН zip файл
	bot, err := tgbotapi.NewBotAPI("7952271289:AAG3CaGezk_Rwr4o9LdIOj1HkKU3nLVsxvU")
	if err != nil {
		log.Printf("❌ Telegram Bot ошибка: %v", err)
		return
	}

	bot.Send(tgbotapi.NewMessage(1819220129, "📎 Müqavilələr ZIP faylında:"))

	zipFile := tgbotapi.FileBytes{
		Name:  "contracts.zip",
		Bytes: zipBuffer.Bytes(),
	}

	doc := tgbotapi.NewDocument(1819220129, zipFile)
	_, err = bot.Send(doc)
	if err != nil {
		log.Printf("❌ Ошибка отправки в Telegram: %v", err)
		return
	}

	w.Write([]byte("OK"))
}
