package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FormData struct {
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	MiddleName string   `json:"middle_name"`
	BirthDate  string   `json:"birth_date"`
	Phone      string   `json:"phone"`
	Gender     string   `json:"gender"`
	Answers    []string `json:"answers"` // Ответы на вопросы
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	http.HandleFunc("/submit", handleSubmit)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Невозможно прочитать JSON", http.StatusBadRequest)
		return
	}

	replacements := map[string]string{
		"{{Имя}}":      data.FirstName,
		"{{Фамилия}}":  data.LastName,
		"{{Отчество}}": data.MiddleName,
		"{{Дата}}":     data.BirthDate,
		"{{Телефон}}":  data.Phone,
		"{{Пол}}":      data.Gender,
	}

	// Добавим ответы на {{Ответ1}}, {{Ответ2}}, ...
	for i, ans := range data.Answers {
		key := fmt.Sprintf("{{Ответ%d}}", i+1)
		replacements[key] = ans
	}

	inputDir := "templates"
	outputDir := "output"
	os.MkdirAll(outputDir, os.ModePerm)

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		http.Error(w, "Ошибка при чтении шаблонов", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "agreement_") && strings.HasSuffix(file.Name(), ".rtf") {
			inPath := filepath.Join(inputDir, file.Name())
			outPath := filepath.Join(outputDir, "filled_"+file.Name())

			err := fillRTF(inPath, outPath, replacements)
			if err != nil {
				log.Printf("Ошибка при заполнении %s: %v", file.Name(), err)
			}
		}
	}

	w.Write([]byte("Документы успешно созданы!"))
}

func fillRTF(inputPath, outputPath string, replacements map[string]string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	str := string(content)
	for placeholder, value := range replacements {
		str = strings.ReplaceAll(str, placeholder, value)
	}

	return os.WriteFile(outputPath, []byte(str), 0644)
}
