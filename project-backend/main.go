package main

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 👤 Структура данных формы
type FormData struct {
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	MiddleName string   `json:"middle_name"`
	BirthDate  string   `json:"birth_date"`
	Phone      string   `json:"phone"`
	Gender     string   `json:"gender"`
	Answers    []string `json:"answers"`
}

// 🧠 Безопасная функция index
func indexSafe(slice []string, i int) string {
	if i >= 0 && i < len(slice) {
		return slice[i]
	}
	return ""
}

func main() {
	// 📂 Обслуживаем frontend
	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	// 📥 Обработка формы
	http.HandleFunc("/submit", handleSubmit)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusBadRequest)
		return
	}

	outputDir := "output"
	os.MkdirAll(outputDir, os.ModePerm)

	templatesDir := "templates"

	err := filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if strings.HasPrefix(d.Name(), "contract_") && strings.HasSuffix(d.Name(), ".html") {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("❌ Ошибка чтения %s: %v", d.Name(), err)
				return nil
			}

			str := string(content)
			if strings.Contains(str, "{{") {
				// 🔧 Вот здесь ПОДКЛЮЧАЕМ indexSafe 💡
				tmpl, err := template.New(d.Name()).
					Funcs(template.FuncMap{"indexSafe": indexSafe}).
					Parse(str)
				if err != nil {
					log.Printf("❌ Ошибка шаблона %s: %v", d.Name(), err)
					return nil
				}

				outputName := "filled_" + d.Name()
				outputPath := filepath.Join(outputDir, outputName)

				outFile, err := os.Create(outputPath)
				if err != nil {
					log.Printf("❌ Ошибка при создании файла %s: %v", outputPath, err)
					return nil
				}
				defer outFile.Close()

				err = tmpl.Execute(outFile, data)
				if err != nil {
					log.Printf("❌ Ошибка при рендеринге %s: %v", d.Name(), err)
					return nil
				}

				log.Printf("✅ Сохранён: %s", outputPath)
			} else {
				log.Printf("ℹ️ Пропущен: %s (без плейсхолдеров)", d.Name())
			}
		}

		return nil
	})

	if err != nil {
		http.Error(w, "Ошибка генерации документов", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("🎉 Bütün sənədlər uğurla yaradıldı!"))
}
