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

type FormData struct {
	FirstName  string            `json:"firstName"`
	LastName   string            `json:"lastName"`
	MiddleName string            `json:"middleName"`
	BirthDate  string            `json:"birthDate"`
	Gender     string            `json:"gender"`
	Phone      string            `json:"phone"`
	Answers    map[string]string `json:"answers"`
}

func main() {
	router := gin.Default()

	// Настройка CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.StaticFile("/", "./index.html")
	router.POST("/generate", handleGenerate)

	router.Run(":8080")
}

func handleGenerate(c *gin.Context) {
	var data FormData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templatesDir := "./project-backend/templates" // Исправленный путь к шаблонам
	outputFiles := []string{}

	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".docx") {
			return nil
		}

		doc, err := document.Open(path)
		if err != nil {
			return err
		}

		for _, para := range doc.Paragraphs() {
			for _, run := range para.Runs() {
				text := run.Text()
				text = strings.ReplaceAll(text, "{{.FirstName}}", data.FirstName)
				text = strings.ReplaceAll(text, "{{.LastName}}", data.LastName)
				text = strings.ReplaceAll(text, "{{.MiddleName}}", data.MiddleName)
				text = strings.ReplaceAll(text, "{{.BirthDate}}", data.BirthDate)
				text = strings.ReplaceAll(text, "{{.Gender}}", data.Gender)
				text = strings.ReplaceAll(text, "{{.Phone}}", data.Phone)

				// Вставляем ответы на тесты
				for key, val := range data.Answers {
					placeholder := fmt.Sprintf("{{.Q_%s}}", key)
					text = strings.ReplaceAll(text, placeholder, val)
				}

				run.ClearContent()
				run.AddText(text)
			}
		}

		outputPath := filepath.Join(os.TempDir(), info.Name())
		err = doc.SaveToFile(outputPath)
		if err != nil {
			return err
		}
		outputFiles = append(outputFiles, outputPath)
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Template processing failed: " + err.Error()})
		return
	}

	// Упаковываем в zip
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, file := range outputFiles {
		fileToZip, err := os.Open(file)
		if err != nil {
			continue
		}
		defer fileToZip.Close()

		info, _ := fileToZip.Stat()
		header, _ := zip.FileInfoHeader(info)
		header.Name = filepath.Base(file)
		writer, _ := zipWriter.CreateHeader(header)
		io.Copy(writer, fileToZip)
	}
	zipWriter.Close()

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=agreements.zip")
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}
