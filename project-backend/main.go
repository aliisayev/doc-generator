package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/unidoc/unioffice/document"
)

type ClientData struct {
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	MiddleName string   `json:"middleName"`
	BirthDate  string   `json:"birthDate"`
	Phone      string   `json:"phone"`
	Gender     string   `json:"gender"`
	Answers    []string `json:"answers"`
}

func fillTemplate(templatePath string, data ClientData) (*bytes.Buffer, error) {
	doc, err := document.Open(templatePath)
	if err != nil {
		return nil, err
	}

	placeholders := map[string]string{
		"{{.FirstName}}":  data.FirstName,
		"{{.LastName}}":   data.LastName,
		"{{.MiddleName}}": data.MiddleName,
		"{{.BirthDate}}":  data.BirthDate,
		"{{.Phone}}":      data.Phone,
		"{{.Gender}}":     data.Gender,
	}

	for i, answer := range data.Answers {
		key := fmt.Sprintf("{{.Answer%d}}", i+1)
		placeholders[key] = answer
	}

	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text := run.Text()
			for placeholder, value := range placeholders {
				if strings.Contains(text, placeholder) {
					text = strings.ReplaceAll(text, placeholder, value)
				}
			}
			run.ClearContent()
			run.AddText(text)
		}
	}

	var buf bytes.Buffer
	if err := doc.Save(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}

func main() {
	router := gin.Default()

	// Статические файлы
	router.StaticFile("/", "../frontend/index.html")
	router.StaticFile("/style.css", "../frontend/style.css")

	// Обработка формы
	router.POST("/generate", func(c *gin.Context) {
		var data ClientData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		files, err := filepath.Glob("templates/*.docx")
		if err != nil || len(files) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Шаблоны не найдены"})
			return
		}

		var zipBuffer bytes.Buffer
		zipWriter := zip.NewWriter(&zipBuffer)

		for _, file := range files {
			buf, err := fillTemplate(file, data)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации документа"})
				return
			}

			f, err := zipWriter.Create(filepath.Base(file))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания ZIP"})
				return
			}
			_, err = io.Copy(f, buf)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи файла в ZIP"})
				return
			}
		}

		zipWriter.Close()

		c.Header("Content-Disposition", "attachment; filename=documents.zip")
		c.Data(http.StatusOK, "application/zip", zipBuffer.Bytes())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
