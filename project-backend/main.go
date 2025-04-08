package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/phpdave11/gofpdf"
)

type ClientData struct {
	FirstName  string            `json:"firstName"`
	LastName   string            `json:"lastName"`
	MiddleName string            `json:"middleName"`
	BirthDate  string            `json:"birthDate"`
	Phone      string            `json:"phone"`
	Gender     string            `json:"gender"`
	Answers    map[string]string `json:"answers"`
}

func main() {
	router := gin.Default()

	// ✅ CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// ✅ Отдача frontend
	router.StaticFile("/", "../frontend/index.html")
	router.StaticFile("/style.css", "../frontend/style.css")

	// ✅ Обработка запроса на генерацию PDF
	router.POST("/generate", func(c *gin.Context) {
		var data ClientData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
			return
		}

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "", 14)

		// Заголовок
		pdf.Cell(0, 10, "Анкета клиента")
		pdf.Ln(12)

		// Инфо
		pdf.Cell(0, 10, "Фамилия: "+data.LastName)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Имя: "+data.FirstName)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Отчество: "+data.MiddleName)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Дата рождения: "+data.BirthDate)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Пол: "+data.Gender)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Телефон: "+data.Phone)
		pdf.Ln(12)

		// Вопросы
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "Ответы на вопросы:")
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 13)

		for i := 1; i <= 15; i++ {
			key := "question" + strconv.Itoa(i)
			answer := data.Answers[key]
			if answer != "" {
				pdf.Cell(0, 8, "Вопрос "+strconv.Itoa(i)+": "+answer)
				pdf.Ln(6)
			}
		}

		// Вывод PDF
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=agreement.pdf")
		if err := pdf.Output(c.Writer); err != nil {
			c.String(http.StatusInternalServerError, "Ошибка генерации PDF: "+err.Error())
		}
	})

	// ✅ Запуск
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
