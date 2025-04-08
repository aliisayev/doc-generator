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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
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

	// 📦 Отдача frontend
	router.StaticFile("/", "../frontend/index.html")
	router.StaticFile("/style.css", "../frontend/style.css")

	// 📬 Прием данных + генерация PDF
	router.POST("/generate", func(c *gin.Context) {
		var data ClientData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{"error": "Ошибка данных формы: " + err.Error()})
			return
		}

		files, err := filepath.Glob("templates/agreement_*.pdf")
		if err != nil || len(files) == 0 {
			c.JSON(500, gin.H{"error": "Файлы agreement_*.pdf не найдены в templates/"})
			return
		}

		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		for _, inputPath := range files {
			tmpOut := inputPath + "_out.pdf"

			if err := copyFile(inputPath, tmpOut); err != nil {
				fmt.Println("Ошибка копирования файла:", err)
				continue
			}

			replacements := map[string]string{
				"{{Имя}}":      data.FirstName,
				"{{Фамилия}}":  data.LastName,
				"{{Отчество}}": data.MiddleName,
				"{{Дата}}":     data.BirthDate,
				"{{Пол}}":      data.Gender,
				"{{Телефон}}":  data.Phone,
			}

			for i := 1; i <= 15; i++ {
				key := fmt.Sprintf("{{Вопрос%d}}", i)
				replacements[key] = data.Answers[fmt.Sprintf("question%d", i)]
			}

			for key, value := range replacements {
				if value == "" {
					continue
				}

				wm, err := pdfcpu.ParseTextWatermarkDetails(
					value,
					"scale:1, pos:tl, op:0.95, replace:"+key,
					true,
					types.POINTS,
				)
				if err != nil {
					fmt.Println("Ошибка watermark:", err)
					continue
				}

				err = api.AddWatermarksFile(tmpOut, tmpOut, nil, wm, model.NewDefaultConfiguration())
				if err != nil {
					fmt.Println("Ошибка замены:", err)
					continue
				}
			}

			modified, err := os.ReadFile(tmpOut)
			if err == nil {
				outFile, _ := zipWriter.Create(filepath.Base(inputPath))
				outFile.Write(modified)
			}
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
