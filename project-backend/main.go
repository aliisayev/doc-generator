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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	router.StaticFile("/", "../frontend/index.html")
	router.StaticFile("/style.css", "../frontend/style.css")

	router.POST("/generate", func(c *gin.Context) {
		var data ClientData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
			return
		}

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddUTF8Font("DejaVu", "", "fonts/DejaVuSans.ttf")
		pdf.AddUTF8Font("DejaVu", "B", "fonts/DejaVuSans-Bold.ttf")
		pdf.SetFont("DejaVu", "", 14)
		pdf.AddPage()

		pdf.Cell(0, 10, "Anket məlumatları")
		pdf.Ln(12)
		pdf.Cell(0, 10, "Soyad: "+data.LastName)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Ad: "+data.FirstName)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Ata adı: "+data.MiddleName)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Doğum tarixi: "+data.BirthDate)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Cins: "+data.Gender)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Telefon: "+data.Phone)
		pdf.Ln(12)

		pdf.SetFont("DejaVu", "B", 14)
		pdf.Cell(0, 10, "Suallar və cavablar:")
		pdf.Ln(10)
		pdf.SetFont("DejaVu", "", 13)

		for i := 1; i <= 15; i++ {
			key := "question" + strconv.Itoa(i)
			answer := data.Answers[key]
			if answer != "" {
				pdf.Cell(0, 8, "Sual "+strconv.Itoa(i)+": "+answer)
				pdf.Ln(6)
			}
		}

		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=agreement.pdf")
		if err := pdf.Output(c.Writer); err != nil {
			c.String(http.StatusInternalServerError, "Ошибка генерации PDF: "+err.Error())
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
