package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type ClientData struct {
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

	// Разрешаем CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	router.POST("/generate", func(c *gin.Context) {
		var data ClientData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{"error": "Ошибка данных формы: " + err.Error()})
			return
		}

		// Сканируем шаблоны PDF
		files, err := filepath.Glob("templates/agreement_*.pdf")
		if err != nil || len(files) == 0 {
			c.JSON(500, gin.H{"error": "Файлы agreement_*.pdf не найдены в templates/"})
			return
		}

		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		// Проходим по каждому PDF-шаблону
		for _, inputPath := range files {
			tmpOut := inputPath + "_out.pdf"

			// Копируем оригинал
			err := copyFile(inputPath, tmpOut)
			if err != nil {
				continue
			}

			// Список подстановок
			replacements := map[string]string{
				"{{Имя}}":      data.FirstName,
				"{{Фамилия}}":  data.LastName,
				"{{Отчество}}": data.MiddleName,
				"{{Дата}}":     data.BirthDate,
				"{{Пол}}":      data.Gender,
				"{{Телефон}}":  data.Phone,
			}

			// Добавим вопросы
			for i := 1; i <= 15; i++ {
				key := fmt.Sprintf("{{Вопрос%d}}", i)
				answer := data.Answers[fmt.Sprintf("question%d", i)]
				replacements[key] = answer
			}

			// Заменяем плейсхолдеры
			for key, value := range replacements {
				if value == "" {
					continue
				}
				api.AddWatermarksFile(
					tmpOut, tmpOut,
					nil,
					fmt.Sprintf("text:%s, scale:1, pos:tl, op:0.95, replace:%s", value, key),
					pdfcpu.DefaultConfiguration(),
				)
			}

			// Читаем финальный файл
			modified, _ := os.ReadFile(tmpOut)
			outFile, _ := zipWriter.Create(filepath.Base(inputPath))
			outFile.Write(modified)
			os.Remove(tmpOut)
		}

		zipWriter.Close()

		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", "attachment; filename=agreements.zip")
		c.Data(200, "application/zip", buf.Bytes())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// Копирование PDF
func copyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}
